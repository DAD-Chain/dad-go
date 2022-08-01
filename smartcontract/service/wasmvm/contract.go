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
	"github.com/go-interpreter/wagon/exec"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/payload"
	"github.com/ontio/dad-go/errors"
)

func ContractCreate(proc *exec.Process,
	codePtr uint32,
	codeLen uint32,
	needStorage uint32,
	namePtr uint32,
	nameLen uint32,
	verPtr uint32,
	verLen uint32,
	authorPtr uint32,
	authorLen uint32,
	emailPtr uint32,
	emailLen uint32,
	descPtr uint32,
	descLen uint32,
	newAddressPtr uint32) uint32 {
	self := proc.HostData().(*Runtime)
	code, err := ReadWasmMemory(proc, codePtr, codeLen)
	if err != nil {
		panic(err)
	}

	cost := CONTRACT_CREATE_GAS + uint64(uint64(codeLen)/PER_UNIT_CODE_LEN)*UINT_DEPLOY_CODE_LEN_GAS
	self.checkGas(cost)

	name, err := ReadWasmMemory(proc, namePtr, nameLen)
	if err != nil {
		panic(err)
	}

	version, err := ReadWasmMemory(proc, verPtr, verLen)
	if err != nil {
		panic(err)
	}

	author, err := ReadWasmMemory(proc, authorPtr, authorLen)
	if err != nil {
		panic(err)
	}

	email, err := ReadWasmMemory(proc, emailPtr, emailLen)
	if err != nil {
		panic(err)
	}

	desc, err := ReadWasmMemory(proc, descPtr, descLen)
	if err != nil {
		panic(err)
	}

	dep, err := payload.CreateDeployCode(code, needStorage, name, version, author, email, desc)
	if err != nil {
		panic(err)
	}

	if dep.VmType() != payload.WASMVM_TYPE {
		panic("[Contract] expect WASMVM_TYPE. get NEOVM_TYPE")
	}

	_, err = ReadWasmModule(dep.Code, true)
	if err != nil {
		panic(err)
	}

	contractAddr := dep.Address()
	if self.isContractExist(contractAddr) {
		panic(errors.NewErr("contract has been deployed"))
	}

	self.Service.CacheDB.PutContract(dep)

	length, err := proc.WriteAt(contractAddr[:], int64(newAddressPtr))
	return uint32(length)

}

func ContractMigrate(proc *exec.Process,
	codePtr uint32,
	codeLen uint32,
	needStorage uint32,
	namePtr uint32,
	nameLen uint32,
	verPtr uint32,
	verLen uint32,
	authorPtr uint32,
	authorLen uint32,
	emailPtr uint32,
	emailLen uint32,
	descPtr uint32,
	descLen uint32,
	newAddressPtr uint32) uint32 {

	self := proc.HostData().(*Runtime)

	code, err := ReadWasmMemory(proc, codePtr, codeLen)
	if err != nil {
		panic(err)
	}

	cost := CONTRACT_CREATE_GAS + uint64(uint64(codeLen)/PER_UNIT_CODE_LEN)*UINT_DEPLOY_CODE_LEN_GAS
	self.checkGas(cost)

	name, err := ReadWasmMemory(proc, namePtr, nameLen)
	if err != nil {
		panic(err)
	}

	version, err := ReadWasmMemory(proc, verPtr, verLen)
	if err != nil {
		panic(err)
	}

	author, err := ReadWasmMemory(proc, authorPtr, authorLen)
	if err != nil {
		panic(err)
	}

	email, err := ReadWasmMemory(proc, emailPtr, emailLen)
	if err != nil {
		panic(err)
	}

	desc, err := ReadWasmMemory(proc, descPtr, descLen)
	if err != nil {
		panic(err)
	}

	dep, err := payload.CreateDeployCode(code, needStorage, name, version, author, email, desc)
	if err != nil {
		panic(err)
	}

	if dep.VmType() != payload.WASMVM_TYPE {
		panic("[Contract] expect WASMVM_TYPE. get NEOVM_TYPE")
	}

	_, err = ReadWasmModule(dep.Code, true)
	if err != nil {
		panic(err)
	}

	contractAddr := dep.Address()
	if self.isContractExist(contractAddr) {
		panic(errors.NewErr("contract has been deployed"))
	}
	oldAddress := self.Service.ContextRef.CurrentContext().ContractAddress

	self.Service.CacheDB.PutContract(dep)
	self.Service.CacheDB.DeleteContract(oldAddress)

	iter := self.Service.CacheDB.NewIterator(oldAddress[:])
	for has := iter.First(); has; has = iter.Next() {
		key := iter.Key()
		val := iter.Value()

		newkey := serializeStorageKey(contractAddr, key[20:])

		self.Service.CacheDB.Put(newkey, val)
		self.Service.CacheDB.Delete(key)
	}

	iter.Release()
	if err := iter.Error(); err != nil {
		panic(err)
	}

	length, err := proc.WriteAt(contractAddr[:], int64(newAddressPtr))
	if err != nil {
		panic(err)
	}

	return uint32(length)
}

func ContractDestroy(proc *exec.Process) {
	self := proc.HostData().(*Runtime)
	contractAddress := self.Service.ContextRef.CurrentContext().ContractAddress
	iter := self.Service.CacheDB.NewIterator(contractAddress[:])

	for has := iter.First(); has; has = iter.Next() {
		self.Service.CacheDB.Delete(iter.Key())
	}
	iter.Release()
	if err := iter.Error(); err != nil {
		panic(err)
	}

	//the contract has been deleted ,quit the contract operation
	proc.Terminate()
}

func (self *Runtime) isContractExist(contractAddress common.Address) bool {
	item, err := self.Service.CacheDB.GetContract(contractAddress)
	if err != nil {
		panic(err)
	}
	return item != nil
}
