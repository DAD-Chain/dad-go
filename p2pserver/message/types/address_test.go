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
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */

package types

import (
	"bytes"
	"net"
	"testing"

	comm "github.com/ontio/dad-go/p2pserver/common"
	"github.com/stretchr/testify/assert"
)

func MessageTest(t *testing.T, msg Message) {
	p := new(bytes.Buffer)
	err := WriteMessage(p, msg)
	assert.Nil(t, err)

	demsg, err := ReadMessage(p)
	assert.Nil(t, err)

	assert.Equal(t, msg, demsg)
}

func TestAddressSerializationDeserialization(t *testing.T) {
	var msg Addr
	var addr [16]byte
	ip := net.ParseIP("192.168.0.1")
	ip.To16()
	copy(addr[:], ip[:16])
	nodeAddr := comm.PeerAddr{
		Time:          12345678,
		Services:      100,
		IpAddr:        addr,
		Port:          8080,
		ConsensusPort: 8081,
		ID:            987654321,
	}
	msg.NodeAddrs = append(msg.NodeAddrs, nodeAddr)

	MessageTest(t, &msg)
}
