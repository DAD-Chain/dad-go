package wasm

import (
	vmtypes "github.com/dad-go/smartcontract/types"
	"github.com/dad-go/smartcontract/storage"
	"github.com/dad-go/core/store"
	scommon "github.com/dad-go/core/store/common"
	"github.com/dad-go/core/types"
)

type WasmStateMachine struct {
	*WasmStateReader
	ldgerStore store.ILedgerStore
	CloneCache *storage.CloneCache
	trigger    vmtypes.TriggerType
	block       *types.Block
}


func NewWasmStateMachine(ldgerStore store.ILedgerStore, dbCache scommon.IStateStore, trigger vmtypes.TriggerType, block *types.Block) *WasmStateMachine {
	var stateMachine WasmStateMachine
	stateMachine.ldgerStore = ldgerStore
	stateMachine.CloneCache = storage.NewCloneCache(dbCache)
	stateMachine.WasmStateReader = NewWasmStateReader(ldgerStore,trigger)
	stateMachine.trigger = trigger
	stateMachine.block = block

	//stateMachine.Register("getBlockHeight",bcGetHeight)
	//todo add and register services
	return &stateMachine
}

//======================some block api ===============
/*
func  bcGetHeight(engine *exec.ExecutionEngine) (bool, error) {
	vm := engine.GetVM()
	var i uint32
	if ledger.DefaultLedger == nil {
		i = 0
	} else {
		i = ledger.DefaultLedger.Store.GetHeight()
	}
	//engine.vm.ctx = envCall.envPreCtx
	vm.RestoreCtx()
	if vm.GetEnvCall().GetReturns(){
		vm.PushResult(uint64(i))
	}
	return true,nil
}*/
