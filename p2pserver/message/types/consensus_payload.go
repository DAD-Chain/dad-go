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
	"fmt"
	"io"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/signature"
	"github.com/ontio/dad-go/errors"
)

type ConsensusPayload struct {
	Version         uint32
	PrevHash        common.Uint256
	Height          uint32
	BookkeeperIndex uint16
	Timestamp       uint32
	Data            []byte
	Owner           keypair.PublicKey
	Signature       []byte
	PeerId          uint64
	hash            common.Uint256
}

//get the consensus payload hash
func (this *ConsensusPayload) Hash() common.Uint256 {
	return common.Uint256{}
}

//Check whether header is correct
func (this *ConsensusPayload) Verify() error {
	sink := common.NewZeroCopySink(nil)
	this.SerializationUnsigned(sink)

	err := signature.Verify(this.Owner, sink.Bytes(), this.Signature)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetVerifyFail, fmt.Sprintf("signature verify error. buf:%v", sink.Bytes()))
	}
	return nil
}

//serialize the consensus payload
func (this *ConsensusPayload) ToArray() []byte {
	return common.SerializeToBytes(this)
}

//return inventory type
func (this *ConsensusPayload) InventoryType() common.InventoryType {
	return common.CONSENSUS
}

func (this *ConsensusPayload) GetMessage() []byte {
	//TODO: GetMessage
	//return sig.GetHashData(cp)
	return []byte{}
}

func (this *ConsensusPayload) Type() common.InventoryType {

	//TODO:Temporary add for Interface signature.SignableData use.
	return common.CONSENSUS
}

func (this *ConsensusPayload) Serialization(sink *common.ZeroCopySink) {
	this.SerializationUnsigned(sink)
	buf := keypair.SerializePublicKey(this.Owner)
	sink.WriteVarBytes(buf)
	sink.WriteVarBytes(this.Signature)
}

//Deserialize message payload
func (this *ConsensusPayload) Deserialization(source *common.ZeroCopySource) error {
	err := this.DeserializationUnsigned(source)
	if err != nil {
		return err
	}
	buf, _, irregular, eof := source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}

	this.Owner, err = keypair.DeserializePublicKey(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, "deserialize publickey error")
	}

	this.Signature, _, irregular, eof = source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}

	return nil
}

func (this *ConsensusPayload) SerializationUnsigned(sink *common.ZeroCopySink) {
	sink.WriteUint32(this.Version)
	sink.WriteHash(this.PrevHash)
	sink.WriteUint32(this.Height)
	sink.WriteUint16(this.BookkeeperIndex)
	sink.WriteUint32(this.Timestamp)
	sink.WriteVarBytes(this.Data)
}

func (this *ConsensusPayload) DeserializationUnsigned(source *common.ZeroCopySource) error {
	var irregular, eof bool
	this.Version, eof = source.NextUint32()
	this.PrevHash, eof = source.NextHash()
	this.Height, eof = source.NextUint32()
	this.BookkeeperIndex, eof = source.NextUint16()
	this.Timestamp, eof = source.NextUint32()
	this.Data, _, irregular, eof = source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}

	return nil
}
