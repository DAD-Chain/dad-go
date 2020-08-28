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
	"errors"
	"time"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go-eventbus/actor"
	"github.com/ontio/dad-go/core/types"
	ontErrors "github.com/ontio/dad-go/errors"
	netActor "github.com/ontio/dad-go/net/actor"
	txpool "github.com/ontio/dad-go/txnpool/common"
)

type TxPoolActor struct {
	Pool *actor.PID
}

func (self *TxPoolActor) GetTxnPool(byCount bool, height uint32) []*txpool.TXEntry {
	poolmsg := &txpool.GetTxnPoolReq{ByCount: byCount, Height: height}
	future := self.Pool.RequestFuture(poolmsg, time.Second*10)
	entry, err := future.Result()
	if err != nil {
		return nil
	}

	txs := entry.(*txpool.GetTxnPoolRsp).TxnPool
	return txs
}

func (self *TxPoolActor) VerifyBlock(txs []*types.Transaction, height uint32) error {
	poolmsg := &txpool.VerifyBlockReq{Txs: txs, Height: height}
	future := self.Pool.RequestFuture(poolmsg, time.Second*10)
	entry, err := future.Result()
	if err != nil {
		return err
	}

	txentry := entry.(*txpool.VerifyBlockRsp).TxnPool
	for _, entry := range txentry {
		if entry.ErrCode != ontErrors.ErrNoError {
			return errors.New(entry.ErrCode.Error())
		}
	}

	return nil
}

type P2PActor struct {
	P2P *actor.PID
}

func (self *P2PActor) Broadcast(msg interface{}) {
	self.P2P.Tell(msg)
}

func (self *P2PActor) Transmit(target *keypair.PublicKey, msg []byte) {
	self.P2P.Tell(&netActor.TransmitConsensusMsgReq{
		Target: target,
		Msg:    msg,
	})
}

type LedgerActor struct {
	Ledger *actor.PID
}
