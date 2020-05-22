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

package net

import (
	"github.com/dad-go/crypto"
	"github.com/dad-go/events"
	ns "github.com/dad-go/net/actor"
	"github.com/dad-go/net/node"
	"github.com/dad-go/net/protocol"
	"github.com/ontio/dad-go-eventbus/actor"
)

type Neter interface {
	//GetTxnPool(byCount bool) (map[Uint256]*types.Transaction, Fixed64)
	Xmit(interface{}) error
	GetEvent(eventName string) *events.Event
	GetBookkeepersAddrs() ([]*crypto.PubKey, uint64)
	//CleanTransactions(txns []*types.Transaction) error
	GetNeighborNoder() []protocol.Noder
	Tx(buf []byte)
	//AppendTxnPool(*types.Transaction) ErrCode
}

func SetTxnPoolPid(txnPid *actor.PID) {
	ns.SetTxnPoolPid(txnPid)
}

func SetConsensusPid(conPid *actor.PID) {
	ns.SetConsensusPid(conPid)
}

func SetLedgerPid(conPid *actor.PID) {
	ns.SetLedgerPid(conPid)
}

func InitNetServerActor(noder protocol.Noder) (*actor.PID, error) {
	netServerPid, err := ns.InitNetServer(noder)
	return netServerPid, err
}

func StartProtocol(pubKey *crypto.PubKey) protocol.Noder {
	net := node.InitNode(pubKey)
	net.ConnectSeeds()
	return net
}
