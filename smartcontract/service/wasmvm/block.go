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
	"github.com/ontio/dad-go/vm/wasmvm/exec"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/types"
	"bytes"
)

func (this *WasmVmService)blockGetCurrentHeaderHash(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()

	headerHash:= this.Store.GetCurrentHeaderHash().ToArray()
	idx,err:= vm.SetPointerMemory(headerHash)
	if err != nil{
		return false,err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true,nil
}

func (this *WasmVmService)blockGetCurrentHeaderHight(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()
	headerHight:= this.Store.GetCurrentHeaderHeight()
	vm.RestoreCtx()
	vm.PushResult(uint64(headerHight))
	return true,nil
}

func (this *WasmVmService)blockGetCurrentBlockHash(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()

	bHash:= this.Store.GetCurrentBlockHash().ToArray()
	idx,err:= vm.SetPointerMemory(bHash)
	if err != nil{
		return false,err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true,nil
}

func (this *WasmVmService)blockGetCurrentBlockHight(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	vm.RestoreCtx()
	bHight:= this.Store.GetCurrentBlockHeight()
	vm.RestoreCtx()
	vm.PushResult(uint64(bHight))
	return true,nil
}

func (this *WasmVmService)blockGetTransactionByHash(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[RuntimeLog]parameter count error ")
	}

	hashbytes,err:= vm.GetPointerMemory(params[0])
	if err != nil{
		return false,err
	}

	thash,err := common.Uint256ParseFromBytes(hashbytes)
	if err != nil{
		return false,err
	}
	tx,_,err:= this.Store.GetTransaction(thash)
	txbytes:= tx.ToArray()
	idx,err := vm.SetPointerMemory(txbytes)
	if err != nil{
		return false,err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))
	return true,nil

}
// BlockGetTransactionCount put block's transactions count to vm stack
func (this *WasmVmService)blockGetTransactionCount(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[RuntimeLog]parameter count error ")
	}

	blockbytes,err := vm.GetPointerMemory(params[0])
	if err != nil{
		return false,err
	}
	block := &types.Block{}
	err =block.Deserialize(bytes.NewBuffer(blockbytes))
	if err != nil{
		return false,err
	}

	length := len(block.Transactions)

	vm.RestoreCtx()
	vm.PushResult(uint64(length))
	return true,nil
}

// BlockGetTransactions put block's transactions to vm stack
func (this *WasmVmService)BlockGetTransactions( engine *exec.ExecutionEngine)  (bool, error)  {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 1 {
		return false, errors.NewErr("[BlockGetTransactions]parameter count error ")
	}

	blockbytes,err := vm.GetPointerMemory(params[0])
	if err != nil{
		return false,err
	}
	block := &types.Block{}
	err =block.Deserialize(bytes.NewBuffer(blockbytes))
	if err != nil{
		return false,err
	}
	transactionList := make([][]byte,len(block.Transactions))
	for i,tx := range block.Transactions {
		transactionList[i] = tx.ToArray()
	}

	idx,err := vm.SetPointerMemory(transactionList)
	if err != nil{
		return false,err
	}
	vm.RestoreCtx()
	vm.PushResult(uint64(idx))

	return true,nil
}
