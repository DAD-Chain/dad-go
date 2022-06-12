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
	"github.com/ontio/dad-go/errors"
	vm "github.com/ontio/dad-go/vm/neovm"
	vmtypes "github.com/ontio/dad-go/vm/neovm/types"
)

// BlockGetTransactionCount put block's transactions count to vm stack
func BlockGetTransactionCount(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	if block, ok := i.Data.(*types.Block); ok {
		val := vmtypes.VmValueFromInt64(int64(len(block.Transactions)))
		engine.EvalStack.Push(val)
		return nil
	}
	return errors.NewErr("[BlockGetTransactionCount] Wrong type ")
}

// BlockGetTransactions put block's transactions to vm stack
func BlockGetTransactions(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	if block, ok := i.Data.(*types.Block); ok {
		transactions := block.Transactions
		transactionList := make([]vmtypes.VmValue, 0)

		for _, v := range transactions {
			transactionList = append(transactionList, vmtypes.VmValueFromInteropValue(vmtypes.NewInteropValue(v)))
		}

		engine.EvalStack.PushAsArray(transactionList)
	}
	return errors.NewErr("[BlockGetTransactions] Wrong type ")
}

// BlockGetTransaction put block's transaction to vm stack
func BlockGetTransaction(service *NeoVmService, engine *vm.Executor) error {
	i, err := engine.EvalStack.PopAsInteropValue()
	if err != nil {
		return err
	}
	index, err := engine.EvalStack.PopAsInt64()
	if err != nil {
		return err
	}
	if block, ok := i.Data.(*types.Block); ok {
		if index < 0 || int(index) >= len(block.Transactions) {
			return errors.NewErr("[BlockGetTransaction] index out of bounds")
		}
		return engine.EvalStack.PushAsInteropValue(block.Transactions[index])

	}
	return errors.NewErr("[BlockGetTransaction] Wrong type ")
}
