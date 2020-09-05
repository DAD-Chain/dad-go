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

package link

import (
	"net"
	"time"

	"github.com/ontio/dad-go/p2pserver/common"
)

// The RX buffer of this node to solve mutliple packets problem
type RxBuf struct {
	p   []byte //buffer
	len int    //patload length in buffer
}

//Link used to establish
type Link struct {
	id       uint64
	addr     string    // The address of the node
	conn     net.Conn  // Connect socket with the peer node
	port     uint16    // The server port of the node
	rxBuf    RxBuf     // recv buffer
	time     time.Time // The latest time the node activity
	recvChan chan common.MsgPayload
}
