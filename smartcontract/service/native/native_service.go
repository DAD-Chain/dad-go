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

package native

import (
	"bytes"
	"fmt"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/smartcontract/context"
	"github.com/ontio/dad-go/smartcontract/event"
	"github.com/ontio/dad-go/smartcontract/storage"
	sstates "github.com/ontio/dad-go/smartcontract/states"
)

type (
	Handler         func(native *NativeService) error
	RegisterService func(native *NativeService)
)

var (
	Contracts = make(map[common.Address]RegisterService)
)

// Native service struct
// Invoke a native smart contract, new a native service
type NativeService struct {
	CloneCache    *storage.CloneCache
	ServiceMap    map[string]Handler
	Notifications []*event.NotifyEventInfo
	Code          []byte
	Input         []byte
	Tx            *types.Transaction
	Height        uint32
	Time          uint32
	ContextRef    context.ContextRef
}

func (this *NativeService) Register(methodad-gome string, handler Handler) {
	this.ServiceMap[methodad-gome] = handler
}

func (this *NativeService) Invoke() (interface{}, error) {
	bf := bytes.NewBuffer(this.Code)
	contract := new(sstates.Contract)
	if err := contract.Deserialize(bf); err != nil {
		return false, err
	}
	services, ok := Contracts[contract.Address]
	if !ok {
		return false, fmt.Errorf("Native contract address %x haven't been registered.", contract.Address)
	}
	services(this)
	service, ok := this.ServiceMap[contract.Method]
	if !ok {
		return false, fmt.Errorf("Native contract %x doesn't support this function %s.",
			contract.Address, contract.Method)
	}
	this.Input = contract.Args
	this.ContextRef.PushContext(&context.Context{ContractAddress: contract.Address})
	if err := service(this); err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode,
			"[Invoke] Native serivce function execute error!")
	}
	this.ContextRef.PopContext()
	this.ContextRef.PushNotifications(this.Notifications)
	return true, nil
}
