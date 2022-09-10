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
package util

import (
	"bytes"

	"errors"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/utils"
	"github.com/ontio/dad-go/smartcontract/context"
	neovms "github.com/ontio/dad-go/smartcontract/service/neovm"
	"github.com/ontio/dad-go/vm/crossvm_codec"
	"github.com/ontio/dad-go/vm/neovm"
)

func BuildNeoVMParamEvalStack(params []interface{}) (*neovm.ValueStack, error) {
	builder := neovm.NewParamsBuilder(new(bytes.Buffer))
	err := utils.BuildNeoVMParam(builder, params)
	if err != nil {
		return nil, err
	}

	exec := neovm.NewExecutor(builder.ToArray(), neovm.VmFeatureFlag{true, true})
	err = exec.Execute()
	if err != nil {
		return nil, err
	}
	return exec.EvalStack, nil
}

//create paramters for neovm contract
func GenerateNeoVMParamEvalStack(input []byte) (*neovm.ValueStack, error) {
	params, err := crossvm_codec.DeserializeCallParam(input)
	if err != nil {
		return nil, err
	}

	list, ok := params.([]interface{})
	if ok == false {
		return nil, errors.New("invoke neovm param is not list type")
	}

	stack, err := BuildNeoVMParamEvalStack(list)
	if err != nil {
		return nil, err
	}

	return stack, nil
}

func SetNeoServiceParamAndEngine(addr common.Address, engine context.Engine, stack *neovm.ValueStack) error {
	service, ok := engine.(*neovms.NeoVmService)
	if ok == false {
		return errors.New("engine should be NeoVmService")
	}

	code, err := service.GetNeoContract(addr)
	if err != nil {
		return err
	}

	feature := service.Engine.Features
	service.Engine = neovm.NewExecutor(code, feature)
	service.Code = code

	service.Engine.EvalStack = stack

	return nil
}
