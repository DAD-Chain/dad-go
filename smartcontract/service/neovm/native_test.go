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

package neovm

import (
	"bytes"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/vm/neovm/types"
	"github.com/stretchr/testify/assert"
	"math"
	"math/big"
	"testing"
)

func TestBuildParamToNative(t *testing.T) {
	inte := types.NewInteger(new(big.Int).SetUint64(math.MaxUint64))
	boo := types.NewBoolean(false)
	bs := types.NewByteArray([]byte("hello"))
	s := make([]types.StackItems, 0)
	s = append(s, inte)
	s = append(s, boo)
	s = append(s, bs)
	stru := types.NewStruct(s)
	arr := types.NewArray(nil)
	arr.Add(stru)

	buff := new(bytes.Buffer)
	err := BuildParamToNative(buff, arr)
	assert.Nil(t, err)
	assert.Equal(t, "010109ffffffffffffffff00000568656c6c6f", common.ToHexString(buff.Bytes()))
}
