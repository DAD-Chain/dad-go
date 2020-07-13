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
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/payload"
	"github.com/ontio/dad-go/core/types"
)

var (
	txn *types.Transaction
)

func init() {
	log.Init(log.PATH, log.Stdout)

	bookKeepingPayload := &payload.Bookkeeping{
		Nonce: uint64(time.Now().UnixNano()),
	}

	txn = &types.Transaction{
		Version:    0,
		Attributes: []*types.TxAttribute{},
		TxType:     types.Bookkeeper,
		Payload:    bookKeepingPayload,
	}

	tempStr := "3369930accc1ddd067245e8edadcd9bea207ba5e1753ac18a51df77a343bfe92"
	hex, _ := hex.DecodeString(tempStr)
	var hash common.Uint256
	hash.Deserialize(bytes.NewReader(hex))
	txn.SetHash(hash)
}

func TestTxPool(t *testing.T) {
	txPool := &TXPool{}
	txPool.Init()

	txEntry := &TXEntry{
		Tx:    txn,
		Attrs: []*TXAttr{},
		Fee:   txn.GetTotalFee(),
	}

	ret := txPool.AddTxList(txEntry)
	if ret == false {
		t.Error("Failed to add tx to the pool")
		return
	}

	ret = txPool.AddTxList(txEntry)
	if ret == true {
		t.Error("Failed to add tx to the pool")
		return
	}

	txList := txPool.GetTxPool(true)
	for _, v := range txList {
		fmt.Println(v)
	}

	entry := txPool.GetTransaction(txn.Hash())
	if entry == nil {
		t.Error("Failed to get the transaction")
		return
	}

	status := txPool.GetTxStatus(txn.Hash())
	if status == nil {
		t.Error("failed to get the status")
		return
	}

	count := txPool.GetTransactionCount()
	fmt.Println(count)

	err := txPool.CleanTransactionList([]*types.Transaction{txn})
	if err != nil {
		t.Error("Failed to clean transaction list")
		return
	}
}
