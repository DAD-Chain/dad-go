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

package jsonrpc

import (
	"net/http"
	"strconv"

	cfg "github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	"github.com/dad-go/http/base/rpc"
)

func StartRPCServer() {
	log.Debug()
	http.HandleFunc("/", rpc.Handle)

	rpc.HandleFunc("getgenerateblocktime", rpc.GetGenerateBlockTime)
	rpc.HandleFunc("getbestblockhash", rpc.GetBestBlockHash)
	rpc.HandleFunc("getblock", rpc.GetBlock)
	rpc.HandleFunc("getblockcount", rpc.GetBlockCount)
	rpc.HandleFunc("getblockhash", rpc.GetBlockHash)
	rpc.HandleFunc("getconnectioncount", rpc.GetConnectionCount)
	//HandleFunc("getrawmempool", GetRawMemPool)

	rpc.HandleFunc("getrawtransaction", rpc.GetRawTransaction)
	rpc.HandleFunc("sendrawtransaction", rpc.SendRawTransaction)
	rpc.HandleFunc("getstorage", rpc.GetStorage)
	rpc.HandleFunc("getversion", rpc.GetNodeVersion)

	rpc.HandleFunc("getblocksysfee", rpc.GetSystemFee)
	rpc.HandleFunc("getcontractstate", rpc.GetContractState)
	rpc.HandleFunc("getmempooltxstate", rpc.GetMemPoolTxState)
	rpc.HandleFunc("getsmartcodeevent", rpc.GetSmartCodeEvent)
	rpc.HandleFunc("getblockheightbytxhash", rpc.GetBlockHeightByTxHash)

	rpc.HandleFunc("getbalance", rpc.GetBalance)

	err := http.ListenAndServe(":"+strconv.Itoa(cfg.Parameters.HttpJsonPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
