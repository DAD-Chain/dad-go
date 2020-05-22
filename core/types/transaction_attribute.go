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
	"errors"
	"io"

	"github.com/dad-go/common/serialization"
	. "github.com/dad-go/errors"
)

type TransactionAttributeUsage byte

const (
	Nonce          TransactionAttributeUsage = 0x00
	Script         TransactionAttributeUsage = 0x20
	DescriptionUrl TransactionAttributeUsage = 0x81
	Description    TransactionAttributeUsage = 0x90
)

func IsValidAttributeType(usage TransactionAttributeUsage) bool {
	return usage == Nonce || usage == Script ||
		usage == DescriptionUrl || usage == Description
}

type TxAttribute struct {
	Usage TransactionAttributeUsage
	Data  []byte
	Size  uint32
}

func NewTxAttribute(u TransactionAttributeUsage, d []byte) TxAttribute {
	tx := TxAttribute{u, d, 0}
	tx.Size = tx.GetSize()
	return tx
}

func (u *TxAttribute) GetSize() uint32 {
	if u.Usage == DescriptionUrl {
		return uint32(len([]byte{(byte(0xff))}) + len([]byte{(byte(0xff))}) + len(u.Data))
	}
	return 0
}

func (tx *TxAttribute) Serialize(w io.Writer) error {
	if err := serialization.WriteUint8(w, byte(tx.Usage)); err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction attribute Usage serialization error.")
	}
	if !IsValidAttributeType(tx.Usage) {
		return NewDetailErr(errors.New("[TxAttribute] error"), ErrNoCode, "Unsupported attribute Description.")
	}
	if err := serialization.WriteVarBytes(w, tx.Data); err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction attribute Data serialization error.")
	}
	return nil
}

func (tx *TxAttribute) Deserialize(r io.Reader) error {
	val, err := serialization.ReadBytes(r, 1)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction attribute Usage deserialization error.")
	}
	tx.Usage = TransactionAttributeUsage(val[0])
	if !IsValidAttributeType(tx.Usage) {
		return NewDetailErr(errors.New("[TxAttribute] error"), ErrNoCode, "Unsupported attribute Description.")
	}
	tx.Data, err = serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction attribute Data deserialization error.")
	}
	return nil

}

func (tx *TxAttribute) ToArray() []byte {
	bf := new(bytes.Buffer)
	tx.Serialize(bf)
	return bf.Bytes()
}
