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

package ledgerstore

import (
	"testing"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/account"
	"github.com/ontio/dad-go/core/payload"
	"github.com/ontio/dad-go/core/states"
	scommon "github.com/ontio/dad-go/core/store/common"
	"github.com/ontio/dad-go/core/store/statestore"
	vmtypes "github.com/ontio/dad-go/smartcontract/types"
)

func TestContractState(t *testing.T) {
	batch, err := getStateBatch()
	if err != nil {
		t.Errorf("NewStateBatch error %s", err)
		return
	}
	testCode := []byte("testcode")

	vmCode := vmtypes.VmCode{
		VmType: vmtypes.NEOVM,
		Code:   testCode,
	}
	deploy := &payload.DeployCode{
		Code:        vmCode,
		NeedStorage: false,
		Name:        "testsm",
		Version:     "v1.0",
		Author:      "",
		Email:       "",
		Description: "",
	}
	code := &vmtypes.VmCode{
		Code:   testCode,
		VmType: vmtypes.NEOVM,
	}
	codeHash := code.AddressFromVmCode()
	err = batch.TryGetOrAdd(
		scommon.ST_CONTRACT,
		codeHash[:],
		deploy)
	if err != nil {
		t.Errorf("TryGetOrAdd contract error %s", err)
		return
	}

	err = batch.CommitTo()
	if err != nil {
		t.Errorf("batch.CommitTo error %s", err)
		return
	}
	err = testStateStore.CommitTo()
	if err != nil {
		t.Errorf("testStateStore.CommitTo error %s", err)
		return
	}
	contractState1, err := testStateStore.GetContractState(codeHash)
	if err != nil {
		t.Errorf("GetContractState error %s", err)
		return
	}
	if contractState1.Name != deploy.Name ||
		contractState1.Version != deploy.Version ||
		contractState1.Author != deploy.Author ||
		contractState1.Description != deploy.Description ||
		contractState1.Email != deploy.Email {
		t.Errorf("TestContractState failed %+v != %+v", contractState1, deploy)
		return
	}
}

func TestBookkeeperState(t *testing.T) {
	batch, err := getStateBatch()
	if err != nil {
		t.Errorf("NewStateBatch error %s", err)
		return
	}
	acc1 := account.NewAccount("")
	acc2 := account.NewAccount("")

	currBookkeepers := make([]keypair.PublicKey, 0)
	currBookkeepers = append(currBookkeepers, acc1.PublicKey)
	currBookkeepers = append(currBookkeepers, acc2.PublicKey)
	nextBookkeepers := make([]keypair.PublicKey, 0)
	nextBookkeepers = append(nextBookkeepers, acc1.PublicKey)
	nextBookkeepers = append(nextBookkeepers, acc2.PublicKey)

	bookkeeperState := &states.BookkeeperState{
		CurrBookkeeper: currBookkeepers,
		NextBookkeeper: nextBookkeepers,
	}

	batch.TryAdd(scommon.ST_BOOKKEEPER, BOOKKEEPER, bookkeeperState)
	err = batch.CommitTo()
	if err != nil {
		t.Errorf("batch.CommitTo error %s", err)
		return
	}
	err = testStateStore.CommitTo()
	if err != nil {
		t.Errorf("testStateStore.CommitTo error %s", err)
		return
	}
	bookState, err := testStateStore.GetBookkeeperState()
	if err != nil {
		t.Errorf("GetBookkeeperState error %s", err)
		return
	}
	currBookkeepers1 := bookState.CurrBookkeeper
	nextBookkeepers1 := bookState.NextBookkeeper
	for index, pk := range currBookkeepers {
		pk1 := currBookkeepers1[index]
		if !keypair.ComparePublicKey(pk, pk1) {
			t.Errorf("TestBookkeeperState currentBookkeeper failed")
			return
		}
	}
	for index, pk := range nextBookkeepers {
		pk1 := nextBookkeepers1[index]
		if !keypair.ComparePublicKey(pk, pk1) {
			t.Errorf("TestBookkeeperState nextBookkeeper failed")
			return
		}
	}
}

func getStateBatch() (*statestore.StateBatch, error) {
	testStateStore.NewBatch()
	batch := testStateStore.NewStateBatch()
	return batch, nil
}
