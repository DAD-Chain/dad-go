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

package ledger

import (
	"fmt"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/genesis"
	"github.com/ontio/dad-go/core/payload"
	"github.com/ontio/dad-go/core/states"
	"github.com/ontio/dad-go/core/store"
	"github.com/ontio/dad-go/core/store/ledgerstore"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/smartcontract/event"
)

var DefLedger *Ledger

type Ledger struct {
	ldgStore store.LedgerStore
}

func NewLedger() (*Ledger, error) {
	ldgStore, err := ledgerstore.NewLedgerStore()
	if err != nil {
		return nil, fmt.Errorf("NewLedgerStore error %s", err)
	}
	return &Ledger{
		ldgStore: ldgStore,
	}, nil
}

func (self *Ledger) GetStore() store.LedgerStore {
	return self.ldgStore
}

func (self *Ledger) Init(defaultBookkeeper []keypair.PublicKey) error {
	genesisBlock, err := genesis.GenesisBlockInit(defaultBookkeeper)
	if err != nil {
		return fmt.Errorf("genesisBlock error %s", err)
	}
	err = self.ldgStore.InitLedgerStoreWithGenesisBlock(genesisBlock, defaultBookkeeper)
	if err != nil {
		return fmt.Errorf("InitLedgerStoreWithGenesisBlock error %s", err)
	}
	return nil
}

func (self *Ledger) AddHeaders(headers []*types.Header) error {
	return self.ldgStore.AddHeaders(headers)
}

func (self *Ledger) AddBlock(block *types.Block) error {
	err := self.ldgStore.AddBlock(block)
	if err != nil {
		log.Errorf("Ledger AddBlock BlockHeight:%d BlockHash:%x error:%s", block.Header.Height, block.Hash(), err)
	}
	return err
}

func (self *Ledger) GetBlockRootWithNewTxRoot(txRoot common.Uint256) common.Uint256 {
	return self.ldgStore.GetBlockRootWithNewTxRoot(txRoot)
}

func (self *Ledger) GetBlockByHeight(height uint32) (*types.Block, error) {
	return self.ldgStore.GetBlockByHeight(height)
}

func (self *Ledger) GetBlockByHash(blockHash common.Uint256) (*types.Block, error) {
	return self.ldgStore.GetBlockByHash(blockHash)
}

func (self *Ledger) GetHeaderByHeight(height uint32) (*types.Header, error) {
	return self.ldgStore.GetHeaderByHeight(height)
}

func (self *Ledger) GetHeaderByHash(blockHash common.Uint256) (*types.Header, error) {
	return self.ldgStore.GetHeaderByHash(blockHash)
}

func (self *Ledger) GetBlockHash(height uint32) common.Uint256 {
	return self.ldgStore.GetBlockHash(height)
}

func (self *Ledger) GetTransaction(txHash common.Uint256) (*types.Transaction, error) {
	tx, _, err := self.ldgStore.GetTransaction(txHash)
	return tx, err
}

func (self *Ledger) GetTransactionWithHeight(txHash common.Uint256) (*types.Transaction, uint32, error) {
	return self.ldgStore.GetTransaction(txHash)
}

func (self *Ledger) GetCurrentBlockHeight() uint32 {
	return self.ldgStore.GetCurrentBlockHeight()
}

func (self *Ledger) GetCurrentBlockHash() common.Uint256 {
	return self.ldgStore.GetCurrentBlockHash()
}

func (self *Ledger) GetCurrentHeaderHeight() uint32 {
	return self.ldgStore.GetCurrentHeaderHeight()
}

func (self *Ledger) GetCurrentHeaderHash() common.Uint256 {
	return self.ldgStore.GetCurrentHeaderHash()
}

func (self *Ledger) IsContainTransaction(txHash common.Uint256) (bool, error) {
	return self.ldgStore.IsContainTransaction(txHash)
}

func (self *Ledger) IsContainBlock(blockHash common.Uint256) (bool, error) {
	return self.ldgStore.IsContainBlock(blockHash)
}

func (self *Ledger) GetCurrentStateRoot() (common.Uint256, error) {
	return common.Uint256{}, nil
}

func (self *Ledger) GetBookkeeperState() (*states.BookkeeperState, error) {
	return self.ldgStore.GetBookkeeperState()
}

func (self *Ledger) GetStorageItem(codeHash common.Address, key []byte) ([]byte, error) {
	storageKey := &states.StorageKey{
		CodeHash: codeHash,
		Key:      key,
	}
	storageItem, err := self.ldgStore.GetStorageItem(storageKey)
	if err != nil {
		return nil, fmt.Errorf("GetStorageItem error %s", err)
	}
	if storageItem == nil {
		return nil, nil
	}
	return storageItem.Value, nil
}

func (self *Ledger) FindStorageItem(codeHash common.Address, key []byte) ([][]byte, error) {
	storageKey := &states.StorageKey{
		CodeHash: codeHash,
		Key:      key,
	}
	storageItem, err := self.ldgStore.FindStorageItem(storageKey)
	if err != nil {
		return nil, fmt.Errorf("FindStorageItem error %s", err)
	}
	var value [][]byte
	for _, storageitem := range storageItem {
		value = append(value, storageitem.Value)
	}
	return value, nil
}

func (self *Ledger) GetContractState(contractHash common.Address) (*payload.DeployCode, error) {
	return self.ldgStore.GetContractState(contractHash)
}

func (self *Ledger) GetMerkleProof(proofHeight, rootHeight uint32) ([]common.Uint256, error) {
	return self.ldgStore.GetMerkleProof(proofHeight, rootHeight)
}

func (self *Ledger) PreExecuteContract(tx *types.Transaction) (interface{}, error) {
	return self.ldgStore.PreExecuteContract(tx)
}

func (self *Ledger) GetEventNotifyByTx(tx common.Uint256) (*event.ExecuteNotify, error) {
	return self.ldgStore.GetEventNotifyByTx(tx)
}

func (self *Ledger) GetEventNotifyByBlock(height uint32) ([]common.Uint256, error) {
	return self.ldgStore.GetEventNotifyByBlock(height)
}

func (self *Ledger) Close() error {
	return self.ldgStore.Close()
}
