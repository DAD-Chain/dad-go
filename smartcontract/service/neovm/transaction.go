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
	"github.com/ontio/dad-go/core/types"
	vm "github.com/ontio/dad-go/vm/neovm"
	vmtypes "github.com/ontio/dad-go/vm/neovm/types"
)

// GetExecutingAddress push transaction's hash to vm stack
func TransactionGetHash(service *NeoVmService, engine *vm.ExecutionEngine) error {
	txn := vm.PopInteropInterface(engine).(*types.Transaction)
	txHash := txn.Hash()
	vm.PushData(engine, txHash.ToArray())
	return nil
}

// TransactionGetType push transaction's type to vm stack
func TransactionGetType(service *NeoVmService, engine *vm.ExecutionEngine) error {
	txn := vm.PopInteropInterface(engine).(*types.Transaction)
	vm.PushData(engine, int(txn.TxType))
	return nil
}

// TransactionGetAttributes push transaction's attributes to vm stack
func TransactionGetAttributes(service *NeoVmService, engine *vm.ExecutionEngine) error {
	txn := vm.PopInteropInterface(engine).(*types.Transaction)
	attributes := txn.Attributes
	attributList := make([]vmtypes.StackItems, 0)
	for _, v := range attributes {
		attributList = append(attributList, vmtypes.NewInteropInterface(v))
	}
	vm.PushData(engine, attributList)
	return nil
}
