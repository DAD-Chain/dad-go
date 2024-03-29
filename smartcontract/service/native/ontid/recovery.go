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
package ontid

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/states"
	"github.com/ontio/dad-go/smartcontract/service/native"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
)

const (
	_VERSION_0 byte = 0x00
	_VERSION_1 byte = 0x01
)

func setRecovery(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("setRecovery: argument 0 error")
	}
	// arg1: recovery struct
	arg1, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("setRecovery: argument 1 error")
	}
	// arg2: operator's public key index
	arg2, err := utils.DecodeVarUint(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("setRecovery: argument 2 error")
	}

	encId, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("setRecovery: " + err.Error())
	}
	if err := checkWitnessByIndex(srvc, encId, uint32(arg2)); err != nil {
		return utils.BYTE_FALSE, errors.New("setRecovery: authentication failed: " + err.Error())
	}

	re, err := getRecovery(srvc, encId)
	if err == nil && re != nil {
		return utils.BYTE_FALSE, errors.New("setRecovery: recovery is already set")
	}

	re, err = putRecovery(srvc, encId, arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("setRecovery: " + err.Error())
	}

	newEvent(srvc, []interface{}{"recovery", "set", string(arg0), re.ToJson()})
	return utils.BYTE_TRUE, nil
}

func updateRecovery(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("updateRecovery: argument 0 error")
	}
	// arg1: new recovery
	arg1, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("updateRecovery: argument 1 error")
	}
	// arg2: signers
	arg2, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("updateRecovery: argument 2 error")
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("update recovery: " + err.Error())
	}
	re, err := getRecovery(srvc, key)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("update recovery: get old recovery error, " + err.Error())
	}
	signers, err := deserializeSigners(arg2)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("update recovery: signers error: " + err.Error())
	}

	if !verifyGroupSignature(srvc, re, signers) {
		return utils.BYTE_FALSE, errors.New("update recovery: verification failed")
	}
	re, err = putRecovery(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("update recovery: " + err.Error())
	}

	newEvent(srvc, []interface{}{"Recovery", "update", string(arg0), re.ToJson()})
	return utils.BYTE_TRUE, nil
}

func addKeyByRecovery(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: id
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("argument 0 error")
	}
	// arg1: public key
	arg1, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("argument 1 error")
	}
	// arg2: signers
	arg2, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("argument 2 error")
	}

	encId, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	signers, err := deserializeSigners(arg2)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("signers error, %s", err)
	}

	rec, err := getRecovery(srvc, encId)
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if !verifyGroupSignature(srvc, rec, signers) {
		return utils.BYTE_FALSE, errors.New("verification failed")
	}

	index, err := insertPk(srvc, encId, arg1)
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	triggerPublicEvent(srvc, "add", arg0, arg1, index)
	return utils.BYTE_TRUE, nil
}

func removeKeyByRecovery(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: id
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("argument 0 error")
	}
	// arg1: public key index
	arg1, err := utils.DecodeVarUint(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("argument 1 error")
	}
	// arg2: signers
	arg2, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("argument 2 error")
	}

	encId, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	signers, err := deserializeSigners(arg2)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("signers error, %s", err)
	}

	rec, err := getRecovery(srvc, encId)
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if !verifyGroupSignature(srvc, rec, signers) {
		return utils.BYTE_FALSE, errors.New("verification failed")
	}

	pk, err := revokePkByIndex(srvc, encId, uint32(arg1))
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	triggerPublicEvent(srvc, "remove", arg0, pk, uint32(arg1))
	return utils.BYTE_TRUE, nil
}

func putRecovery(srvc *native.NativeService, encID, data []byte) (*Group, error) {
	rec, err := deserializeGroup(data)
	if err != nil {
		return nil, err
	}
	err = validateMembers(srvc, rec)
	if err != nil {
		return nil, fmt.Errorf("invalid recovery member, %s", err)
	}
	key := append(encID, FIELD_RECOVERY)
	item := states.StorageItem{}
	item.Value = data
	item.StateVersion = _VERSION_1 // storage version
	srvc.CacheDB.Put(key, item.ToArray())
	return rec, nil
}

func getRecovery(srvc *native.NativeService, encID []byte) (*Group, error) {
	key := append(encID, FIELD_RECOVERY)
	item, err := utils.GetStorageItem(srvc, key)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.New("empty storage item")
	}
	if item.StateVersion != _VERSION_1 {
		return nil, errors.New("unexpected storage version")
	}
	return deserializeGroup(item.Value)
}

// deprecated
// retain for conpatibility
func addRecovery(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add recovery failed: argument 0 error")
	}
	// arg1: recovery address
	arg1, err := utils.DecodeAddress(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add recovery failed: argument 1 error")
	}
	// arg2: operator's public key
	arg2, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add recovery failed: argument 2 error")
	}

	err = checkWitness(srvc, arg2)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add recovery failed: " + err.Error())
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add recovery failed: " + err.Error())
	}
	if !isValid(srvc, key) {
		return utils.BYTE_FALSE, errors.New("add recovery failed: ID not registered")
	}
	if !isOwner(srvc, key, arg2) {
		return utils.BYTE_FALSE, errors.New("add recovery failed: not authorized")
	}

	re, err := getOldRecovery(srvc, key)
	if err == nil && len(re) > 0 {
		return utils.BYTE_FALSE, errors.New("add recovery failed: already set recovery")
	}

	err = setOldRecovery(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add recovery failed: " + err.Error())
	}

	triggerRecoveryEvent(srvc, "add", arg0, arg1)

	return utils.BYTE_TRUE, nil
}

// deprecated
// retain for conpatibility
func changeRecovery(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("change recovery failed: argument 0 error")
	}
	// arg1: new recovery address
	arg1, err := utils.DecodeAddress(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("change recovery failed: argument 1 error")
	}
	// arg2: operator's address, who should be the old recovery
	arg2, err := utils.DecodeAddress(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("change recovery failed: argument 2 error")
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("change recovery failed: " + err.Error())
	}
	re, err := getOldRecovery(srvc, key)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("change recovery failed: recovery not set")
	}
	if !bytes.Equal(re, arg2[:]) {
		return utils.BYTE_FALSE, errors.New("change recovery failed: operator is not the recovery")
	}
	err = checkWitness(srvc, arg2[:])
	if err != nil {
		return utils.BYTE_FALSE, errors.New("change recovery failed: " + err.Error())
	}
	if !isValid(srvc, key) {
		return utils.BYTE_FALSE, errors.New("change recovery failed: ID not registered")
	}
	err = setOldRecovery(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("change recovery failed: " + err.Error())
	}

	triggerRecoveryEvent(srvc, "change", arg0, arg1)
	return utils.BYTE_TRUE, nil
}

// deprecated
// retain for conpatibility
func setOldRecovery(srvc *native.NativeService, encID []byte, recovery common.Address) error {
	key := append(encID, FIELD_RECOVERY)
	val := states.StorageItem{Value: recovery[:]}
	val.StateVersion = _VERSION_0
	srvc.CacheDB.Put(key, val.ToArray())
	return nil
}

// deprecated
// retain for conpatibility
func getOldRecovery(srvc *native.NativeService, encID []byte) ([]byte, error) {
	key := append(encID, FIELD_RECOVERY)
	item, err := utils.GetStorageItem(srvc, key)
	if err != nil {
		return nil, errors.New("get recovery error: " + err.Error())
	} else if item == nil {
		return nil, nil
	}
	if item.StateVersion != _VERSION_0 {
		return nil, errors.New("unexpected storage version")
	}
	return item.Value, nil
}
