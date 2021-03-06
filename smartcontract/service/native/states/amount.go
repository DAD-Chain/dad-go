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
	"math/big"
	"io"
	"github.com/dad-go/common/serialization"
	"github.com/dad-go/errors"
)

type Amount struct {
	Value *big.Int
}

func(this *Amount) Serialize(w io.Writer) error {
	return serialization.WriteVarBytes(w, this.Value.Bytes())
}

func(this *Amount) Deserialize(r io.Reader) error {
	if this == nil {
		this = new(Amount)
	}
	bs, err := serialization.ReadVarBytes(r)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[TotalSupply Deserialize] read value error!")
	}
	this.Value = new(big.Int).SetBytes(bs)
	return nil
}
