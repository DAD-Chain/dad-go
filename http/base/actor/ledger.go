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

package actor

import (
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/ledger"
	"github.com/ontio/dad-go/core/payload"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/smartcontract/event"
	cstate "github.com/ontio/dad-go/smartcontract/states"
)

const (
	REQ_TIMEOUT    = 5
	ERR_ACTOR_COMM = "[http] Actor comm error: %v"
)

//GetHeaderByHeight from ledger
func GetHeaderByHeight(height uint32) (*types.Header, error) {
	return ledger.DefLedger.GetHeaderByHeight(height)
}

//GetBlockByHeight from ledger
func GetBlockByHeight(height uint32) (*types.Block, error) {
	return ledger.DefLedger.GetBlockByHeight(height)
}

//GetBlockHashFromStore from ledger
func GetBlockHashFromStore(height uint32) common.Uint256 {
	return ledger.DefLedger.GetBlockHash(height)
}

//CurrentBlockHash from ledger
func CurrentBlockHash() common.Uint256 {
	return ledger.DefLedger.GetCurrentBlockHash()
}

//GetBlockFromStore from ledger
func GetBlockFromStore(hash common.Uint256) (*types.Block, error) {
	return ledger.DefLedger.GetBlockByHash(hash)
}

//GetCurrentBlockHeight from ledger
func GetCurrentBlockHeight() uint32 {
	return ledger.DefLedger.GetCurrentBlockHeight()
}

//GetTransaction from ledger
func GetTransaction(hash common.Uint256) (*types.Transaction, error) {
	return ledger.DefLedger.GetTransaction(hash)
}

//GetStorageItem from ledger
func GetStorageItem(address common.Address, key []byte) ([]byte, error) {
	return ledger.DefLedger.GetStorageItem(address, key)
}

//GetContractStateFromStore from ledger
func GetContractStateFromStore(hash common.Address) (*payload.DeployCode, error) {
	hash = updateNativeSCAddr(hash)
	return ledger.DefLedger.GetContractState(hash)
}

//GetTxnWithHeightByTxHash from ledger
func GetTxnWithHeightByTxHash(hash common.Uint256) (uint32, *types.Transaction, error) {
	tx, height, err := ledger.DefLedger.GetTransactionWithHeight(hash)
	return height, tx, err
}

//PreExecuteContract from ledger
func PreExecuteContract(tx *types.Transaction) (*cstate.PreExecResult, error) {
	return ledger.DefLedger.PreExecuteContract(tx)
}

//GetEventNotifyByTxHash from ledger
func GetEventNotifyByTxHash(txHash common.Uint256) (*event.ExecuteNotify, error) {
	return ledger.DefLedger.GetEventNotifyByTx(txHash)
}

//GetEventNotifyByHeight from ledger
func GetEventNotifyByHeight(height uint32) ([]*event.ExecuteNotify, error) {
	return ledger.DefLedger.GetEventNotifyByBlock(height)
}

//GetMerkleProof from ledger
func GetMerkleProof(proofHeight uint32, rootHeight uint32) ([]common.Uint256, error) {
	return ledger.DefLedger.GetMerkleProof(proofHeight, rootHeight)
}
