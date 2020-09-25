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

package global_params

import (
	"bytes"

	"github.com/ontio/dad-go/common"
	cstates "github.com/ontio/dad-go/core/states"
	"github.com/ontio/dad-go/smartcontract/event"
	"github.com/ontio/dad-go/smartcontract/service/native"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
)

const (
	SET_PARAM = "SetGlobalParam"
	PARAM     = "param"
	TRANSFER  = "transfer"
	ADMIN     = "admin"
)

func getAdminStorageItem(admin *Admin) *cstates.StorageItem {
	bf := new(bytes.Buffer)
	admin.Serialize(bf)
	return &cstates.StorageItem{Value: bf.Bytes()}
}

func getParamStorageItem(params *Params) *cstates.StorageItem {
	bf := new(bytes.Buffer)
	params.Serialize(bf)
	return &cstates.StorageItem{Value: bf.Bytes()}
}

func getParamKey(contract common.Address, valueType paramType) []byte {
	key := append(contract[:], PARAM...)
	key = append(key[:], byte(valueType))
	return key
}

func getAdminKey(contract common.Address, isTransferAdmin bool) []byte {
	if isTransferAdmin {
		return append(contract[:], TRANSFER...)
	} else {
		return append(contract[:], ADMIN...)
	}
}

func notifyParamSetSuccess(native *native.NativeService, contract common.Address, params Params) {
	native.Notifications = append(native.Notifications,
		&event.NotifyEventInfo{
			TxHash:          native.Tx.Hash(),
			ContractAddress: contract,
			States:          []interface{}{SET_PARAM, params},
		})
}

func getStorageParam(native *native.NativeService, key []byte) (*Params, error) {
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	params := new(Params)
	bf := bytes.NewBuffer(item.Value)
	params.Deserialize(bf)
	return params, nil
}

func getStorageAdmin(native *native.NativeService, key []byte) (*Admin, error) {
	item, err := utils.GetStorageItem(native, key)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	admin := new(Admin)
	bf := bytes.NewBuffer(item.Value)
	admin.Deserialize(bf)
	return admin, nil
}
