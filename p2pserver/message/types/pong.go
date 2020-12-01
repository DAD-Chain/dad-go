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
	"fmt"

	"github.com/ontio/dad-go/common/serialization"
	"github.com/ontio/dad-go/errors"
)

type Pong struct {
	MsgHdr
	Height uint64
}

//Check whether header is correct
func (this Pong) Verify(buf []byte) error {
	err := this.MsgHdr.Verify(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetVerifyFail, fmt.Sprintf("verify error. buf:%v", buf))
	}
	return nil
}

//Serialize message payload
func (this Pong) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := serialization.WriteUint64(p, this.Height)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. Height:%v", this.Height))
	}
	checkSumBuf := CheckSum(p.Bytes())
	this.MsgHdr.Init("pong", checkSumBuf, uint32(len(p.Bytes())))

	hdrBuf, err := this.MsgHdr.Serialization()
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialization error. MsgHdr:%v", this.MsgHdr))
	}
	buf := bytes.NewBuffer(hdrBuf)
	data := append(buf.Bytes(), p.Bytes()...)
	return data, nil

}

//Deserialize message payload
func (this *Pong) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &(this.MsgHdr))
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read Hdr error. buf:%v", buf))
	}

	this.Height, err = serialization.ReadUint64(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read Height error. buf:%v", buf))
	}
	return nil
}
