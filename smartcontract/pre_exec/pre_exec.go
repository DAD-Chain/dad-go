package pre_exec

import (
	"github.com/dad-go/smartcontract/service"
	"github.com/dad-go/vm/neovm"
	"github.com/dad-go/vm/neovm/interfaces"
	"github.com/dad-go/core/store/statestore"
	"github.com/dad-go/smartcontract/types"
	"github.com/dad-go/core/store/ChainStore"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/smartcontract/common"
)

func PreExec(code []byte, container interfaces.ICodeContainer) ([]byte, error) {
	var (
		crypto interfaces.ICrypto
		err error
	)
	crypto = new(neovm.ECDsaCrypto)
	stateStore := ChainStore.NewStateStore(statestore.NewMemDatabase(), ledger.DefaultLedger.Store, nil, nil)
	stateMachine := service.NewStateMachine(stateStore, types.Application, nil)
	se := neovm.NewExecutionEngine(container, crypto, &ChainStore.NewCacheCodeTable(stateStore), stateMachine)
	se.LoadCode(code, false)
	err = se.Execute()
	if err != nil {
		return []byte{}, err
	}
	if se.GetEvaluationStackCount() == 0 {
		return []byte{}, err
	}
	if neovm.Peek(se).GetStackItem() == nil {
		return []byte{}, err
	}
	return common.ConvertReturnTypes(neovm.Peek(se).GetStackItem()), nil
}
