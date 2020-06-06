// Copyright 2017 The dad-go Authors
// This file is part of the dad-go library.
//
// The dad-go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The dad-go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the dad-go library. If not, see <http://www.gnu.org/licenses/>.

package smartcontract

import (
	vmtypes "github.com/dad-go/vm/types"
	"github.com/dad-go/vm/neovm/interfaces"
	ctypes "github.com/dad-go/core/types"
	"github.com/dad-go/smartcontract/service/native"
	scommon "github.com/dad-go/core/store/common"
	sneovm "github.com/dad-go/smartcontract/service/neovm"
	"github.com/dad-go/core/store"
	stypes "github.com/dad-go/smartcontract/types"
	"github.com/dad-go/vm/neovm"
	"github.com/dad-go/smartcontract/context"
	"github.com/dad-go/smartcontract/event"
	"github.com/dad-go/common"
	"github.com/dad-go/smartcontract/service/wasm"
	"github.com/dad-go/vm/wasmvm/exec"
	"github.com/dad-go/vm/wasmvm/util"
)

type SmartContract struct {
	Context []*context.Context
	Config *Config
	Engine Engine
	Notifications []*event.NotifyEventInfo
}

type Config struct {
	Time uint32
	Height uint32
	Tx *ctypes.Transaction
	Table interfaces.CodeTable
	DBCache scommon.StateStore
	Store store.LedgerStore
}

type Engine interface {
	StepInto()
}


//put current context to smart contract
func(sc *SmartContract) PushContext(context *context.Context) {
	sc.Context = append(sc.Context, context)
}

//get smart contract current context
func(sc *SmartContract) CurrentContext() *context.Context {
	if len(sc.Context) < 1 {
		return nil
	}
	return sc.Context[len(sc.Context) - 1]
}

//get smart contract caller context
func(sc *SmartContract) CallingContext() *context.Context {
	if len(sc.Context) < 2 {
		return nil
	}
	return sc.Context[len(sc.Context) - 2]
}

//get smart contract entry entrance context
func(sc *SmartContract) EntryContext() *context.Context {
	if len(sc.Context) < 1 {
		return nil
	}
	return sc.Context[0]
}

//pop smart contract current context
func(sc *SmartContract) PopContext() {
	sc.Context = sc.Context[:len(sc.Context) - 1]
}

func(sc *SmartContract) PushNotifications(notifications []*event.NotifyEventInfo) {
	sc.Notifications = append(sc.Notifications, notifications...)
}

func (sc *SmartContract) Execute() error {
	ctx := sc.CurrentContext()
	switch ctx.Code.VmType {
	case vmtypes.Native:
		service := native.NewNativeService(sc.Config.DBCache, sc.Config.Height, sc.Config.Tx, sc)
		if err := service.Invoke(); err != nil {
			return err
		}
	case vmtypes.NEOVM:
		stateMachine := sneovm.NewStateMachine(sc.Config.Store, sc.Config.DBCache, stypes.Application, sc.Config.Time)
		engine := neovm.NewExecutionEngine(
			sc.Config.Tx,
			new(neovm.ECDsaCrypto),
			sc.Config.Table,
			stateMachine,
		)
		engine.LoadCode(ctx.Code.Code, false)
		if err := engine.Execute(); err != nil {
			return err
		}
		stateMachine.CloneCache.Commit()
		sc.Notifications = append(sc.Notifications, stateMachine.Notifications...)
	case vmtypes.WASMVM:
		stateMachine:= wasm.NewWasmStateMachine(sc.Config.Store, sc.Config.DBCache, stypes.Application,sc.Config.Time)

		engine := exec.NewExecutionEngine(
			sc.Config.Tx,
			new(util.ECDsaCrypto),
			sc.Config.Table,
			stateMachine,
			"product",
		)
		//todo how to get the input
		input:= []byte{}
		engine.Call(ctx.ContractAddress,ctx.Code.Code,input)
		//fmt.Println(engine)
		sc.Notifications = append(sc.Notifications, stateMachine.Notifications...)
	}
	return nil
}

func (sc *SmartContract) CheckWitness(address common.Address) bool {
	if vmtypes.IsVmCodeAddress(address) {
		for _, v := range sc.Context {
			if v.ContractAddress == address {
				return true
			}
		}
	} else {
		addresses := sc.Config.Tx.GetSignatureAddresses()
		for _, v := range addresses {
			if v == address {
				return true
			}
		}
	}

	return false
}
