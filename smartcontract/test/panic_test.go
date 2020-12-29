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

package test

import (
	"crypto/rand"
	"fmt"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/types"
	. "github.com/ontio/dad-go/smartcontract"
	"os"
	"testing"
)

func TestRandomCodeCrash(t *testing.T) {
	dbFile := "test"
	log.InitLog(4)
	defer func() {
		os.RemoveAll(dbFile)
	}()
	//testLevelDB, err := leveldbstore.NewLevelDBStore(dbFile)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//store := statestore.NewMemDatabase()
	//
	//testBatch := statestore.NewStateStoreBatch(store, testLevelDB)
	config := &Config{
		Time:   10,
		Height: 10,
		Tx:     &types.Transaction{},
	}

	var code []byte
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("code %x \n", code)
		}
	}()

	for i := 1; i < 10; i++ {
		fmt.Printf("test round:%d \n", i)
		code := make([]byte, i)
		for j := 0; j < 10; j++ {
			rand.Read(code)

			//cache := storage.NewCloneCache(testBatch)
			sc := SmartContract{
				Config:     config,
				Gas:        10000,
				CloneCache: nil,
			}
			engine, _ := sc.NewExecuteEngine(code)
			engine.Invoke()
		}
	}
}
