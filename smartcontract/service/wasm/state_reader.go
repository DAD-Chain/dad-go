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
	"errors"
	"github.com/dad-go/smartcontract/event"
	trigger "github.com/dad-go/smartcontract/types"

	"github.com/dad-go/vm/wasmvm/exec"
	"github.com/dad-go/core/store"
)



type WasmStateReader struct{
	serviceMap map[string]func(*exec.ExecutionEngine) (bool, error)
	trigger    trigger.TriggerType
	Notifications []*event.NotifyEventInfo
	ldgerStore    store.ILedgerStore
}

func NewWasmStateReader(ldgerStore store.ILedgerStore,trigger trigger.TriggerType) *WasmStateReader {

	i := &WasmStateReader{
		ldgerStore:ldgerStore,
		serviceMap: make(map[string]func(*exec.ExecutionEngine) (bool, error)),
		trigger:trigger,
		}
	return i
}

func (i *WasmStateReader) Register(name string, handler func(*exec.ExecutionEngine) (bool, error)) bool {
	if _, ok := i.serviceMap[name]; ok {
		return false
	}
	i.serviceMap[name] = handler
	return true
}

func (i *WasmStateReader) Invoke(methodad-gome string,engine *exec.ExecutionEngine) (bool, error) {

	if v, ok := i.serviceMap[methodad-gome]; ok {
		return v(engine)
	}
	return true, errors.New("Not supported method:" + methodad-gome)
}

func (i *WasmStateReader) MergeMap(mMap map[string]func(*exec.ExecutionEngine) (bool, error)) bool {

	for k, v := range mMap {
		if _, ok := i.serviceMap[k]; !ok {
			i.serviceMap[k] = v
		}
	}
	return true
}

func (i *WasmStateReader) GetServiceMap() map[string]func(*exec.ExecutionEngine) (bool, error) {
	return i.serviceMap
}

func (i *WasmStateReader) Exists(name string) bool {
	_, ok := i.serviceMap[name]
	return ok
}
