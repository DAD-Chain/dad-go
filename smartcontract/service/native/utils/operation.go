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
	"encoding/hex"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/smartcontract/event"
	"github.com/ontio/dad-go/smartcontract/service/native"
)

func AddCommonEvent(native *native.NativeService, contract common.Address, name string, params interface{}) {
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			TxHash:          native.Tx.Hash(),
			ContractAddress: contract,
			States:          []interface{}{name, params},
		})
}

func ConcatKey(contract common.Address, args ...[]byte) []byte {
	temp := contract[:]
	for _, arg := range args {
		temp = append(temp, arg...)
	}
	return temp
}

func ValidateOwner(native *native.NativeService, address string) error {
	addrBytes, err := hex.DecodeString(address)
	if err != nil {
		return errors.NewErr("[validateOwner] Decode address hex string to bytes failed!")
	}
	addr, err := common.AddressParseFromBytes(addrBytes)
	if err != nil {
		return errors.NewErr("[validateOwner] Decode bytes to address failed!")
	}
	if native.ContextRef.CheckWitness(addr) == false {
		return errors.NewErr("[validateOwner] Authentication failed!")
	}
	return nil
}


