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

package states

import (
	"io"
	"math/big"

	"github.com/dad-go/common"
	"github.com/dad-go/errors"
	"github.com/dad-go/common/serialization"
)

type TransferFrom struct {
	Version byte
	Sender common.Address
	From common.Address
	To common.Address
	Value *big.Int
}

func (this *TransferFrom) Serialize(w io.Writer) error {
	if err := serialization.WriteByte(w, byte(this.Version)); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Serialize version error!")
	}
	if err := this.Sender.Serialize(w); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Serialize sender error!")
	}
	if err := this.From.Serialize(w); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Serialize from error!")
	}
	if err := this.To.Serialize(w); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Serialize to error!")
	}
	if err := serialization.WriteVarBytes(w, this.Value.Bytes()); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Serialize value error!")
	}
	return nil
}

func (this *TransferFrom) Deserialize(r io.Reader) error {
	version, err := serialization.ReadByte(r); if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Deserialize version error!")
	}
	this.Version = version

	sender := new(common.Address)
	if err := sender.Deserialize(r); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Deserialize sender error!")
	}
	this.Sender = *sender

	from := new(common.Address)
	if err := from.Deserialize(r); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Deserialize from error!")
	}
	this.From = *from

	to := new(common.Address)
	if err := to.Deserialize(r); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Deserialize to error!")
	}
	this.To = *to

	value, err := serialization.ReadVarBytes(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TransferFrom] Deserialize value error!")
	}

	this.Value = new(big.Int).SetBytes(value)
	return nil
}
