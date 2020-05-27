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

	"github.com/dad-go/common/serialization"
	. "github.com/dad-go/errors"
	"github.com/ontio/dad-go-crypto/keypair"
)

type ValidatorState struct {
	StateBase
	PublicKey keypair.PublicKey
}

func (this *ValidatorState) Serialize(w io.Writer) error {
	this.StateBase.Serialize(w)
	buf := keypair.SerializePublicKey(this.PublicKey)
	if err := serialization.WriteVarBytes(w, buf); err != nil {
		return err
	}
	return nil
}

func (this *ValidatorState) Deserialize(r io.Reader) error {
	if this == nil {
		this = new(ValidatorState)
	}
	err := this.StateBase.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[ValidatorState], StateBase Deserialize failed.")
	}
	buf, err := serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[ValidatorState], PublicKey Deserialize failed.")
	}
	pk, err := keypair.DeserializePublicKey(buf)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "[ValidatorState], PublicKey Deserialize failed.")
	}
	this.PublicKey = pk
	return nil
}
