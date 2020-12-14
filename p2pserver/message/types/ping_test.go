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

package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"testing"

	"github.com/ontio/dad-go/common/serialization"
	"github.com/ontio/dad-go/p2pserver/common"
)

func TestPingSerializationDeserialization(t *testing.T) {
	var msg Ping
	msg.Hdr.Magic = common.NETMAGIC
	copy(msg.Hdr.CMD[0:7], "ping")
	msg.Height = 1
	t.Log("new ping message before serialize msg.Height = 1")
	tmpBuffer := bytes.NewBuffer([]byte{})
	serialization.WriteUint64(tmpBuffer, msg.Height)
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, tmpBuffer.Bytes())
	if err != nil {
		t.Error("Binary Write failed at new Msg")
		return
	}
	s := sha256.Sum256(b.Bytes())
	s2 := s[:]
	s = sha256.Sum256(s2)
	buf := bytes.NewBuffer(s[:4])
	binary.Read(buf, binary.LittleEndian, &(msg.Hdr.Checksum))
	msg.Hdr.Length = uint32(len(b.Bytes()))

	p, err := msg.Serialization()
	if err != nil {
		t.Error("Error Convert net message ", err.Error())
		return
	}
	var demsg Ping
	err = demsg.Deserialization(p)
	if err != nil {
		t.Error(err)
		return
	} else {
		t.Log("ping Test_Deserialization successful")
	}

	t.Log("deserialize ping message, msg.Height = ", demsg.Height)
}
