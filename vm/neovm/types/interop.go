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
	"math/big"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/vm/neovm/interfaces"
)

type Interop struct {
	_object interfaces.Interop
}

func NewInteropInterface(value interfaces.Interop) *Interop {
	var ii Interop
	ii._object = value
	return &ii
}

func (this *Interop) Equals(other StackItems) bool {
	if _, ok := other.(*Interop); !ok {
		return false
	}
	if !common.IsEqualBytes(this._object.ToArray(), other.GetInterface().ToArray()) {
		return false
	}
	return true
}

func (this *Interop) GetBigInteger() *big.Int {
	return big.NewInt(0)
}

func (this *Interop) GetBoolean() bool {
	if this._object == nil {
		return false
	}
	return true
}

func (this *Interop) GetByteArray() []byte {
	return this._object.ToArray()
}

func (this *Interop) GetInterface() interfaces.Interop {
	return this._object
}

func (this *Interop) GetArray() []StackItems {
	return []StackItems{this}
}

func (this *Interop) GetStruct() []StackItems {
	return []StackItems{this}
}
