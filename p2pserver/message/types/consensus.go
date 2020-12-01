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

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/errors"
)

type Consensus struct {
	MsgHdr
	Cons ConsensusPayload
}

type PeerStateUpdate struct {
	PeerPubKey keypair.PublicKey
	Connected  bool
}

//Serialize message payload
func (this *Consensus) Serialization() ([]byte, error) {

	p := bytes.NewBuffer([]byte{})
	err := this.Cons.Serialize(p)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. consensus:%v", this.Cons))
	}
	checkSumBuf := CheckSum(p.Bytes())
	this.MsgHdr.Init("consensus", checkSumBuf, uint32(len(p.Bytes())))
	log.Debug("NewConsensus The message payload length is ", this.MsgHdr.Length)

	hdrBuf, err := this.MsgHdr.Serialization()
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialization error. MsgHdr:%v", this.MsgHdr))
	}
	buf := bytes.NewBuffer(hdrBuf)
	data := append(buf.Bytes(), p.Bytes()...)
	return data, nil
}

//Deserialize message payload
func (this *Consensus) Deserialization(p []byte) error {
	log.Debug()
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &(this.MsgHdr))
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read MsgHdr error. buf:%v", buf))
	}
	err = this.Cons.Deserialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("deserialize Cons error. buf:%v", buf))
	}
	return nil
}
