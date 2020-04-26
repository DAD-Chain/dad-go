package httpjsonrpc

import (
	. "github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	. "github.com/dad-go/http/common/rpc"
	. "github.com/dad-go/http/common"
	. "github.com/dad-go/consensus"
	. "github.com/dad-go/net/protocol"
	"net/http"
	"strconv"
)

func RegistRpcNode(n Noder) {
	if CNoder == nil {
		CNoder = n
	}
}

func RegistConsensusService(consensus ConsensusService) {
	ConsensusSrv = consensus
}
func StartRPCServer() {
	log.Debug()
	http.HandleFunc("/", Handle)

	HandleFunc("getbestblockhash", GetBestBlockHash)
	HandleFunc("getblock", GetBlock)
	HandleFunc("getblockcount", GetBlockCount)
	HandleFunc("getblockhash", GetBlockHash)
	//HandleFunc("getunspendoutput", getUnspendOutput)
	HandleFunc("getconnectioncount", GetConnectionCount)
	HandleFunc("getrawmempool", GetRawMemPool)
	HandleFunc("getrawtransaction", GetRawTransaction)
	HandleFunc("sendrawtransaction", SendRawTransaction)
	HandleFunc("getstorage", GetStorage)
	HandleFunc("getbalance", GetBalance)
	HandleFunc("submitblock", SubmitBlock)
	HandleFunc("getversion", GetVersion)
	HandleFunc("getdataile", GetDataFile)
	HandleFunc("catdatarecord", CatDataRecord)
	HandleFunc("regdatafile", RegDataFile)
	HandleFunc("uploadDataFile", UploadDataFile)
	HandleFunc("getsmartcodeevent", GetSmartCodeEvent)

	err := http.ListenAndServe(":"+strconv.Itoa(Parameters.HttpJsonPort), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}
}
