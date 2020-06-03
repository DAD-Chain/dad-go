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

package common

import (
	"github.com/dad-go/common"
	"github.com/dad-go/core/types"
	"github.com/dad-go/errors"
)

const (
	MAXCAPACITY    = 100140                       // The tx pool's capacity that holds the verified txs
	MAXPENDINGTXN  = 2048 * 10                    // The max length of pending txs
	MAXWORKERNUM   = 2                            // The max concurrent workers
	MAXRCVTXNLEN   = MAXWORKERNUM * MAXPENDINGTXN // The max length of the queue that server can hold
	MAXRETRIES     = 0                            // The retry times to verify tx
	EXPIREINTERVAL = 9                            // The timeout that verify tx
	STATELESSMASK  = 0x1                          // The mask of stateless validator
	STATEFULMASK   = 0x2                          // The mask of stateful validator
	VERIFYMASK     = STATELESSMASK | STATEFULMASK
	MAXLIMITATION  = 10000
)

type ActorType uint8

const (
	_              ActorType = iota
	TxActor                  // Actor that handles new transaction
	TxPoolActor              // Actor that handles consensus msg
	VerifyRspActor           // Actor that handles the response from valdiators
	NetActor                 // Actor to send msg to the net actor
	MAXACTOR
)

type SenderType uint8

const (
	NilSender  SenderType = iota
	NetSender             // Net sends tx req
	HttpSender            // Http sends tx req
)

func (sender SenderType) Sender() string {
	switch sender {
	case NilSender:
		return "nil sender"
	case NetSender:
		return "net sender"
	case HttpSender:
		return "http sender"
	default:
		return "unknown sender"
	}
}

type TxnStatsType uint8

const (
	_              TxnStatsType = iota
	RcvStats                    // The count that the tx pool receive from the actor bus
	SuccessStats                // The count that the transctions are verified successfully
	FailureStats                // The count that the transactions are invalid
	DuplicateStats              // The count that the transactions are duplicated input
	SigErrStats                 // The count that the transactions' signature error
	StateErrStats               // The count that the transactions are invalid in database

	MAXSTATS
)

type TxStatus struct {
	Hash  common.Uint256 // transaction hash
	Attrs []*TXAttr      // transaction's status
}

type TxReq struct {
	Tx     *types.Transaction
	Sender SenderType
}

type TxRsp struct {
	Hash    common.Uint256
	ErrCode errors.ErrCode
}

// restful api
type GetTxnReq struct {
	Hash common.Uint256
}

type GetTxnRsp struct {
	Txn *types.Transaction
}

type CheckTxnReq struct {
	Hash common.Uint256
}

type CheckTxnRsp struct {
	Ok bool
}

type GetTxnStatusReq struct {
	Hash common.Uint256
}

type GetTxnStatusRsp struct {
	Hash     common.Uint256
	TxStatus []*TXAttr
}

type GetTxnStats struct {
}

type GetTxnStatsRsp struct {
	Count []uint64
}

type GetPendingTxnReq struct {
	ByCount bool
}

type GetPendingTxnRsp struct {
	Txs []*types.Transaction
}

// consensus messages
type GetTxnPoolReq struct {
	ByCount bool
	Height  uint32
}

type GetTxnPoolRsp struct {
	TxnPool []*TXEntry
}

type VerifyBlockReq struct {
	Height uint32
	Txs    []*types.Transaction
}

type VerifyTxResult struct {
	Height  uint32
	Tx      *types.Transaction
	ErrCode errors.ErrCode
}

type VerifyBlockRsp struct {
	TxnPool []*VerifyTxResult
}

/*
 * Implement sort.Interface
 */
type LB struct {
	Size     int
	WorkerID uint8
}

type LBSlice []LB

func (this LBSlice) Len() int {
	return len(this)
}

func (this LBSlice) Swap(i, j int) {
	this[i].Size, this[j].Size = this[j].Size, this[i].Size
	this[i].WorkerID, this[j].WorkerID = this[j].WorkerID, this[i].WorkerID
}

func (this LBSlice) Less(i, j int) bool {
	return this[i].Size < this[j].Size
}
