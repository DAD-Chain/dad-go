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
	"encoding/binary"

	"github.com/ontio/dad-go/common/serialization"
)

type Ping struct {
	Hdr    MsgHdr
	Height uint64
}

//Check whether header is correct
func (msg Ping) Verify(buf []byte) error {
	err := msg.Hdr.Verify(buf)
	return err
}

//Serialize message payload
func (msg Ping) Serialization() ([]byte, error) {
	hdrBuf, err := msg.Hdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	err = serialization.WriteUint64(buf, msg.Height)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), err

}

//Deserialize message payload
func (msg *Ping) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &(msg.Hdr))
	if err != nil {
		return err
	}

	msg.Height, err = serialization.ReadUint64(buf)
	return err
}