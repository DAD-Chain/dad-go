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
	"fmt"

	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/vm/neovm/errors"
)

func opNop(e *ExecutionEngine) (VMState, error) {
	return NONE, nil
}

func opJmp(e *ExecutionEngine) (VMState, error) {
	offset := int(e.Context.OpReader.ReadInt16())

	offset = e.Context.GetInstructionPointer() + offset - 3

	if offset < 0 || offset > len(e.Context.Code) {
		log.Error(fmt.Sprintf("[opJmp] offset:%v > e.contex.Code len:%v error", offset, len(e.Context.Code)))
		return FAULT, errors.ERR_FAULT
	}
	var fValue = true

	if e.OpCode > JMP {
		if EvaluationStackCount(e) < 1 {
			log.Error(fmt.Sprintf("[opJmp] stack count:%v > 1 error", EvaluationStackCount(e)))
			return FAULT, errors.ERR_UNDER_STACK_LEN
		}
		fValue = PopBoolean(e)
		if e.OpCode == JMPIFNOT {
			fValue = !fValue
		}
	}

	if fValue {
		e.Context.SetInstructionPointer(int64(offset))
	}
	return NONE, nil
}

func opCall(e *ExecutionEngine) (VMState, error) {
	context := e.Context.Clone()
	e.Context.SetInstructionPointer(int64(e.Context.GetInstructionPointer() + 2))
	e.OpCode = JMP
	e.PushContext(context)
	return opJmp(e)
}

func opRet(e *ExecutionEngine) (VMState, error) {
	e.PopContext()
	return NONE, nil
}
