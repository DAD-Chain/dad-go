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

package rpc

import (
	"github.com/dad-go/common/log"
	. "github.com/dad-go/http/base/common"
	. "github.com/dad-go/http/base/actor"
	Err "github.com/dad-go/http/base/error"
	"os"
	"path/filepath"
)

const (
	RANDBYTELEN = 4
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func GetNeighbor(params []interface{}) map[string]interface{} {
	addr, _ := GetNeighborAddrs()
	return responseSuccess(addr)
}

func GetNodeState(params []interface{}) map[string]interface{} {
	state,err := GetConnectionState()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	t,err := GetNodeTime()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	port,err := GetNodePort()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	id,err := GetID()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	ver,err := GetVersion()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	tpe,err := GetNodeType()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	relay,err := GetRelayState()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	height,err := BlockHeight()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	txnCnt,err := GetTxnCnt()
	if err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	n := NodeInfo{
		NodeState:    uint(state),
		NodeTime:     t,
		NodePort:     port,
		ID:       id,
		NodeVersion:  ver,
		NodeType: tpe,
		Relay:    relay,
		Height:   height,
		TxnCnt:   txnCnt,
	}
	return responseSuccess(n)
}

func StartConsensus(params []interface{}) map[string]interface{} {
	if err := ConsensusSrvStart(); err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	return responsePacking(Err.SUCCESS, true)
}

func StopConsensus(params []interface{}) map[string]interface{} {
	if err := ConsensusSrvHalt(); err != nil {
		return responsePacking(Err.INTERNAL_ERROR, false)
	}
	return responsePacking(Err.SUCCESS, true)
}

func SendSampleTransaction(params []interface{}) map[string]interface{} {
	panic("need reimplementation")
	return nil

	/*
		if len(params) < 1 {
			return dad-goRpcNil
		}
		var txType string
		switch params[0].(type) {
		case string:
			txType = params[0].(string)
		default:
			return dad-goRpcInvalidParameter
		}

		issuer, err := account.NewAccount()
		if err != nil {
			return dad-goRpc("Failed to create account")
		}
		admin := issuer

		rbuf := make([]byte, RANDBYTELEN)
		rand.Read(rbuf)
		switch string(txType) {
		case "perf":
			num := 1
			if len(params) == 2 {
				switch params[1].(type) {
				case float64:
					num = int(params[1].(float64))
				}
			}
			for i := 0; i < num; i++ {
				regTx := NewRegTx(ToHexString(rbuf), i, admin, issuer)
				SignTx(admin, regTx)
				VerifyAndSendTx(regTx)
			}
			return dad-goRpc(fmt.Sprintf("%d transaction(s) was sent", num))
		default:
			return dad-goRpc("Invalid transacion type")
		}
	*/
}

func SetDebugInfo(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return responsePacking(Err.INVALID_PARAMS, "")
	}
	switch params[0].(type) {
	case float64:
		level := params[0].(float64)
		if err := log.Log.SetDebugLevel(int(level)); err != nil {
			return responsePacking(Err.INVALID_PARAMS, "")
		}
	default:
		return responsePacking(Err.INVALID_PARAMS, "")
	}
	return responsePacking(Err.SUCCESS, true)
}
