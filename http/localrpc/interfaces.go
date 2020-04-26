package localrpc

import (
	"github.com/dad-go/common/log"
	. "github.com/dad-go/http/base/common"
	. "github.com/dad-go/http/base/rpc"
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
	addr, _ := CNoder.GetNeighborAddrs()
	return dad-goRpc(addr)
}

func GetNodeState(params []interface{}) map[string]interface{} {
	n := NodeInfo{
		State:    uint(CNoder.GetState()),
		Time:     CNoder.GetTime(),
		Port:     CNoder.GetPort(),
		ID:       CNoder.GetID(),
		Version:  CNoder.Version(),
		Services: CNoder.Services(),
		Relay:    CNoder.GetRelay(),
		Height:   CNoder.GetHeight(),
		TxnCnt:   CNoder.GetTxnCnt(),
		RxTxnCnt: CNoder.GetRxTxnCnt(),
	}
	return dad-goRpc(n)
}

func StartConsensus(params []interface{}) map[string]interface{} {
	if err := ConsensusSrv.Start(); err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpcSuccess
}

func StopConsensus(params []interface{}) map[string]interface{} {
	if err := ConsensusSrv.Halt(); err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpcSuccess
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
		return dad-goRpcInvalidParameter
	}
	switch params[0].(type) {
	case float64:
		level := params[0].(float64)
		if err := log.Log.SetDebugLevel(int(level)); err != nil {
			return dad-goRpcInvalidParameter
		}
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpcSuccess
}
