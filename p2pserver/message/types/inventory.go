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
	"io"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/serialization"
	p2pCommon "github.com/ontio/dad-go/p2pserver/common"
)

var LastInvHash common.Uint256

type InvPayload struct {
	InvType common.InventoryType
	Cnt     uint32
	Blk     []byte
}

type Inv struct {
	Hdr MsgHdr
	P   InvPayload
}

func (this *InvPayload) Serialization(w io.Writer) {
	serialization.WriteUint8(w, uint8(this.InvType))
	serialization.WriteUint32(w, this.Cnt)

	binary.Write(w, binary.LittleEndian, this.Blk)
}

//Check whether header is correct
func (this Inv) Verify(buf []byte) error {
	err := this.Hdr.Verify(buf)
	return err
}

func (this Inv) invType() common.InventoryType {
	return this.P.InvType
}

//Serialize message payload
func (this Inv) Serialization() ([]byte, error) {

	tmpBuffer := bytes.NewBuffer([]byte{})
	this.P.Serialization(tmpBuffer)

	checkSumBuf := CheckSum(tmpBuffer.Bytes())
	this.Hdr.Init("inv", checkSumBuf, uint32(len(tmpBuffer.Bytes())))

	hdrBuf, err := this.Hdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	err = binary.Write(buf, binary.LittleEndian, tmpBuffer.Bytes())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

//Deserialize message payload
func (this *Inv) Deserialization(p []byte) error {
	err := this.Hdr.Deserialization(p)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(p[p2pCommon.MSG_HDR_LEN:])
	invType, err := serialization.ReadUint8(buf)
	if err != nil {
		return err
	}
	this.P.InvType = common.InventoryType(invType)
	this.P.Cnt, err = serialization.ReadUint32(buf)
	if err != nil {
		return err
	}

	this.P.Blk = make([]byte, this.P.Cnt*p2pCommon.HASH_LEN)
	err = binary.Read(buf, binary.LittleEndian, &(this.P.Blk))

	return err
}
