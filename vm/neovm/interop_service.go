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
	"github.com/ontio/dad-go/vm/neovm/errors"
)

type InteropServices interface {
	Register(method string, handler func(*ExecutionEngine) (bool, error)) bool
	GetServiceMap() map[string]func(*ExecutionEngine) (bool, error)
}

type InteropService struct {
	serviceMap map[string]func(*ExecutionEngine) (bool, error)
}

func NewInteropService() *InteropService {
	var i InteropService
	i.serviceMap = make(map[string]func(*ExecutionEngine) (bool, error), 0)
	i.Register("System.ExecutionEngine.GetScriptContainer", i.GetCodeContainer)
	i.Register("System.ExecutionEngine.GetExecutingScriptHash", i.GetExecutingCodeHash)
	i.Register("System.ExecutionEngine.GetCallingScriptHash", i.GetCallingCodeHash)
	i.Register("System.ExecutionEngine.GetEntryScriptHash", i.GetEntryCodeHash)
	return &i
}

func (is *InteropService) Register(methodad-gome string, handler func(*ExecutionEngine) (bool, error)) bool {
	if _, ok := is.serviceMap[methodad-gome]; ok {
		return false
	}
	is.serviceMap[methodad-gome] = handler
	return true
}

func (i *InteropService) MergeMap(dictionary map[string]func(*ExecutionEngine) (bool, error)) {
	for k, v := range dictionary {
		if _, ok := i.serviceMap[k]; !ok {
			i.serviceMap[k] = v
		}
	}
}

func (i *InteropService) GetServiceMap() map[string]func(*ExecutionEngine) (bool, error) {
	return i.serviceMap
}

func (i *InteropService) Invoke(methodad-gome string, engine *ExecutionEngine) (bool, error) {
	if v, ok := i.serviceMap[methodad-gome]; ok {
		return v(engine)
	}
	return false, errors.ERR_NOT_SUPPORT_SERVICE
}

func (i *InteropService) GetCodeContainer(engine *ExecutionEngine) (bool, error) {
	PushData(engine, engine.codeContainer)
	return true, nil
}

func (i *InteropService) GetExecutingCodeHash(engine *ExecutionEngine) (bool, error) {
	context, err := engine.CurrentContext()
	if err != nil {
		return false, err
	}
	codeHash, err := context.GetCodeHash()
	if err != nil {
		return false, err
	}
	PushData(engine, codeHash[:])
	return true, nil
}

func (i *InteropService) GetCallingCodeHash(engine *ExecutionEngine) (bool, error) {
	context, err := engine.CallingContext()
	if err != nil {
		return false, err
	}
	codeHash, err := context.GetCodeHash()
	if err != nil {
		return false, err
	}
	PushData(engine, codeHash[:])
	return true, nil
}
func (i *InteropService) GetEntryCodeHash(engine *ExecutionEngine) (bool, error) {
	context, err := engine.EntryContext()
	if err != nil {
		return false, err
	}
	codeHash, err := context.GetCodeHash()
	if err != nil {
		return false, err
	}
	PushData(engine, codeHash[:])
	return true, nil
}
