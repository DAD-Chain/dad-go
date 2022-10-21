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
package testsuite

import (
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/smartcontract/service/native"
	_ "github.com/ontio/dad-go/smartcontract/service/native/init"
	"github.com/ontio/dad-go/smartcontract/service/native/ont"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
	"github.com/ontio/dad-go/smartcontract/storage"
	"github.com/stretchr/testify/assert"

	"testing"
)

func setOntBalance(db *storage.CacheDB, addr common.Address, value uint64) {
	balanceKey := ont.GenBalanceKey(utils.OntContractAddress, addr)
	item := utils.GenUInt64StorageItem(value)
	db.Put(balanceKey, item.ToArray())
}

func ontBalanceOf(native *native.NativeService, addr common.Address) int {
	sink := common.NewZeroCopySink(nil)
	utils.EncodeAddress(sink, addr)
	native.Input = sink.Bytes()
	buf, _ := ont.OntBalanceOf(native)
	val := common.BigIntFromNeoBytes(buf)
	return int(val.Uint64())
}

func ontTransfer(native *native.NativeService, from, to common.Address, value uint64) error {
	native.Tx.SignedAddr = append(native.Tx.SignedAddr, from)

	state := ont.State{from, to, value}
	native.Input = common.SerializeToBytes(&ont.Transfers{States: []ont.State{state}})

	_, err := ont.OntTransfer(native)
	return err
}

func TestTransfer(t *testing.T) {
	InvokeNativeContract(t, utils.OntContractAddress, func(native *native.NativeService) ([]byte, error) {
		a := RandomAddress()
		b := RandomAddress()
		c := RandomAddress()
		setOntBalance(native.CacheDB, a, 10000)

		assert.Equal(t, ontBalanceOf(native, a), 10000)
		assert.Equal(t, ontBalanceOf(native, b), 0)
		assert.Equal(t, ontBalanceOf(native, c), 0)

		assert.Nil(t, ontTransfer(native, a, b, 10))
		assert.Equal(t, ontBalanceOf(native, a), 9990)
		assert.Equal(t, ontBalanceOf(native, b), 10)

		assert.Nil(t, ontTransfer(native, b, c, 10))
		assert.Equal(t, ontBalanceOf(native, b), 0)
		assert.Equal(t, ontBalanceOf(native, c), 10)

		return nil, nil
	})
}
