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
package wasmvm

import (
	"bytes"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/vm/wasmvm/exec"
)

func (this *WasmVmService) transactionGetHash(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[transactionGetHash] parameter count error")
	}

	transbytes, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}

	trans := types.Transaction{}
	err = trans.Deserialize(bytes.NewBuffer(transbytes))
	if err != nil {
		return false, err
	}
	hash := trans.Hash()
	idx, err := vm.SetPointerMemory(hash.ToArray())
	if err != nil {
		return false, err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true, nil
}
func (this *WasmVmService) transactionGetType(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[transactionGetType] parameter count error")
	}

	transbytes, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}

	trans := types.Transaction{}
	err = trans.Deserialize(bytes.NewBuffer(transbytes))
	if err != nil {
		return false, err
	}
	txtype := trans.TxType

	vm.RestoreCtx()
	vm.PushResult(uint64(txtype))
	return true, nil
}
func (this *WasmVmService) transactionGetAttributes(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[transactionGetAttributes] parameter count error")
	}

	transbytes, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}

	trans := types.Transaction{}
	err = trans.Deserialize(bytes.NewBuffer(transbytes))
	if err != nil {
		return false, err
	}
	attributes := make([][]byte, 0)

	idx, err := vm.SetPointerMemory(attributes)
	if err != nil {
		return false, err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true, nil
}
