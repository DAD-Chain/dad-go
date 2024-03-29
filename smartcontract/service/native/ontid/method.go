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
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/account"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/smartcontract/service/native"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
)

func regIdWithPublicKey(srvc *native.NativeService) ([]byte, error) {
	log.Debug("registerIdWithPublicKey")
	log.Debug("srvc.Input:", srvc.Input)
	// parse arguments
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: parsing argument 0 failed")
	} else if len(arg0) == 0 {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: invalid length of argument 0")
	}
	log.Debug("arg 0:", hex.EncodeToString(arg0), string(arg0))
	// arg1: public key
	arg1, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: parsing argument 1 failed")
	}

	log.Debug("arg 1:", hex.EncodeToString(arg1))

	if len(arg0) == 0 || len(arg1) == 0 {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: invalid argument")
	}

	if !account.VerifyID(string(arg0)) {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: invalid ID")
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: " + err.Error())
	}

	if checkIDState(srvc, key) != flag_not_exist {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: already registered")
	}

	public, err := keypair.DeserializePublicKey(arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: invalid public key")
	}
	addr := types.AddressFromPubKey(public)
	if !srvc.ContextRef.CheckWitness(addr) {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: checking witness failed")
	}

	// insert public key
	_, err = insertPk(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: store public key error, " + err.Error())
	}
	// set flags
	utils.PutBytes(srvc, key, []byte{flag_valid})

	triggerRegisterEvent(srvc, arg0)

	return utils.BYTE_TRUE, nil
}

func regIdWithAttributes(srvc *native.NativeService) ([]byte, error) {
	// parse arguments
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(source)

	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: argument 0 error, " + err.Error())
	} else if len(arg0) == 0 {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: argument 0 error, invalid length")
	}
	if !account.VerifyID(string(arg0)) {
		return utils.BYTE_FALSE, errors.New("register ONT ID error: invalid ID")
	}

	// arg1: public key
	arg1, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: argument 1 error," + err.Error())
	}
	// arg2: attributes
	// first get number
	num, err := utils.DecodeVarUint(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: argument 2 error, " + err.Error())
	}

	// next parse each attribute
	var arg2 = make([]attribute, 0)
	for i := 0; i < int(num); i++ {
		var v attribute
		err = v.Deserialization(source)
		if err != nil {
			return utils.BYTE_FALSE, errors.New("register ID with attributes error: argument 2 error, " + err.Error())
		}
		arg2 = append(arg2, v)
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: " + err.Error())
	}

	if checkIDState(srvc, key) != flag_not_exist {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: already registered")
	}
	public, err := keypair.DeserializePublicKey(arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: invalid public key: " + err.Error())
	}
	addr := types.AddressFromPubKey(public)
	if !srvc.ContextRef.CheckWitness(addr) {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: check witness failed")
	}

	_, err = insertPk(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: store pubic key error: " + err.Error())
	}

	err = batchInsertAttr(srvc, key, arg2)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("register ID with attributes error: insert attribute error: " + err.Error())
	}

	utils.PutBytes(srvc, key, []byte{flag_valid})
	triggerRegisterEvent(srvc, arg0)
	return utils.BYTE_TRUE, nil
}

func addKey(srvc *native.NativeService) ([]byte, error) {
	log.Debug("ID contract: AddKey")
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: id
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add key failed: argument 0 error, " + err.Error())
	}
	log.Debug("arg 0:", hex.EncodeToString(arg0))

	// arg1: public key
	arg1, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add key failed: argument 1 error, " + err.Error())
	}
	log.Debug("arg 1:", hex.EncodeToString(arg1))
	_, err = keypair.DeserializePublicKey(arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add key error: invalid key")
	}

	// arg2: operator's public key
	arg2, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add key failed: argument 2 error, " + err.Error())
	}
	log.Debug("arg 2:", hex.EncodeToString(arg2))

	if err = checkWitness(srvc, arg2); err != nil {
		return utils.BYTE_FALSE, errors.New("add key failed: check witness failed, " + err.Error())
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add key failed: " + err.Error())
	}
	if !isValid(srvc, key) {
		return utils.BYTE_FALSE, errors.New("add key failed: ID not registered")
	}
	var auth bool = false
	rec, _ := getOldRecovery(srvc, key)
	if len(rec) > 0 {
		auth = bytes.Equal(rec, arg2)
	}
	if !auth {
		if !isOwner(srvc, key, arg2) {
			return utils.BYTE_FALSE, errors.New("add key failed: operator has no authorization")
		}
	}

	keyID, err := insertPk(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("add key failed: insert public key error, " + err.Error())
	}

	triggerPublicEvent(srvc, "add", arg0, arg1, keyID)

	return utils.BYTE_TRUE, nil
}

func removeKey(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: id
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("remove key failed: argument 0 error, %s", err)
	}

	// arg1: public key
	arg1, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("remove key failed: argument 1 error, %s", err)
	}

	// arg2: operator's public key
	arg2, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("remove key failed: argument 2 error, %s", err)
	}
	if err = checkWitness(srvc, arg2); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("remove key failed: check witness failed, %s", err)
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("remove key failed: %s", err)
	}
	if !isValid(srvc, key) {
		return utils.BYTE_FALSE, errors.New("remove key failed: ID not registered")
	}
	var auth = false
	rec, err := getOldRecovery(srvc, key)
	if len(rec) > 0 {
		auth = bytes.Equal(rec, arg2)
	}
	if !auth {
		if !isOwner(srvc, key, arg2) {
			return utils.BYTE_FALSE, errors.New("remove key failed: operator has no authorization")
		}
	}

	keyID, err := revokePk(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("remove key failed: %s", err)
	}

	triggerPublicEvent(srvc, "remove", arg0, arg1, keyID)

	return utils.BYTE_TRUE, nil
}

func addAttributes(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("add attributes failed, argument 0, error: %s", err)
	}
	// arg1: attributes
	// first get number
	num, err := utils.DecodeVarUint(source)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("add attributes failed, argument 1 error: %s", err)
	}
	// next parse each attribute
	var arg1 = make([]attribute, 0)
	for i := 0; i < int(num); i++ {
		var v attribute
		err = v.Deserialization(source)
		if err != nil {
			return utils.BYTE_FALSE, fmt.Errorf("add attributes failed, argument 1 error: %s", err)
		}
		arg1 = append(arg1, v)
	}
	// arg2: opperator's public key
	arg2, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("add attributes failed, argument 2, error: %s", err)
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("add attributes failed: %s", err)
	}
	if !isValid(srvc, key) {
		return utils.BYTE_FALSE, errors.New("add attributes failed, ID not registered")
	}
	if !isOwner(srvc, key, arg2) {
		return utils.BYTE_FALSE, errors.New("add attributes failed, no authorization")
	}
	err = checkWitness(srvc, arg2)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("add attributes failed, %s", err)
	}

	err = batchInsertAttr(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("add attributes failed, %s", err)
	}

	paths := getAttrKeys(arg1)
	triggerAttributeEvent(srvc, "add", arg0, paths)
	return utils.BYTE_TRUE, nil
}

func removeAttribute(srvc *native.NativeService) ([]byte, error) {
	args := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("remove attribute failed: argument 0 error")
	}
	// arg1: path
	arg1, err := utils.DecodeVarBytes(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("remove attribute failed: argument 1 error")
	}
	// arg2: operator's public key
	arg2, err := utils.DecodeVarBytes(args)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("remove attribute failed: argument 2 error")
	}

	err = checkWitness(srvc, arg2)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("remove attribute failed: " + err.Error())
	}
	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("remove attribute failed: " + err.Error())
	}
	if !isValid(srvc, key) {
		return utils.BYTE_FALSE, errors.New("remove attribute failed: ID not registered")
	}
	if !isOwner(srvc, key, arg2) {
		return utils.BYTE_FALSE, errors.New("remove attribute failed: no authorization")
	}

	err = deleteAttr(srvc, key, arg1)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("remove attribute failed: " + err.Error())
	}

	triggerAttributeEvent(srvc, "remove", arg0, [][]byte{arg1})
	return utils.BYTE_TRUE, nil
}

func verifySignature(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: ID
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("verify signature error: argument 0 error, error: " + err.Error())
	}
	// arg1: index of public key
	arg1, err := utils.DecodeVarUint(source)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("verify signature error: argument 1 error, " + err.Error())
	}

	key, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, errors.New("verify signature error: " + err.Error())
	}
	if err := checkWitnessByIndex(srvc, key, uint32(arg1)); err != nil {
		return utils.BYTE_FALSE, errors.New("verify signature failed: " + err.Error())
	}

	return utils.BYTE_TRUE, nil
}

func revokeID(srvc *native.NativeService) ([]byte, error) {
	source := common.NewZeroCopySource(srvc.Input)
	// arg0: id
	arg0, err := utils.DecodeVarBytes(source)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("argument 0 error")
	}
	// arg1: index of public key
	arg1, err := utils.DecodeVarUint(source)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("argument 1 error")
	}

	encID, err := encodeID(arg0)
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	if !isValid(srvc, encID) {
		return utils.BYTE_FALSE, fmt.Errorf("%s is not registered or already revoked", string(arg0))
	}

	if err := checkWitnessByIndex(srvc, encID, uint32(arg1)); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("authorization failed, %s", err)
	}

	err = deleteID(srvc, encID)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("delete id error, %s", err)
	}
	newEvent(srvc, []interface{}{"Revoke", string(arg0)})
	return utils.BYTE_TRUE, nil
}
