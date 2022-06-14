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
	"fmt"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/payload"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/smartcontract/service/util"
	"github.com/ontio/dad-go/vm/crossvm_codec"
	vm "github.com/ontio/dad-go/vm/neovm"
)

//neovm contract call wasmvm contract
func WASMInvoke(service *NeoVmService, engine *vm.Executor) error {

	address, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}

	contractAddress, err := common.AddressParseFromBytes(address)
	if err != nil {
		return fmt.Errorf("invoke wasm contract:%s, address invalid", address)
	}

	dp, err := service.CacheDB.GetContract(contractAddress)
	if err != nil {
		return err
	}
	if dp == nil {
		return fmt.Errorf("wasm contract does not exist")
	}
	if dp.VmType != payload.WASMVM_TYPE {
		return fmt.Errorf("not a wasm contract")
	}

	parambytes, err := engine.EvalStack.PopAsBytes()
	if err != nil {
		return err
	}
	list, err := crossvm_codec.DeserializeInput(parambytes)
	if err != nil {
		return err
	}

	inputs, err := util.BuildWasmVMInvokeCode(contractAddress, list)
	if err != nil {
		return err
	}

	newservice, err := service.ContextRef.NewExecuteEngine(inputs, types.InvokeWasm)
	if err != nil {
		return err
	}

	tmpRes, err := newservice.Invoke()
	if err != nil {
		return err
	}

	engine.EvalStack.PushBytes(tmpRes.([]byte))
	return nil

}
