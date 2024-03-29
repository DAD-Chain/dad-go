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
	"github.com/ontio/dad-go/vm/neovm/errors"
	"math/big"
)

func opNop(e *ExecutionEngine) (VMState, error) {
	return NONE, nil
}

func opJmp(e *ExecutionEngine) (VMState, error) {
	num, err := e.Context.OpReader.ReadInt16()
	if err != nil {
		return FAULT, err
	}
	offset := int(num)

	offset = e.Context.GetInstructionPointer() + offset - 3

	if offset < 0 || offset > len(e.Context.Code) {
		return FAULT, errors.ERR_FAULT
	}
	var fValue = true

	if e.OpCode > JMP {
		if EvaluationStackCount(e) < 1 {
			return FAULT, errors.ERR_UNDER_STACK_LEN
		}
		var err error
		fValue, err = PopBoolean(e)
		if err != nil {
			return FAULT, err
		}
		if e.OpCode == JMPIFNOT {
			fValue = !fValue
		}
	}

	if fValue {
		err := e.Context.SetInstructionPointer(int64(offset))
		if err != nil {
			return FAULT, err
		}
	}
	return NONE, nil
}

func opCall(e *ExecutionEngine) (VMState, error) {
	context := e.Context.Clone()
	err := e.Context.SetInstructionPointer(int64(e.Context.GetInstructionPointer() + 2))
	if err != nil {
		return FAULT, err
	}
	e.OpCode = JMP
	e.PushContext(context)
	return opJmp(e)
}

func opDCALL(e *ExecutionEngine) (VMState, error) {
	context := e.Context.Clone()
	e.PushContext(context)

	dest, err := PopBigInt(e)
	if err != nil {
		return FAULT, errors.ERR_DCALL_OFFSET_ERROR
	}

	if dest.Sign() < 0 || dest.Cmp(big.NewInt(int64(len(e.Context.Code)))) >= 0 {
		return FAULT, errors.ERR_DCALL_OFFSET_ERROR
	}

	target := dest.Int64()

	err = e.Context.SetInstructionPointer(target)
	if err != nil {
		return FAULT, err
	}

	return NONE, nil
}

func opRet(e *ExecutionEngine) (VMState, error) {
	e.PopContext()
	return NONE, nil
}
