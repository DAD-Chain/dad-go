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
)

// AttributeGetUsage put attribute's usage to vm stack
func AttributeGetUsage(service *NeoVmService, engine *vm.ExecutionEngine) error {
	i, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	vm.PushData(engine, int(i.(*types.TxAttribute).Usage))
	return nil
}

// AttributeGetData put attribute's data to vm stack
func AttributeGetData(service *NeoVmService, engine *vm.ExecutionEngine) error {
	i, err := vm.PopInteropInterface(engine)
	if err != nil {
		return err
	}
	vm.PushData(engine, i.(*types.TxAttribute).Data)
	return nil
}
