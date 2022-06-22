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
package common

import (
	"encoding/json"
	utils2 "github.com/ontio/dad-go/cmd/utils"
	"github.com/ontio/dad-go/core/utils"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/types"
)

type TestEnv struct {
	Witness []common.Address `json:"witness"`
}

func (self TestEnv) MarshalJSON() ([]byte, error) {
	var witness []string
	for _, addr := range self.Witness {
		witness = append(witness, addr.ToBase58())
	}
	env := struct {
		Witness []string `json:"witness"`
	}{Witness: witness}

	return json.Marshal(env)
}

func (self *TestEnv) UnmarshalJSON(buf []byte) error {
	env := struct {
		Witness []string `json:"witness"`
	}{}
	err := json.Unmarshal(buf, &env)
	if err != nil {
		return err
	}
	var witness []common.Address
	for _, addr := range env.Witness {
		wit, err := common.AddressFromBase58(addr)
		if err != nil {
			return err
		}

		witness = append(witness, wit)
	}

	self.Witness = witness
	return nil
}

type TestCase struct {
	Env    TestEnv `json:"env"`
	Method string  `json:"method"`
	Param  string  `json:"param"`
	Expect string  `json:"expected"`
}

func GenWasmTransaction(testCase TestCase, contract common.Address) (*types.Transaction, error) {
	params, err := utils2.ParseParams(testCase.Param)
	if err != nil {
		return nil, err
	}
	allParam := append([]interface{}{}, testCase.Method)
	allParam = append(allParam, params...)
	tx, err := utils.NewWasmVMInvokeTransaction(0, 100000000, contract, allParam)
	if err != nil {
		return nil, err
	}

	imt, err := tx.IntoImmutable()
	if err != nil {
		return nil, err
	}

	imt.SignedAddr = append(imt.SignedAddr, testCase.Env.Witness...)

	return imt, nil
}