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

package proc

import (
	"fmt"
	"reflect"

	"github.com/ontio/dad-go-eventbus/actor"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/config"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/ledger"
	tx "github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/events/message"
	hComm "github.com/ontio/dad-go/http/base/common"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
	"github.com/ontio/dad-go/smartcontract/service/neovm"
	tc "github.com/ontio/dad-go/txnpool/common"
	"github.com/ontio/dad-go/validator/types"
)

// NewTxActor creates an actor to handle the transaction-based messages from
// network and http
func NewTxActor(s *TXPoolServer) *TxActor {
	a := &TxActor{}
	a.setServer(s)
	return a
}

// NewTxPoolActor creates an actor to handle the messages from the consensus
func NewTxPoolActor(s *TXPoolServer) *TxPoolActor {
	a := &TxPoolActor{}
	a.setServer(s)
	return a
}

// NewVerifyRspActor creates an actor to handle the verified result from validators
func NewVerifyRspActor(s *TXPoolServer) *VerifyRspActor {
	a := &VerifyRspActor{}
	a.setServer(s)
	return a
}

// isBalanceEnough checks if the tranactor has enough to cover gas cost
func isBalanceEnough(address common.Address, gas uint64) bool {
	balance, err := hComm.GetContractBalance(0, utils.OngContractAddress, address)
	if err != nil {
		log.Debugf("failed to get contract balance %s err %v",
			address.ToHexString(), err)
		return false
	}
	return balance >= gas
}

func replyTxResult(txResultCh chan *tc.TxResult, hash common.Uint256,
	err errors.ErrCode, desc string) {
	result := &tc.TxResult{
		Err:  err,
		Hash: hash,
		Desc: desc,
	}
	select {
	case txResultCh <- result:
	default:
		log.Debugf("handleTransaction: duplicated result")
	}
}

// TxnActor: Handle the low priority msg from P2P and API
type TxActor struct {
	server *TXPoolServer
}

// handleTransaction handles a transaction from network and http
func (ta *TxActor) handleTransaction(sender tc.SenderType, self *actor.PID,
	txn *tx.Transaction, txResultCh chan *tc.TxResult) {
	ta.server.increaseStats(tc.RcvStats)
	if len(txn.ToArray()) > tc.MAX_TX_SIZE {
		log.Debugf("handleTransaction: reject a transaction due to size over 1M")
		if sender == tc.HttpSender && txResultCh != nil {
			replyTxResult(txResultCh, txn.Hash(), errors.ErrUnknown, "size is over 1M")
		}
		return
	}

	if ta.server.getTransaction(txn.Hash()) != nil {
		log.Debugf("handleTransaction: transaction %x already in the txn pool",
			txn.Hash())

		ta.server.increaseStats(tc.DuplicateStats)
		if sender == tc.HttpSender && txResultCh != nil {
			replyTxResult(txResultCh, txn.Hash(), errors.ErrDuplicateInput,
				fmt.Sprintf("transaction %x is already in the tx pool", txn.Hash()))
		}
	} else if ta.server.getTransactionCount() >= tc.MAX_CAPACITY {
		log.Debugf("handleTransaction: transaction pool is full for tx %x",
			txn.Hash())

		ta.server.increaseStats(tc.FailureStats)
		if sender == tc.HttpSender && txResultCh != nil {
			replyTxResult(txResultCh, txn.Hash(), errors.ErrTxPoolFull,
				"transaction pool is full")
		}
	} else {
		if _, overflow := common.SafeMul(txn.GasLimit, txn.GasPrice); overflow {
			log.Debugf("handleTransaction: gasLimit %v, gasPrice %v overflow",
				txn.GasLimit, txn.GasPrice)
			if sender == tc.HttpSender && txResultCh != nil {
				replyTxResult(txResultCh, txn.Hash(), errors.ErrUnknown,
					fmt.Sprintf("gasLimit %d * gasPrice %d overflow",
						txn.GasLimit, txn.GasPrice))
			}
			return
		}

		gasLimitConfig := config.DefConfig.Common.GasLimit
		gasPriceConfig := ta.server.getGasPrice()
		if txn.GasLimit < gasLimitConfig || txn.GasPrice < gasPriceConfig {
			log.Debugf("handleTransaction: invalid gasLimit %v, gasPrice %v",
				txn.GasLimit, txn.GasPrice)
			if sender == tc.HttpSender && txResultCh != nil {
				replyTxResult(txResultCh, txn.Hash(), errors.ErrUnknown,
					fmt.Sprintf("Please input gasLimit >= %d and gasPrice >= %d",
						gasLimitConfig, gasPriceConfig))
			}
			return
		}

		if txn.TxType == tx.Deploy && txn.GasLimit < neovm.CONTRACT_CREATE_GAS {
			log.Debugf("handleTransaction: deploy tx invalid gasLimit %v, gasPrice %v",
				txn.GasLimit, txn.GasPrice)
			if sender == tc.HttpSender && txResultCh != nil {
				replyTxResult(txResultCh, txn.Hash(), errors.ErrUnknown,
					fmt.Sprintf("Deploy tx gaslimit should >= %d",
						neovm.CONTRACT_CREATE_GAS))
			}
			return
		}

		if ta.server.preExec {
			result, err := ledger.DefLedger.PreExecuteContract(txn)
			if err != nil {
				log.Debugf("handleTransaction: failed to preExecuteContract tx %x err %v",
					txn.Hash(), err)
			}
			if txn.GasLimit < result.Gas {
				log.Debugf("handleTransaction: transaction's gasLimit %d is less than preExec gasLimit %d",
					txn.GasLimit, result.Gas)
				if sender == tc.HttpSender && txResultCh != nil {
					replyTxResult(txResultCh, txn.Hash(), errors.ErrUnknown,
						fmt.Sprintf("transaction's gasLimit %d is less than preExec gasLimit %d",
							txn.GasLimit, result.Gas))
				}
				return
			}
			gas, overflow := common.SafeMul(txn.GasPrice, result.Gas)
			if overflow {
				log.Debugf("handleTransaction: gasPrice %d preExec gasLimit %d overflow",
					txn.GasPrice, result.Gas)
				if sender == tc.HttpSender && txResultCh != nil {
					replyTxResult(txResultCh, txn.Hash(), errors.ErrUnknown,
						fmt.Sprintf("gasPrice %d * preExec gasLimit %d overflow",
							txn.GasPrice, result.Gas))
				}
				return
			}
			if !isBalanceEnough(txn.Payer, gas) {
				log.Debugf("handleTransaction: transactor %s has no balance enough to cover gas cost %d",
					txn.Payer.ToHexString(), gas)
				if sender == tc.HttpSender && txResultCh != nil {
					replyTxResult(txResultCh, txn.Hash(), errors.ErrUnknown,
						fmt.Sprintf("insufficient balance to cover gas cost %d", gas))
				}
				return
			}
			log.Debugf("handleTransaction: tx %x preExec success", txn.Hash())
		}
		<-ta.server.slots
		ta.server.assignTxToWorker(txn, sender, txResultCh)
	}
}

// Receive implements the actor interface
func (ta *TxActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("txpool-tx actor started and be ready to receive tx msg")

	case *actor.Stopping:
		log.Warn("txpool-tx actor stopping")

	case *actor.Restarting:
		log.Warn("txpool-tx actor restarting")

	case *tc.TxReq:
		sender := msg.Sender

		log.Debugf("txpool-tx actor receives tx from %v ", sender.Sender())

		ta.handleTransaction(sender, context.Self(), msg.Tx, msg.TxResultCh)

	case *tc.GetTxnReq:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives getting tx req from %v", sender)

		res := ta.server.getTransaction(msg.Hash)
		if sender != nil {
			sender.Request(&tc.GetTxnRsp{Txn: res},
				context.Self())
		}

	case *tc.GetTxnStats:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives getting tx stats from %v", sender)

		res := ta.server.getStats()
		if sender != nil {
			sender.Request(&tc.GetTxnStatsRsp{Count: res},
				context.Self())
		}

	case *tc.CheckTxnReq:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives checking tx req from %v", sender)

		res := ta.server.checkTx(msg.Hash)
		if sender != nil {
			sender.Request(&tc.CheckTxnRsp{Ok: res},
				context.Self())
		}

	case *tc.GetTxnStatusReq:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives getting tx status req from %v", sender)

		res := ta.server.getTxStatusReq(msg.Hash)
		if sender != nil {
			if res == nil {
				sender.Request(&tc.GetTxnStatusRsp{Hash: msg.Hash,
					TxStatus: nil}, context.Self())
			} else {
				sender.Request(&tc.GetTxnStatusRsp{Hash: res.Hash,
					TxStatus: res.Attrs}, context.Self())
			}
		}

	case *tc.GetTxnCountReq:
		sender := context.Sender()

		log.Debugf("txpool-tx actor receives getting tx count req from %v", sender)

		res := ta.server.getTxCount()
		if sender != nil {
			sender.Request(&tc.GetTxnCountRsp{Count: res},
				context.Self())
		}

	default:
		log.Debugf("txpool-tx actor: unknown msg %v type %v", msg, reflect.TypeOf(msg))
	}
}

func (ta *TxActor) setServer(s *TXPoolServer) {
	ta.server = s
}

// TxnPoolActor: Handle the high priority request from Consensus
type TxPoolActor struct {
	server *TXPoolServer
}

// Receive implements the actor interface
func (tpa *TxPoolActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("txpool actor started and be ready to receive txPool msg")

	case *actor.Stopping:
		log.Warn("txpool actor stopping")

	case *actor.Restarting:
		log.Warn("txpool actor Restarting")

	case *tc.GetTxnPoolReq:
		sender := context.Sender()

		log.Debugf("txpool actor receives getting tx pool req from %v", sender)

		res := tpa.server.getTxPool(msg.ByCount, msg.Height)
		if sender != nil {
			sender.Request(&tc.GetTxnPoolRsp{TxnPool: res}, context.Self())
		}

	case *tc.GetPendingTxnReq:
		sender := context.Sender()

		log.Debugf("txpool actor receives getting pedning tx req from %v", sender)

		res := tpa.server.getPendingTxs(msg.ByCount)
		if sender != nil {
			sender.Request(&tc.GetPendingTxnRsp{Txs: res}, context.Self())
		}

	case *tc.VerifyBlockReq:
		sender := context.Sender()

		log.Debugf("txpool actor receives verifying block req from %v", sender)

		tpa.server.verifyBlock(msg, sender)

	case *message.SaveBlockCompleteMsg:
		sender := context.Sender()

		log.Debugf("txpool actor receives block complete event from %v", sender)

		if msg.Block != nil {
			tpa.server.cleanTransactionList(msg.Block.Transactions, msg.Block.Header.Height)
		}

	default:
		log.Debugf("txpool actor: unknown msg %v type %v", msg, reflect.TypeOf(msg))
	}
}

func (tpa *TxPoolActor) setServer(s *TXPoolServer) {
	tpa.server = s
}

// VerifyRspActor: Handle the response from the validators
type VerifyRspActor struct {
	server *TXPoolServer
}

// Receive implements the actor interface
func (vpa *VerifyRspActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("txpool-verify actor: started and be ready to receive validator's msg")

	case *actor.Stopping:
		log.Warn("txpool-verify actor: stopping")

	case *actor.Restarting:
		log.Warn("txpool-verify actor: Restarting")

	case *types.RegisterValidator:
		log.Debugf("txpool-verify actor:: validator %v connected", msg.Sender)
		vpa.server.registerValidator(msg)

	case *types.UnRegisterValidator:
		log.Debugf("txpool-verify actor:: validator %d:%v disconnected", msg.Type, msg.Id)

		vpa.server.unRegisterValidator(msg.Type, msg.Id)

	case *types.CheckResponse:
		log.Debug("txpool-verify actor:: Receives verify rsp message")

		vpa.server.assignRspToWorker(msg)

	default:
		log.Debugf("txpool-verify actor:Unknown msg %v type %v", msg, reflect.TypeOf(msg))
	}
}

func (vpa *VerifyRspActor) setServer(s *TXPoolServer) {
	vpa.server = s
}
