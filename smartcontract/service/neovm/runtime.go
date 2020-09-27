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
	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/signature"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/errors"
	scommon "github.com/ontio/dad-go/smartcontract/common"
	"github.com/ontio/dad-go/smartcontract/event"
	vm "github.com/ontio/dad-go/vm/neovm"
)

// HeaderGetNextConsensus put current block time to vm stack
func RuntimeGetTime(service *NeoVmService, engine *vm.ExecutionEngine) error {
	vm.PushData(engine, int(service.Time))
	return nil
}

// RuntimeCheckWitness provide check permissions service
// If param address isn't exist in authorization list, check fail
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
		pk, err := keypair.DeserializePublicKey(data)
		if err != nil {
			return errors.NewDetailErr(err, errors.ErrNoCode, "[RuntimeCheckWitness] data invalid.")
		}
		result = service.ContextRef.CheckWitness(types.AddressFromPubKey(pk))
	}

	vm.PushData(engine, result)
	return nil
}

// RuntimeNotify put smart contract execute event notify to notifications
func RuntimeNotify(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item := vm.PopStackItem(engine)
	context := service.ContextRef.CurrentContext()
	service.Notifications = append(service.Notifications, &event.NotifyEventInfo{ContractAddress: context.ContractAddress, States: scommon.ConvertNeoVmTypeHexString(item)})
	return nil
}

// RuntimeLog push smart contract execute event log to client
func RuntimeLog(service *NeoVmService, engine *vm.ExecutionEngine) error {
	item := vm.PopByteArray(engine)
	context := service.ContextRef.CurrentContext()
	txHash := service.Tx.Hash()
	event.PushSmartCodeEvent(txHash, 0, event.EVENT_LOG, &event.LogEventArgs{TxHash: txHash, ContractAddress: context.ContractAddress, Message: string(item)})
	return nil
}

// RuntimeCheckSig verify whether authorization legal
func RuntimeCheckSig(service *NeoVmService, engine *vm.ExecutionEngine) error {
	pubKey := vm.PopByteArray(engine)
	data := vm.PopByteArray(engine)
	sig := vm.PopByteArray(engine)
	return signature.Verify(pubKey, data, sig)
}
