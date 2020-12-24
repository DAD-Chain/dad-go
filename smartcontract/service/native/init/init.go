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

package init

import (
	"bytes"
	"math"
	"math/big"

	"github.com/ontio/dad-go/common"
	invoke "github.com/ontio/dad-go/core/utils"
	"github.com/ontio/dad-go/smartcontract/service/native/auth"
	params "github.com/ontio/dad-go/smartcontract/service/native/global_params"
	"github.com/ontio/dad-go/smartcontract/service/native/governance"
	"github.com/ontio/dad-go/smartcontract/service/native/ong"
	"github.com/ontio/dad-go/smartcontract/service/native/ont"
	"github.com/ontio/dad-go/smartcontract/service/native/ontid"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
	"github.com/ontio/dad-go/smartcontract/service/neovm"
	vm "github.com/ontio/dad-go/vm/neovm"
)

var (
	COMMIT_DPOS_BYTES = InitBytes(utils.GovernanceContractAddress, governance.COMMIT_DPOS)
)

func init() {
	ong.InitOng()
	ont.InitOnt()
	params.InitGlobalParams()
	ontid.Init()
	auth.Init()
	governance.InitGovernance()
}

func InitBytes(addr common.Address, method string) []byte {
	bf := new(bytes.Buffer)
	builder := vm.NewParamsBuilder(bf)
	builder.EmitPushByteArray([]byte{})
	builder.EmitPushByteArray([]byte(method))
	builder.EmitPushByteArray(addr[:])
	builder.EmitPushInteger(big.NewInt(0))
	builder.Emit(vm.SYSCALL)
	builder.EmitPushByteArray([]byte(neovm.NATIVE_INVOKE_NAME))

	tx := invoke.NewInvokeTransaction(builder.ToArray())
	tx.GasLimit = math.MaxUint64
	return bf.Bytes()
}
