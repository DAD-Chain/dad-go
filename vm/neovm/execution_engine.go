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
)

func NewExecutionEngine() *ExecutionEngine {
	var engine ExecutionEngine
	engine.EvaluationStack = NewRandAccessStack()
	engine.AltStack = NewRandAccessStack()
	engine.State = BREAK
	engine.OpCode = 0
	return &engine
}

type ExecutionEngine struct {
	EvaluationStack *RandomAccessStack
	AltStack        *RandomAccessStack
	State           VMState
	Contexts        []*ExecutionContext
	Context         *ExecutionContext
	OpCode          OpCode
	OpExec          OpExec
}

func (this *ExecutionEngine) CurrentContext() *ExecutionContext {
	return this.Contexts[len(this.Contexts)-1]
}

func (this *ExecutionEngine) PopContext() {
	if len(this.Contexts) != 0 {
		this.Contexts = this.Contexts[:len(this.Contexts)-1]
	}
	if len(this.Contexts) != 0 {
		this.Context = this.CurrentContext()
	} else {
		this.Context = nil
	}
}

func (this *ExecutionEngine) PushContext(context *ExecutionContext) {
	this.Contexts = append(this.Contexts, context)
	this.Context = this.CurrentContext()
}

func (this *ExecutionEngine) Execute() error {
	this.State = this.State & (^BREAK)
	for {
		if this.State == FAULT || this.State == HALT || this.State == BREAK {
			break
		}
		if this.Context == nil {
			break
		}
		err := this.ExecuteCode()
		if err != nil {
			break
		}

		if this.OpCode >= PUSHBYTES1 && this.OpCode <= PUSHBYTES75 {
			bs, err := this.Context.OpReader.ReadBytes(int(this.OpCode))
			if err != nil {
				this.State = FAULT
				return err
			}
			PushData(this, bs)
			continue
		}

		err = this.ValidateOp()
		if err != nil {
			this.State = FAULT
			return err
		}

		err = this.StepInto()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *ExecutionEngine) ExecuteCode() error {
	code, err := this.Context.OpReader.ReadByte()
	if err != nil {
		this.State = FAULT
		return err
	}
	this.OpCode = OpCode(code)
	return nil
}

func (this *ExecutionEngine) ValidateOp() error {
	opExec := OpExecList[this.OpCode]
	if opExec.Name == "" {
		return errors.ERR_NOT_SUPPORT_OPCODE
	}
	this.OpExec = opExec
	return nil
}

func (this *ExecutionEngine) StepInto() error {
	state, err := this.ExecuteOp()
	this.State = state
	if err != nil {
		return err
	}
	return nil
}

func (this *ExecutionEngine) ExecuteOp() (VMState, error) {
	if this.OpCode >= PUSHBYTES1 && this.OpCode <= PUSHBYTES75 {
		bs, err := this.Context.OpReader.ReadBytes(int(this.OpCode))
		if err != nil {
			return FAULT, err
		}
		PushData(this, bs)
		return NONE, nil
	}

	if this.OpExec.Validator != nil {
		if err := this.OpExec.Validator(this); err != nil {
			return FAULT, err
		}
	}
	return this.OpExec.Exec(this)
}
