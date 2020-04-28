package jsonrpc

import (
	. "github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	. "github.com/dad-go/http/base/rpc"
	//. "github.com/dad-go/http/base/common"
	//. "github.com/dad-go/consensus"
	"net/http"
	"strconv"
)




func StartRPCServer() {
	log.Debug()
	http.HandleFunc("/", Handle)

	HandleFunc("getbestblockhash", getBestBlockHash)
	HandleFunc("getblock", getBlock)
	HandleFunc("getblockcount", getBlockCount)
	HandleFunc("getblockhash", getBlockHash)
	//HandleFunc("getunspendoutput", getUnspendOutput)
	HandleFunc("getconnectioncount", getConnectionCount)
	HandleFunc("getrawmempool", getRawMemPool)
	HandleFunc("getmempooltx", getMemPoolTx)
	HandleFunc("getrawtransaction", getRawTransaction)
	HandleFunc("sendrawtransaction", sendRawTransaction)
	HandleFunc("getstorage", getStorage)
	HandleFunc("getbalance", getBalance)
	HandleFunc("submitblock", submitBlock)
	HandleFunc("getversion", getNodeVersion)
	HandleFunc("getdataile", getDataFile)
	HandleFunc("catdatarecord", catDataRecord)
	HandleFunc("regdatafile", regDataFile)
	HandleFunc("uploadDataFile", uploadDataFile)
	HandleFunc("getsmartcodeevent", getSmartCodeEvent)

	err := http.ListenAndServe(":"+strconv.Itoa(Parameters.HttpJsonPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
