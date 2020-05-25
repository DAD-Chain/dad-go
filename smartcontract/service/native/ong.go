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

package native

import (
	scommon "github.com/dad-go/core/store/common"
	"github.com/dad-go/errors"
	"math/big"
	"github.com/dad-go/smartcontract/service/native/states"
	cstates "github.com/dad-go/core/states"
	"bytes"
	"github.com/dad-go/core/genesis"
)

var (
	totalSupplyName = []byte("totalSupply")
	decimals = big.NewInt(8)
	totalSupply = new(big.Int).Mul(big.NewInt(1000000000), (new(big.Int).Exp(big.NewInt(10), decimals, nil)))
)

func OngInit(native *NativeService) error {
	contract := native.ContextRef.CurrentContext().ContractAddress
	amount, err := getStorageBigInt(native, getTotalSupplyKey(contract))
	if err != nil {
		return err
	}

	if amount != nil && amount.Sign() != 0 {
		return errors.NewErr("Init ong has been completed!")
	}
	native.CloneCache.Add(scommon.ST_Storage, append(contract[:], getOntContext()...), &cstates.StorageItem{Value: totalSupply.Bytes()})
	addNotifications(native, contract, &states.State{To: genesis.OntContractAddress, Value: totalSupply})
	return nil
}

func OngTransfer(native *NativeService) error {
	transfers := new(states.Transfers)
	if err := transfers.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[OngTransfer] Transfers deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	for _, v := range transfers.States {
		if err := transfer(native, contract, v); err != nil {
			return err
		}
		addNotifications(native, contract, v)
	}
	return nil
}

func OngApprove(native *NativeService) error {
	state := new(states.State)
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[OngApprove] state deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	native.CloneCache.Add(scommon.ST_Storage, getApproveKey(contract, state), &cstates.StorageItem{Value: state.Value.Bytes()})
	return nil
}

func OngTransferFrom(native *NativeService) error {
	state := new(states.TransferFrom)
	if err := state.Deserialize(bytes.NewBuffer(native.Input)); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[OntTransferFrom] State deserialize error!")
	}
	contract := native.ContextRef.CurrentContext().ContractAddress
	if err := transferFrom(native, contract, state); err != nil {
		return err
	}
	addNotifications(native, contract, &states.State{From: state.From, To: state.To, Value: state.Value})
	return nil
}

func getOntContext() []byte {
	return genesis.OntContractAddress[:]
}


