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

package netserver

import (
	"github.com/ontio/dad-go-crypto/keypair"

	"github.com/ontio/dad-go/p2pserver/net/protocol"
	"github.com/ontio/dad-go/p2pserver/peer"
)

//NewNetServer return the net object in p2p
func NewNetServer(pubKey keypair.PublicKey) p2p.P2P {

	return nil
}

//NetServer represent all the actions in net layer
type NetServer struct {
	base peer.PeerCom
}

//init initializes attribute of network server
func (n *NetServer) init(pubKey keypair.PublicKey) error {

	return nil
}
