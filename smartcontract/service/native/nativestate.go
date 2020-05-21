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
	"github.com/dad-go/smartcontract/storage"
	scommon "github.com/dad-go/core/store/common"
	"bytes"
	"github.com/dad-go/common/serialization"
	"github.com/dad-go/core/types"
	"github.com/dad-go/smartcontract/event"
	"fmt"
)

type (
	Handler func(native *NativeService) (bool, error)
)

type NativeService struct {
	CloneCache *storage.CloneCache
	ServiceMap  map[string]Handler
	Notifications []*event.NotifyEventInfo
	Input []byte
	Tx *types.Transaction
}

func NewNativeService(dbCache scommon.IStateStore, input []byte, tx *types.Transaction) *NativeService {
	var nativeService NativeService
	nativeService.CloneCache = storage.NewCloneCache(dbCache)
	nativeService.Input = input
	nativeService.Tx = tx
	nativeService.ServiceMap = make(map[string]Handler)
	nativeService.Register("Token.Common.Transfer", Transfer)
	nativeService.Register("Token.Ont.Init", OntInit)
	return &nativeService
}

func(native *NativeService) Register(methodad-gome string, handler Handler) {
	native.ServiceMap[methodad-gome] = handler
}

func(native *NativeService) Invoke() (bool, error){
	bf := bytes.NewBuffer(native.Input)
	serviceName, err := serialization.ReadVarBytes(bf); if err != nil {
		return false, err
	}
	service, ok := native.ServiceMap[string(serviceName)]; if !ok {
		return false, fmt.Errorf("Native does not support this service:%s !",serviceName)
	}
	native.Input = bf.Bytes()
	return service(native)
}








