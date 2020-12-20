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

package utils

import (
	"bytes"
	"math/big"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/payload"
	"github.com/ontio/dad-go/core/types"
	neovm "github.com/ontio/dad-go/smartcontract/service/neovm"
	vm "github.com/ontio/dad-go/vm/neovm"
	"math"
)

// NewDeployTransaction returns a deploy Transaction
func NewDeployTransaction(code []byte, name, version, author, email, desp string, needStorage bool) *types.Transaction {
	//TODO: check arguments
	DeployCodePayload := &payload.DeployCode{
		Code:        code,
		NeedStorage: needStorage,
		Name:        name,
		Version:     version,
		Author:      author,
		Email:       email,
		Description: desp,
	}

	return &types.Transaction{
		TxType:  types.Deploy,
		Payload: DeployCodePayload,
	}
}

// NewInvokeTransaction returns an invoke Transaction
func NewInvokeTransaction(code []byte) *types.Transaction {
	//TODO: check arguments
	invokeCodePayload := &payload.InvokeCode{
		Code: code,
	}

	return &types.Transaction{
		TxType:  types.Invoke,
		Payload: invokeCodePayload,
	}
}

func Buildad-gotiveTransaction(addr common.Address, initMethod string, args []byte) *types.Transaction {
	bf := new(bytes.Buffer)
	builder := vm.NewParamsBuilder(bf)
	builder.EmitPushByteArray(args)
	builder.EmitPushByteArray([]byte(initMethod))
	builder.EmitPushByteArray(addr[:])
	builder.EmitPushInteger(big.NewInt(0))
	builder.Emit(vm.SYSCALL)
	builder.EmitPushByteArray([]byte(neovm.NATIVE_INVOKE_NAME))

	tx := NewInvokeTransaction(builder.ToArray())
	tx.GasLimit = math.MaxUint64
	return tx
}
