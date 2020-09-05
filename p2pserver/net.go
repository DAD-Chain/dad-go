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

package p2pserver

import (
	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go-eventbus/actor"
	ns "github.com/ontio/dad-go/p2pserver/actor/req"
	"github.com/ontio/dad-go/p2pserver/node"
	"github.com/ontio/dad-go/p2pserver/protocol"
)

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
	//netServerPid, err := ns.InitNetServer(noder)
	return nil, nil
}

func StartProtocol(pubKey keypair.PublicKey) protocol.Noder {
	net := node.InitNode(pubKey)
	net.ConnectSeeds()
	return net
}
