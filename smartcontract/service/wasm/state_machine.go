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

package wasm

import (
	"bytes"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/states"
	"github.com/ontio/dad-go/core/store"
	scommon "github.com/ontio/dad-go/core/store/common"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/smartcontract/storage"
	"github.com/ontio/dad-go/vm/wasmvm/exec"
	"github.com/ontio/dad-go/vm/wasmvm/memory"
	"github.com/ontio/dad-go/vm/wasmvm/util"
	"github.com/ontio/dad-go/vm/wasmvm/wasm"
)

type WasmStateMachine struct {
	*WasmStateReader
	ldgerStore store.LedgerStore
	CloneCache *storage.CloneCache
	time       uint32
}

func NewWasmStateMachine(ldgerStore store.LedgerStore, dbCache scommon.StateStore, time uint32) *WasmStateMachine {

	var stateMachine WasmStateMachine
	stateMachine.ldgerStore = ldgerStore
	stateMachine.CloneCache = storage.NewCloneCache(dbCache)
	stateMachine.WasmStateReader = NewWasmStateReader(ldgerStore)
	stateMachine.time = time

	stateMachine.Register("PutStorage", stateMachine.putstore)
	stateMachine.Register("GetStorage", stateMachine.getstore)
	stateMachine.Register("DeleteStorage", stateMachine.deletestore)
	stateMachine.Register("CallContract", stateMachine.callContract)

	return &stateMachine
}

//======================store apis here============================================
func (s *WasmStateMachine) putstore(engine *exec.ExecutionEngine) (bool, error) {

	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()

	if len(params) != 2 {
		return false, errors.NewErr("[putstore] parameter count error")
	}

	key, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}

	if len(key) > 1024 {
		return false, errors.NewErr("[putstore] Get Storage key to long")
	}

	value, err := vm.GetPointerMemory(params[1])
	if err != nil {
		return false, err
	}

	k, err := serializeStorageKey(vm.CodeHash, key)
	if err != nil {
		return false, err
	}

	s.CloneCache.Add(scommon.ST_STORAGE, k, &states.StorageItem{Value: value})

	vm.RestoreCtx()

	return true, nil
}

func (s *WasmStateMachine) getstore(engine *exec.ExecutionEngine) (bool, error) {

	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()

	if len(params) != 1 {
		return false, errors.NewErr("[getstore] parameter count error ")
	}

	key, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}
	k, err := serializeStorageKey(vm.CodeHash, key)
	if err != nil {
		return false, err
	}
	item, err := s.CloneCache.Get(scommon.ST_STORAGE, k)
	if err != nil {
		return false, err
	}

	if item == nil {
		vm.RestoreCtx()
		if envCall.GetReturns() {
			vm.PushResult(uint64(memory.VM_NIL_POINTER))
		}
		return true, nil
	}

	idx, err := vm.SetPointerMemory(item.(*states.StorageItem).Value)
	if err != nil {
		return false, err
	}

	vm.RestoreCtx()
	if envCall.GetReturns() {
		vm.PushResult(uint64(idx))
	}
	return true, nil
}

func (s *WasmStateMachine) deletestore(engine *exec.ExecutionEngine) (bool, error) {

	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()

	if len(params) != 1 {
		return false, errors.NewErr("[deletestore] parameter count error")
	}

	key, err := vm.GetPointerMemory(params[0])
	if err != nil {
		return false, err
	}

	k, err := serializeStorageKey(vm.CodeHash, key)
	if err != nil {
		return false, err
	}

	s.CloneCache.Delete(scommon.ST_STORAGE, k)
	vm.RestoreCtx()

	return true, nil
}

func (s *WasmStateMachine) GetContractCodeFromAddress(address common.Address) ([]byte, error) {

	dcode, err := s.ldgerStore.GetContractState(address)
	if err != nil {
		return nil, err
	}

	return dcode.Code.Code, nil

}

//call other contract
func (s *WasmStateMachine) callContract(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	envCall := vm.GetEnvCall()
	params := envCall.GetParams()
	if len(params) != 3 {
		return false, errors.NewErr("parameter count error while call readMessage")
	}
	contractAddressIdx := params[0]
	addr, err := vm.GetPointerMemory(contractAddressIdx)
	if err != nil {
		return false, errors.NewErr("get Contract address failed")
	}
	//the contract codes
	contractBytes, err := s.getContractFromAddr(addr)
	if err != nil {
		return false, err
	}
	codeHash := common.ToCodeHash(contractBytes)
	bf := bytes.NewBuffer(contractBytes)

	module, err := wasm.ReadModule(bf, emptyImporter)
	if err != nil {
		return false, errors.NewErr("load Module failed")
	}

	methodad-gome, err := vm.GetPointerMemory(params[1])
	if err != nil {
		return false, errors.NewErr("[callContract]get Contract methodad-gome failed")
	}

	arg, err := vm.GetPointerMemory(params[2])
	if err != nil {
		return false, errors.NewErr("[callContract]get Contract arg failed")
	}
	res, err := vm.CallProductContract(vm.CodeHash, codeHash, module, methodad-gome, arg)
	if err != nil {
		return false, errors.NewErr("[callContract]CallProductContract failed")
	}
	vm.RestoreCtx()
	if envCall.GetReturns() {
		vm.PushResult(uint64(res))
	}
	return true, nil
}

func serializeStorageKey(codeHash common.Address, key []byte) ([]byte, error) {
	bf := new(bytes.Buffer)
	storageKey := &states.StorageKey{CodeHash: codeHash, Key: key}
	if _, err := storageKey.Serialize(bf); err != nil {
		return []byte{}, errors.NewErr("[serializeStorageKey] StorageKey serialize error!")
	}
	return bf.Bytes(), nil
}

func (s *WasmStateMachine) getContractFromAddr(addr []byte) ([]byte, error) {

	//just for test
	/*	contract := util.TrimBuffToString(addr)
		code, err := ioutil.ReadFile(fmt.Sprintf("./testdata2/%s.wasm", contract))
		if err != nil {
			return nil, err
		}

		return code, nil*/
	//Fixme get the contract code from ledger
	addrbytes, err := common.HexToBytes(util.TrimBuffToString(addr))
	if err != nil {
		return nil, errors.NewErr("get contract address error")
	}
	contactaddress, err := common.AddressParseFromBytes(addrbytes)
	if err != nil {
		return nil, errors.NewErr("get contract address error")
	}
	dpcode, err := s.GetContractCodeFromAddress(contactaddress)
	if err != nil {
		return nil, errors.NewErr("get contract  error")
	}
	return dpcode, nil
}

func emptyImporter(name string) (*wasm.Module, error) {
	return nil, nil
}
