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

package dbft

import (
	"github.com/dad-go/common/log"
	ser "github.com/dad-go/common/serialization"
	"github.com/dad-go/core/types"
	. "github.com/dad-go/errors"
	"io"
	"github.com/dad-go/common"
)

type PrepareRequest struct {
	msgData        ConsensusMessageData
	Nonce          uint64
	NextBookKeeper common.Address
	Transactions   []*types.Transaction
	Signature      []byte
}

func (pr *PrepareRequest) Serialize(w io.Writer) error {
	log.Debug()

	pr.msgData.Serialize(w)
	if err := ser.WriteVarUint(w, pr.Nonce); err != nil {
		return NewDetailErr(err, ErrNoCode, "[PrepareRequest] nonce serialization failed")
	}
	if err := pr.NextBookKeeper.Serialize(w); err != nil {
		return NewDetailErr(err, ErrNoCode, "[PrepareRequest] nextbookKeeper serialization failed")
	}
	if err := ser.WriteVarUint(w, uint64(len(pr.Transactions))); err != nil {
		return NewDetailErr(err, ErrNoCode, "[PrepareRequest] length serialization failed")
	}
	for _, t := range pr.Transactions {
		if err := t.Serialize(w); err != nil {
			return NewDetailErr(err, ErrNoCode, "[PrepareRequest] transactions serialization failed")
		}
	}
	if err := ser.WriteVarBytes(w, pr.Signature); err != nil {
		return NewDetailErr(err, ErrNoCode, "[PrepareRequest] signature serialization failed")
	}
	return nil
}

func (pr *PrepareRequest) Deserialize(r io.Reader) error {
	pr.msgData = ConsensusMessageData{}
	pr.msgData.Deserialize(r)
	pr.Nonce, _ = ser.ReadVarUint(r, 0)

	if err := pr.NextBookKeeper.Deserialize(r); err != nil {
		return NewDetailErr(err, ErrNoCode, "[PrepareRequest] nextbookKeeper deserialization failed")
	}

	length, err := ser.ReadVarUint(r, 0)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[PrepareRequest] length deserialization failed")
	}

	pr.Transactions = make([]*types.Transaction, length)
	for i := 0; i < len(pr.Transactions); i++ {
		var t types.Transaction
		if err := t.Deserialize(r); err != nil {
			return NewDetailErr(err, ErrNoCode, "[PrepareRequest] transactions deserialization failed")
		}
		pr.Transactions[i] = &t
	}

	pr.Signature, err = ser.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[PrepareRequest] signature deserialization failed")
	}

	return nil
}

func (pr *PrepareRequest) Type() ConsensusMessageType {
	log.Debug()
	return pr.ConsensusMessageData().Type
}

func (pr *PrepareRequest) ViewNumber() byte {
	log.Debug()
	return pr.msgData.ViewNumber
}

func (pr *PrepareRequest) ConsensusMessageData() *ConsensusMessageData {
	log.Debug()
	return &(pr.msgData)
}
