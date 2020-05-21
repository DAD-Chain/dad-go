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

package common

import (
	"testing"
	"github.com/dad-go/vm/neovm/types"
	"math/big"
)

func TestConvertTypes(t *testing.T) {
	arr := types.NewArray([]types.StackItemInterface{types.NewByteArray([]byte{1,2,3}), types.NewInteger(big.NewInt(32))})
	var states []States
	for _, v := range arr.GetArray() {
		states = append(states, ConvertTypes(v)...)
	}
	t.Log("result:", states)
}

func TestConvertReturnTypes(t *testing.T) {
	arr := types.NewArray([]types.StackItemInterface{types.NewByteArray([]byte{1,2,3}), types.NewInteger(big.NewInt(32)), types.NewArray([]types.StackItemInterface{types.NewByteArray([]byte{1,2,3}), types.NewInteger(big.NewInt(32))})})
	var states []interface{}
	for _, v := range arr.GetArray() {
		states = append(states, ConvertReturnTypes(v)...)
	}
	t.Log("result:", states)
}
