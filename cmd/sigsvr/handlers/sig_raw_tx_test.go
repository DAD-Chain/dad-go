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

package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"github.com/ontio/dad-go/account"
	clisvrcom "github.com/ontio/dad-go/cmd/server/common"
	"github.com/ontio/dad-go/cmd/utils"
	"github.com/ontio/dad-go/common/log"
	"os"
	"testing"
)

var (
	wallet *account.ClientImpl
	passwd = []byte("123456")
)

func TestMain(m *testing.M) {
	log.InitLog(0, os.Stdout)
	clisvrcom.DefAccount = account.NewAccount("")
	m.Run()
	os.RemoveAll("./ActorLog")
}

func TestSigRawTx(t *testing.T) {
	acc := account.NewAccount("")
	defAcc := clisvrcom.DefAccount
	tx, err := utils.TransferTx(0, 0, "ont", defAcc.Address.ToBase58(), acc.Address.ToBase58(), 10)
	if err != nil {
		t.Errorf("TransferTx error:%s", err)
		return
	}
	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	if err != nil {
		t.Errorf("tx.Serialize error:%s", err)
		return
	}
	rawReq := &SigRawTransactionReq{
		RawTx: hex.EncodeToString(buf.Bytes()),
	}
	data, err := json.Marshal(rawReq)
	if err != nil {
		t.Errorf("json.Marshal SigRawTransactionReq error:%s", err)
		return
	}
	req := &clisvrcom.CliRpcRequest{
		Qid:    "t",
		Method: "sigrawtx",
		Params: data,
	}
	resp := &clisvrcom.CliRpcResponse{}
	SigRawTransaction(req, resp)
	if resp.ErrorCode != 0 {
		t.Errorf("SigRawTransaction failed. ErrorCode:%d", resp.ErrorCode)
		return
	}
}
