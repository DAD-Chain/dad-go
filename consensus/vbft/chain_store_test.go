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

package vbft

import (
	"fmt"
	"os"
	"testing"

	"github.com/ontio/dad-go/account"
	"github.com/ontio/dad-go/common/config"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/ledger"
)

func newChainStore() *ChainStore {
	log.Init(log.PATH, log.Stdout)
	var err error
	acct := account.NewAccount("SHA256withECDSA")
	if acct == nil {
		fmt.Println("GetDefaultAccount error: acc is nil")
		os.Exit(1)
	}

	ledger.DefLedger, err = ledger.NewLedger(config.DEFAULT_DATA_DIR)
	if err != nil {
		log.Fatalf("NewLedger error %s", err)
		os.Exit(1)
	}
	store, err := OpenBlockStore(ledger.DefLedger)
	if err != nil {
		fmt.Printf("openblockstore failed: %v\n", err)
		return nil
	}
	return store
}

func TestGetChainedBlockNum(t *testing.T) {
	chainstore := newChainStore()
	if chainstore == nil {
		t.Error("newChainStrore error")
		return
	}
	blocknum := chainstore.GetChainedBlockNum()
	t.Logf("TestGetChainedBlockNum :%d", blocknum)
}
