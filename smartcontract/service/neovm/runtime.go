package neovm

import (
	vm "github.com/ontio/dad-go/vm/neovm"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/smartcontract/event"
	scommon "github.com/ontio/dad-go/smartcontract/common"
)

func RuntimeGetTime(service *NeoVmService, engine *vm.ExecutionEngine) error {
	vm.PushData(engine, int(service.Time))
	return nil
}

func RuntimeCheckWitness(service *NeoVmService, engine *vm.ExecutionEngine) error {
	if vm.EvaluationStackCount(engine) < 1 {
		return errors.NewErr("[RuntimeCheckWitness] Too few input parameters ")
	}
	data := vm.PopByteArray(engine)
	var result bool
	if len(data) == 20 {
		address, err := common.AddressParseFromBytes(data)
		if err != nil {
			return err
		}
		result = service.ContextRef.CheckWitness(address)
	} else {
		pk, err := keypair.DeserializePublicKey(data); if err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[RuntimeCheckWitness] data invalid.")
		}
		result = service.ContextRef.CheckWitness(types.AddressFromPubKey(pk))
	}

	vm.PushData(engine, result)
	return nil
}

func RuntimeNotify(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item := vm.PopStackItem(engine)
	context := service.ContextRef.CurrentContext()
	service.Notifications = append(service.Notifications, &event.NotifyEventInfo{TxHash: service.Tx.Hash(), ContractAddress: context.ContractAddress, States: scommon.ConvertReturnTypes(item)})
	return nil
}

func RuntimeLog(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item := vm.PopByteArray(engine)
	context := service.ContextRef.CurrentContext()
	txHash := service.Tx.Hash()
	event.PushSmartCodeEvent(txHash, 0, "InvokeTransaction", &event.LogEventArgs{TxHash:txHash, ContractAddress: context.ContractAddress, Message: string(item)})
	return nil
}




