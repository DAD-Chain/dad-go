package neovm

import (
	vm "github.com/ontio/dad-go/vm/neovm"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/smartcontract/event"
	scommon "github.com/ontio/dad-go/smartcontract/common"
	"github.com/ontio/dad-go/core/signature"
)

// get current time
func RuntimeGetTime(service *NeoVmService, engine *vm.ExecutionEngine) error {
	vm.PushData(engine, int(service.Time))
	return nil
}

// check permissions
// if param address isn't exist in authorization list, check fail
func RuntimeCheckWitness(service *NeoVmService, engine *vm.ExecutionEngine) error {
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

// smart contract execute event notify
func RuntimeNotify(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item := vm.PopStackItem(engine)
	context := service.ContextRef.CurrentContext()
	service.Notifications = append(service.Notifications, &event.NotifyEventInfo{TxHash: service.Tx.Hash(), ContractAddress: context.ContractAddress, States: scommon.ConvertReturnTypes(item)})
	return nil
}

// smart contract execute log
func RuntimeLog(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item := vm.PopByteArray(engine)
	context := service.ContextRef.CurrentContext()
	txHash := service.Tx.Hash()
	event.PushSmartCodeEvent(txHash, 0, "InvokeTransaction", &event.LogEventArgs{TxHash:txHash, ContractAddress: context.ContractAddress, Message: string(item)})
	return nil
}

func RuntimeCheckSig(service *NeoVmService, engine *vm.ExecutionEngine) error {
	pubKey := vm.PopByteArray(engine)
	data := vm.PopByteArray(engine)
	sig := vm.PopByteArray(engine)
	return signature.Verify(pubKey, data, sig)
}




