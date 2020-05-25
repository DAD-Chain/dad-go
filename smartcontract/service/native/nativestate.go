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
	"github.com/dad-go/core/types"
	"github.com/dad-go/smartcontract/event"
	"github.com/dad-go/common"
	"github.com/dad-go/smartcontract/context"
	"github.com/dad-go/core/genesis"
	"github.com/dad-go/smartcontract/service/native/states"
	"github.com/dad-go/errors"
	vmtypes "github.com/dad-go/vm/types"
	"fmt"
)

type (
	Handler func(native *NativeService) error
	RegisterService func(native *NativeService)
)

var (
	Contracts = map[common.Address]RegisterService{
		genesis.OntContractAddress: RegisterOntContract,
		genesis.OngContractAddress: RegisterOngContract,
	}
)

type NativeService struct {
	CloneCache *storage.CloneCache
	ServiceMap  map[string]Handler
	Notifications []*event.NotifyEventInfo
	Input []byte
	Tx *types.Transaction
	Height uint32
	ContextRef context.ContextRef
}

func NewNativeService(dbCache scommon.IStateStore, height uint32, tx *types.Transaction, ctxRef context.ContextRef) *NativeService {
	var nativeService NativeService
	nativeService.CloneCache = storage.NewCloneCache(dbCache)
	nativeService.Tx = tx
	nativeService.Height = height
	nativeService.ContextRef = ctxRef
	nativeService.ServiceMap = make(map[string]Handler)
	return &nativeService
}

func(native *NativeService) Register(methodad-gome string, handler Handler) {
	native.ServiceMap[methodad-gome] = handler
}

func(native *NativeService) Invoke() error {
	ctx := native.ContextRef.CurrentContext()
	if ctx == nil {
		return errors.NewErr("[Invoke] Native service current context doesn't exist!")
	}
	bf := bytes.NewBuffer(ctx.Code.Code)
	contract := new(states.Contract)
	if err := contract.Deserialize(bf); err != nil {
		return err
	}
	services, ok := Contracts[contract.Address]; if !ok {
		return fmt.Errorf("Native contract address %x haven't been registered.", contract.Address)
	}
	services(native)
	service, ok := native.ServiceMap[contract.Method]; if !ok {
		return fmt.Errorf("Native contract %x doesn't support this function %s.", contract.Address, contract.Method)
	}
	native.ContextRef.PushContext(&context.Context{ContractAddress: contract.Address})
	native.Input = contract.Args
	if err := service(native); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[Invoke] Native serivce function execute error!")
	}
	native.ContextRef.PopContext()
	native.CloneCache.Commit()
	return nil
}

func(native *NativeService) AppCall(address common.Address, method string, args []byte) error {
	bf := new(bytes.Buffer)
	contract := &states.Contract{
		Address: address,
		Method: method,
		Args: args,
	}

	if err := contract.Serialize(bf); err != nil {
		return err
	}

	native.ContextRef.PushContext(&context.Context{
		Code: vmtypes.VmCode{
			VmType: vmtypes.Native,
			Code: bf.Bytes(),
		},
	})
	if err := native.ContextRef.Execute(); err != nil {
		return err
	}
	return nil
}

func RegisterOntContract(native *NativeService) {
	native.Register("init", OntInit)
	native.Register("transfer", OntTransfer)
	native.Register("approve", OntApprove)
	native.Register("transferFrom", OntTransferFrom)
}

func RegisterOngContract(native *NativeService) {
	native.Register("init", OngInit)
	native.Register("transfer", OngTransfer)
	native.Register("approve", OngApprove)
	native.Register("transferFrom", OngTransferFrom)
}








