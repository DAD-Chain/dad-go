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
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/serialization"
	"github.com/ontio/dad-go/smartcontract/service/native"
	"github.com/ontio/dad-go/smartcontract/states"
	vm "github.com/ontio/dad-go/vm/neovm"
	"github.com/ontio/dad-go/vm/neovm/types"
)

func NativeInvoke(service *NeoVmService, engine *vm.ExecutionEngine) error {
	count := vm.EvaluationStackCount(engine)
	if count < 4 {
		return fmt.Errorf("invoke native contract invalid parameters %d < 4 ", count)
	}
	version, err := vm.PopInt(engine)
	if err != nil {
		return err
	}
	address, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	addr, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("invoke native contract:%s, address invalid", address)
	}
	method, err := vm.PopByteArray(engine)
	if err != nil {
		return err
	}
	if len(method) > METHOD_LENGTH_LIMIT {
		return fmt.Errorf("invoke native contract:%s method:%s too long, over max length 1024 limit", address, method)
	}
	args := vm.PopStackItem(engine)

	buf := new(bytes.Buffer)
	if err := BuildParamToNative(buf, args); err != nil {
		return err
	}

	contract := states.ContractInvokeParam{
		Version: byte(version),
		Address: addr,
		Method:  string(method),
		Args:    buf.Bytes(),
	}

	native := &native.NativeService{
		CacheDB:     service.CacheDB,
		InvokeParam: contract,
		Tx:          service.Tx,
		Height:      service.Height,
		Time:        service.Time,
		ContextRef:  service.ContextRef,
		ServiceMap:  make(map[string]native.Handler),
	}

	result, err := native.Invoke()
	if err != nil {
		return err
	}
	vm.PushData(engine, result)
	return nil
}

func BuildParamToNative(bf *bytes.Buffer, item types.StackItems) error {
	if CircularRefAndDepthDetection(item) {
		return errors.New("invoke native circular reference!")
	}
	return buildParamToNative(bf, item)
}

func buildParamToNative(bf *bytes.Buffer, item types.StackItems) error {
	switch item.(type) {
	case *types.ByteArray:
		a, _ := item.GetByteArray()
		if err := serialization.WriteVarBytes(bf, a); err != nil {
			return err
		}
	case *types.Integer:
		i, _ := item.GetByteArray()
		if err := serialization.WriteVarBytes(bf, i); err != nil {
			return err
		}
	case *types.Boolean:
		b, _ := item.GetBoolean()
		if err := serialization.WriteBool(bf, b); err != nil {
			return err
		}
	case *types.Array:
		arr, _ := item.GetArray()
		if err := serialization.WriteVarBytes(bf, types.BigIntToBytes(big.NewInt(int64(len(arr))))); err != nil {
			return err
		}
		for _, v := range arr {
			if err := buildParamToNative(bf, v); err != nil {
				return err
			}
		}
	case *types.Struct:
		st, _ := item.GetStruct()
		for _, v := range st {
			if err := buildParamToNative(bf, v); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("convert neovm params to native invalid type support: %s", reflect.TypeOf(item))
	}
	return nil
}
