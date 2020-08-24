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
	"github.com/ontio/dad-go/common"
)

func TestSignMsg(t *testing.T) {
	passwd := string("passwordtest")
	acct := account.Open(account.WALLET_FILENAME, []byte(passwd))
	acc, err := acct.GetDefaultAccount()
	if err != nil {
		fmt.Println("GetDefaultAccount error:", err)
		os.Exit(1)
	}
	msg, err := constructProposalMsg(acc)
	if err != nil {
		t.Errorf("constructProposalMsg failed: %v", err)
		return
	}
	_, err = SignMsg(acc.PrivateKey, msg)
	if err != nil {
		t.Error("TestSignMsg Failed: %v", err)
		return
	}
	t.Log("TestSignMsg succ")
}

func TestHashBlock(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	hash, _ := HashBlock(blk)
	t.Logf("TestHashBlock: %v", hash)
}

func TestHashMsg(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
		return
	}
	blockproposalmsg := &blockProposalMsg{
		Block: blk,
	}
	uint256, err := HashMsg(blockproposalmsg)
	if err != nil {
		t.Errorf("TestHashMsg failed: %v", err)
		return
	}
	t.Logf("TestHashMsg succ: %v\n", uint256)
}

func TestVrf(t *testing.T) {
	blk, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed: %v", err)
	}
	hash := common.Uint256{}
	vrfvalue := vrf(blk, hash)
	if len(vrfvalue) == 0 {
		t.Errorf("TestVrf failed:")
		return
	}
	t.Log("TestVrf succ")
}