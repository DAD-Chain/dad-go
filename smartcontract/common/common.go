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
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/vm/neovm/types"
	"github.com/ontio/dad-go/common/log"
)

// ConvertReturnTypes return neovm stack element value
// According item types convert to hex string value
// Now neovm support type contain: ByteArray/Integer/Boolean/Array/Struct/Interop/StackItems
func ConvertNeoVmTypeHexString(item interface{}) interface{} {
	if item == nil {
		return nil
	}
	switch v := item.(type) {
	case *types.ByteArray:
		return common.ToHexString(v.GetByteArray())
	case *types.Integer:
		if v.GetBigInteger().Sign() == 0 {
			return common.ToHexString([]byte{0})
		} else {
			return common.ToHexString(types.ConvertBigIntegerToBytes(v.GetBigInteger()))
		}
	case *types.Boolean:
		if v.GetBoolean() {
			return common.ToHexString([]byte{1})
		} else {
			return common.ToHexString([]byte{0})
		}
	case *types.Array:
		var arr []interface{}
		for _, val := range v.GetArray() {
			arr = append(arr, ConvertNeoVmTypeHexString(val))
		}
		return arr
	case *types.Struct:
		var arr []interface{}
		for _, val := range v.GetStruct() {
			arr = append(arr, ConvertNeoVmTypeHexString(val))
		}
		return arr
	case *types.Interop:
		return common.ToHexString(v.GetInterface().ToArray())
	default:
		log.Error("[ConvertTypes] Invalid Types!")
		return nil
	}
}

func ConvertNeoVmReturnTypes(item interface{}) interface{} {
	if item == nil {
		return nil
	}
	switch v := item.(type) {
	case *types.ByteArray:
		return v.GetByteArray()
	case *types.Integer:
		return types.ConvertBigIntegerToBytes(v.GetBigInteger())
	case *types.Boolean:
		if v.GetBoolean() {
			return []byte{1}
		}else {
			return []byte{0}
		}
	case *types.Array:
		var arr []interface{}
		for _, val := range v.GetArray() {
			arr = append(arr, ConvertNeoVmReturnTypes(val))
		}
		return arr
	case *types.Struct:
		var arr []interface{}
		for _, val := range v.GetStruct() {
			arr = append(arr, ConvertNeoVmReturnTypes(val))
		}
		return arr
	case *types.Interop:
		return v.GetInterface().ToArray()
	default:
		log.Error("[ConvertTypes] Invalid Types!")
		return nil
	}
}


