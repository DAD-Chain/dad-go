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

package req

import (
	"time"

	"github.com/ontio/dad-go-eventbus/actor"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/errors"
	p2pcommon "github.com/ontio/dad-go/p2pserver/common"
	txnpool "github.com/ontio/dad-go/txnpool/common"
)

var TxnPoolPid *actor.PID

func SetTxnPoolPid(txnPid *actor.PID) {
	TxnPoolPid = txnPid
}

//add txn to txnpool
func AddTransaction(transaction *types.Transaction) {
	TxnPoolPid.Tell(transaction)
}

//get all txns
func GetTxnPool(byCount bool) ([]*txnpool.TXEntry, error) {
	future := TxnPoolPid.RequestFuture(&txnpool.GetTxnPoolReq{ByCount: byCount}, p2pcommon.ACTOR_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("p2p GetTxnPool ERROR: "), err)
		return nil, err
	}
	return result.(txnpool.GetTxnPoolRsp).TxnPool, nil
}

//get txn according to hash
func GetTransaction(hash common.Uint256) (*types.Transaction, error) {
	future := TxnPoolPid.RequestFuture(&txnpool.GetTxnReq{Hash: hash}, p2pcommon.ACTOR_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("p2p GetTransaction ERROR: "), err)
		return nil, err
	}
	return result.(txnpool.GetTxnRsp).Txn, nil
}

//check whether txn in txnpool
func CheckTransaction(hash common.Uint256) (bool, error) {
	future := TxnPoolPid.RequestFuture(&txnpool.CheckTxnReq{Hash: hash}, p2pcommon.ACTOR_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("p2p CheckTransaction ERROR: "), err)
		return false, err
	}
	return result.(txnpool.CheckTxnRsp).Ok, nil
}

//get tx status according to hash
func GetTransactionStatus(hash common.Uint256) ([]*txnpool.TXAttr, error) {
	future := TxnPoolPid.RequestFuture(&txnpool.GetTxnStatusReq{Hash: hash}, p2pcommon.ACTOR_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("p2p GetTransactionStatus ERROR: "), err)
		return nil, err
	}
	return result.(txnpool.GetTxnStatusRsp).TxStatus, nil
}

//get pending txn by count
func GetPendingTxn(byCount bool) ([]*types.Transaction, error) {
	future := TxnPoolPid.RequestFuture(&txnpool.GetPendingTxnReq{ByCount: byCount}, p2pcommon.ACTOR_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("p2p GetPendingTxn ERROR: "), err)
		return nil, err
	}
	return result.(txnpool.GetPendingTxnRsp).Txs, nil
}

//get veritfy block result from txnpool
func VerifyBlock(height uint32, txs []*types.Transaction) ([]*txnpool.VerifyTxResult, error) {
	future := TxnPoolPid.RequestFuture(&txnpool.VerifyBlockReq{Height: height, Txs: txs}, p2pcommon.ACTOR_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("p2p VerifyBlock ERROR: "), err)
		return nil, err
	}
	return result.(txnpool.VerifyBlockRsp).TxnPool, nil
}

//get txn stats according to hash
func GetTransactionStats(hash common.Uint256) ([]uint64, error) {
	future := TxnPoolPid.RequestFuture(&txnpool.GetTxnStats{}, p2pcommon.ACTOR_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("p2p GetTransactionStats ERROR: "), err)
		return nil, err
	}
	return result.(txnpool.GetTxnStatsRsp).Count, nil
}
