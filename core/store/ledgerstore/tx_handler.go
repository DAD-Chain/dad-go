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

package ledgerstore

import (
	"bytes"
	"fmt"
	"math"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/config"
	"github.com/ontio/dad-go/common/serialization"
	"github.com/ontio/dad-go/core/payload"
	"github.com/ontio/dad-go/core/states"
	"github.com/ontio/dad-go/core/store"
	scommon "github.com/ontio/dad-go/core/store/common"
	"github.com/ontio/dad-go/core/store/statestore"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/smartcontract"
	"github.com/ontio/dad-go/smartcontract/event"
	ninit "github.com/ontio/dad-go/smartcontract/service/native/init"
	"github.com/ontio/dad-go/smartcontract/service/native/ont"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
	"github.com/ontio/dad-go/smartcontract/service/neovm"
	sstates "github.com/ontio/dad-go/smartcontract/states"
	"github.com/ontio/dad-go/smartcontract/storage"
	stypes "github.com/ontio/dad-go/smartcontract/types"
)

//HandleDeployTransaction deal with smart contract deploy transaction
func (self *StateStore) HandleDeployTransaction(store store.LedgerStore, stateBatch *statestore.StateBatch,
	tx *types.Transaction, block *types.Block, eventStore scommon.EventStore) error {
	deploy := tx.Payload.(*payload.DeployCode)
	txHash := tx.Hash()
	originAddress := deploy.Code.AddressFromVmCode()

	var (
		notifies []*event.NotifyEventInfo
		err      error
	)
	// mapping native contract origin address to target address
	if deploy.Code.VmType == stypes.Native {
		targetAddress, err := common.AddressParseFromBytes(deploy.Code.Code)
		if err != nil {
			return fmt.Errorf("Invalid native contract address:%s", err)
		}
		originAddress = targetAddress
	} else {
		if err := isBalanceSufficient(tx, stateBatch); err != nil {
			return err
		}

		cache := storage.NewCloneCache(stateBatch)

		// init smart contract configuration info
		config := &smartcontract.Config{
			Time:   block.Header.Timestamp,
			Height: block.Header.Height,
			Tx:     tx,
		}

		notifies, err = costGas(tx.Payer, tx.GasLimit*tx.GasPrice, config, cache, store)
		if err != nil {
			return err
		}
		cache.Commit()
	}

	// store contract message
	err = stateBatch.TryGetOrAdd(scommon.ST_CONTRACT, originAddress[:], deploy)
	if err != nil {
		return err
	}

	SaveNotify(eventStore, txHash, notifies, true)
	return nil
}

//HandleInvokeTransaction deal with smart contract invoke transaction
func (self *StateStore) HandleInvokeTransaction(store store.LedgerStore, stateBatch *statestore.StateBatch,
	tx *types.Transaction, block *types.Block, eventStore scommon.EventStore) error {
	invoke := tx.Payload.(*payload.InvokeCode)
	txHash := tx.Hash()
	code := invoke.Code.Code
	sysTransFlag := bytes.Compare(code, ninit.COMMIT_DPOS_BYTES) == 0 || block.Header.Height == 0

	if !sysTransFlag && tx.GasPrice != 0 {
		if err := isBalanceSufficient(tx, stateBatch); err != nil {
			return err
		}
	}

	// init smart contract configuration info
	config := &smartcontract.Config{
		Time:   block.Header.Timestamp,
		Height: block.Header.Height,
		Tx:     tx,
	}

	cache := storage.NewCloneCache(stateBatch)
	//init smart contract info
	sc := smartcontract.SmartContract{
		Config:     config,
		CloneCache: cache,
		Store:      store,
		Code:       invoke.Code,
		Gas:        tx.GasLimit,
	}

	//start the smart contract executive function
	_, err := sc.Execute()

	if err != nil {
		return err
	}

	var notifies []*event.NotifyEventInfo
	if !sysTransFlag {
		totalGas := tx.GasLimit - sc.Gas
		if totalGas < neovm.TRANSACTION_GAS {
			totalGas = neovm.TRANSACTION_GAS
		}
		notifies, err = costGas(tx.Payer, totalGas*tx.GasPrice, config, sc.CloneCache, store)
		if err != nil {
			return err
		}
	}

	SaveNotify(eventStore, txHash, append(sc.Notifications, notifies...), true)
	sc.CloneCache.Commit()
	return nil
}

func SaveNotify(eventStore scommon.EventStore, txHash common.Uint256, notifies []*event.NotifyEventInfo, execSucc bool) error {
	if !config.DefConfig.Common.EnableEventLog {
		return nil
	}
	var notifyInfo *event.ExecuteNotify
	if execSucc {
		notifyInfo = &event.ExecuteNotify{TxHash: txHash,
			State: event.CONTRACT_STATE_SUCCESS, Notify: notifies}
	} else {
		notifyInfo = &event.ExecuteNotify{TxHash: txHash,
			State: event.CONTRACT_STATE_FAIL, Notify: notifies}
	}
	if err := eventStore.SaveEventNotifyByTx(txHash, notifyInfo); err != nil {
		return fmt.Errorf("SaveEventNotifyByTx error %s", err)
	}
	event.PushSmartCodeEvent(txHash, 0, event.EVENT_NOTIFY, notifyInfo)
	return nil
}

//HandleClaimTransaction deal with ong claim transaction
func (self *StateStore) HandleClaimTransaction(stateBatch *statestore.StateBatch, tx *types.Transaction) error {
	//TODO
	return nil
}

//HandleVoteTransaction deal with vote transaction
func (self *StateStore) HandleVoteTransaction(stateBatch *statestore.StateBatch, tx *types.Transaction) error {
	vote := tx.Payload.(*payload.Vote)
	buf := new(bytes.Buffer)
	vote.Account.Serialize(buf)
	stateBatch.TryAdd(scommon.ST_VOTE, buf.Bytes(), &states.VoteState{PublicKeys: vote.PubKeys})
	return nil
}

func genNativeTransferCode(contract, from, to common.Address, value uint64) stypes.VmCode {
	transfer := ont.Transfers{States: []*ont.State{{From: from, To: to, Value: value}}}
	tr := new(bytes.Buffer)
	transfer.Serialize(tr)
	trans := &sstates.Contract{
		Address: contract,
		Method:  "transfer",
		Args:    tr.Bytes(),
	}
	ts := new(bytes.Buffer)
	trans.Serialize(ts)
	return stypes.VmCode{Code: ts.Bytes(), VmType: stypes.Native}
}

// check whether payer ong balance sufficient
func isBalanceSufficient(tx *types.Transaction, stateBatch *statestore.StateBatch) error {
	balance, err := getBalance(stateBatch, tx.Payer, utils.OngContractAddress)
	if err != nil {
		return err
	}
	if balance < tx.GasLimit*tx.GasPrice {
		return fmt.Errorf("payer gas insufficient, need %d , only have %d", tx.GasLimit*tx.GasPrice, balance)
	}
	return nil
}

func costGas(payer common.Address, gas uint64, config *smartcontract.Config,
	cache *storage.CloneCache, store store.LedgerStore) ([]*event.NotifyEventInfo, error) {

	nativeTransferCode := genNativeTransferCode(utils.OngContractAddress, payer,
		utils.GovernanceContractAddress, gas)

	sc := smartcontract.SmartContract{
		Config:     config,
		CloneCache: cache,
		Store:      store,
		Code:       nativeTransferCode,
		Gas:        math.MaxUint64,
	}

	_, err := sc.Execute()

	if err != nil {
		return nil, err
	}
	return sc.Notifications, nil
}

func getBalance(stateBatch *statestore.StateBatch, address, contract common.Address) (uint64, error) {
	bl, err := stateBatch.TryGet(scommon.ST_STORAGE, append(contract[:], address[:]...))
	if err != nil {
		return 0, fmt.Errorf("get balance error:%s", err)
	}
	if bl == nil || bl.Value == nil {
		return 0, fmt.Errorf("get %s balance fail from %s", address.ToHexString(), contract.ToHexString())
	}
	item, ok := bl.Value.(*states.StorageItem)
	if !ok {
		return 0, fmt.Errorf("%s", "instance doesn't StorageItem!")
	}
	balance, err := serialization.ReadUint64(bytes.NewBuffer(item.Value))
	if err != nil {
		return 0, fmt.Errorf("read balance error:%s", err)
	}
	return balance, nil
}
