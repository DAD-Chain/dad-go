/*
 * Copyright (C) 2018 The dad-go Authors
 * This file is part of The dad-go library.
 *
 * The dad-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The dad-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */

package solo

import (
	"fmt"
	ldgactor "github.com/dad-go/core/ledger/actor"
	"reflect"
	"time"

	"github.com/dad-go/account"
	. "github.com/dad-go/common"
	"github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	actorTypes "github.com/dad-go/consensus/actor"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/signature"
	"github.com/dad-go/core/types"
	"github.com/dad-go/events"
	"github.com/dad-go/events/message"
	"github.com/dad-go/validator/increment"
	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go-eventbus/actor"
)

/*
*Simple consensus for solo node in test environment.
 */
var GenBlockTime = (config.DEFAULTGENBLOCKTIME * time.Second)

const ContextVersion uint32 = 0

type SoloService struct {
	Account       *account.Account
	poolActor     *actorTypes.TxPoolActor
	ledgerActor   *actorTypes.LedgerActor
	incrValidator *increment.IncrementValidator
	existCh       chan interface{}

	pid *actor.PID
	sub *events.ActorSubscriber
}

func NewSoloService(bkAccount *account.Account, txpool *actor.PID, ledger *actor.PID) (*SoloService, error) {
	service := &SoloService{
		Account:       bkAccount,
		poolActor:     &actorTypes.TxPoolActor{Pool: txpool},
		ledgerActor:   &actorTypes.LedgerActor{Ledger: ledger},
		incrValidator: increment.NewIncrementValidator(10),
	}

	props := actor.FromProducer(func() actor.Actor {
		return service
	})

	pid, err := actor.SpawnNamed(props, "consensus_solo")
	service.pid = pid
	service.sub = events.NewActorSubscriber(pid)

	return service, err
}

func (this *SoloService) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Restarting:
		log.Warn("solo actor restarting")
	case *actor.Stopping:
		log.Warn("solo actor stopping")
	case *actor.Stopped:
		log.Warn("solo actor stopped")
	case *actor.Started:
		log.Warn("solo actor started")
	case *actor.Restart:
		log.Warn("solo actor restart")
	case *actorTypes.StartConsensus:
		if this.existCh != nil {
			log.Warn("consensus have started")
			return
		}

		this.sub.Subscribe(message.TopicSaveBlockComplete)

		timer := time.NewTicker(GenBlockTime)
		this.existCh = make(chan interface{})
		go func() {
			defer timer.Stop()
			existCh := this.existCh
			for {
				select {
				case <-timer.C:
					this.pid.Tell(&actorTypes.TimeOut{})
				case <-existCh:
					return
				}
			}
		}()
	case *actorTypes.StopConsensus:
		if this.existCh != nil {
			close(this.existCh)
			this.existCh = nil
			this.incrValidator.Clean()
			this.sub.Unsubscribe(message.TopicSaveBlockComplete)
		}
	case *message.SaveBlockCompleteMsg:
		log.Info("solo actor receives block complete event. block height=", msg.Block.Header.Height)
		this.incrValidator.AddBlock(msg.Block)

	case *actorTypes.TimeOut:
		err := this.genBlock()
		if err != nil {
			log.Errorf("Solo genBlock error %s", err)
		}
	default:
		log.Info("solo actor: Unknown msg ", msg, "type", reflect.TypeOf(msg))
	}
}

func (this *SoloService) GetPID() *actor.PID {
	return this.pid
}

func (this *SoloService) Start() error {
	this.pid.Tell(&actorTypes.StartConsensus{})
	return nil
}

func (this *SoloService) Halt() error {
	this.pid.Tell(&actorTypes.StopConsensus{})
	return nil
}

func (this *SoloService) genBlock() error {
	block, err := this.makeBlock()
	if err != nil {
		return fmt.Errorf("makeBlock error %s", err)
	}

	future := ldgactor.DefLedgerPid.RequestFuture(&ldgactor.AddBlockReq{Block: block}, 30*time.Second)
	result, err := future.Result()
	if err != nil {
		return fmt.Errorf("genBlock DefLedgerPid.RequestFuture Height:%d error:%s", block.Header.Height, err)
	}
	addBlockRsp := result.(*ldgactor.AddBlockRsp)
	if addBlockRsp.Error != nil {
		return fmt.Errorf("AddBlockRsp Height:%d error:%s", block.Header.Height, addBlockRsp.Error)
	}
	return nil
}

func (this *SoloService) makeBlock() (*types.Block, error) {
	log.Debug()
	owner := this.Account.PublicKey
	nextBookkeeper, err := types.AddressFromBookkeepers([]keypair.PublicKey{owner})
	if err != nil {
		return nil, fmt.Errorf("GetBookkeeperAddress error:%s", err)
	}
	prevHash := ledger.DefLedger.GetCurrentBlockHash()
	height := ledger.DefLedger.GetCurrentBlockHeight()

	validHeight := height

	start, end := this.incrValidator.BlockRange()

	if height+1 == end {
		validHeight = start
	} else {
		this.incrValidator.Clean()
		log.Infof("incr validator block height %v != ledger block height %v", end-1, height)
	}

	log.Infof("current block Height %v, incrValidateHeight %v", height, validHeight)

	txs := this.poolActor.GetTxnPool(true, validHeight)
	// todo : fix feesum calcuation
	feeSum := Fixed64(0)

	// TODO: increment checking txs

	nonce := GetNonce()
	txBookkeeping := this.createBookkeepingTransaction(nonce, feeSum)

	transactions := make([]*types.Transaction, 0, len(txs)+1)
	transactions = append(transactions, txBookkeeping)
	for _, txEntry := range txs {
		// TODO optimize to use height in txentry
		if err := this.incrValidator.Verify(txEntry.Tx, validHeight); err == nil {
			transactions = append(transactions, txEntry.Tx)
		}
	}

	txHash := []Uint256{}
	for _, t := range transactions {
		txHash = append(txHash, t.Hash())
	}
	txRoot, err := ComputeRoot(txHash)
	if err != nil {
		return nil, fmt.Errorf("ComputeRoot error:%s", err)
	}

	blockRoot := ledger.DefLedger.GetBlockRootWithNewTxRoot(txRoot)
	header := &types.Header{
		Version:          ContextVersion,
		PrevBlockHash:    prevHash,
		TransactionsRoot: txRoot,
		BlockRoot:        blockRoot,
		Timestamp:        uint32(time.Now().Unix()),
		Height:           height + 1,
		ConsensusData:    nonce,
		NextBookkeeper:   nextBookkeeper,
	}
	block := &types.Block{
		Header:       header,
		Transactions: transactions,
	}

	blockHash := block.Hash()

	sig, err := signature.Sign(this.Account.PrivKey(), blockHash[:])
	if err != nil {
		return nil, fmt.Errorf("[Signature],Sign error:%s.", err)
	}

	block.Header.Bookkeepers = []keypair.PublicKey{owner}
	block.Header.SigData = [][]byte{sig}
	return block, nil
}

func (this *SoloService) createBookkeepingTransaction(nonce uint64, fee Fixed64) *types.Transaction {
	log.Debug()
	//TODO: sysfee
	bookKeepingPayload := &payload.BookKeeping{
		Nonce: uint64(time.Now().UnixNano()),
	}
	tx := &types.Transaction{
		TxType:     types.BookKeeping,
		Payload:    bookKeepingPayload,
		Attributes: []*types.TxAttribute{},
	}
	txHash := tx.Hash()
	acc := this.Account
	s, err := signature.Sign(acc.PrivateKey, txHash[:])
	if err != nil {
		return nil
	}
	sig := &types.Sig{
		PubKeys: []keypair.PublicKey{acc.PublicKey},
		M:       1,
		SigData: [][]byte{s},
	}
	tx.Sigs = []*types.Sig{sig}
	return tx
}
