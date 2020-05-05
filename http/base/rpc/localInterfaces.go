package localrpc

import (
	"github.com/dad-go/common/log"
	. "github.com/dad-go/http/base/common"
	. "github.com/dad-go/http/base/rpc"
	. "github.com/dad-go/http/base/actor"
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

func getNeighbor(params []interface{}) map[string]interface{} {
	addr, _ := GetNeighborAddrs()
	return dad-goRpc(addr)
}

func getNodeState(params []interface{}) map[string]interface{} {
	state,err := GetConnectionState()
	if err != nil {
		return dad-goRpcFailed
	}
	t,err := GetNodeTime()
	if err != nil {
		return dad-goRpcFailed
	}
	port,err := GetNodePort()
	if err != nil {
		return dad-goRpcFailed
	}
	id,err := GetID()
	if err != nil {
		return dad-goRpcFailed
	}
	ver,err := GetNodeVersion()
	if err != nil {
		return dad-goRpcFailed
	}
	tpe,err := GetNodeType()
	if err != nil {
		return dad-goRpcFailed
	}
	relay,err := GetRelayState()
	if err != nil {
		return dad-goRpcFailed
	}
	height,err := BlockHeight()
	if err != nil {
		return dad-goRpcFailed
	}
	txnCnt,err := GetTxnCnt()
	if err != nil {
		return dad-goRpcFailed
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
	return dad-goRpc(n)
}

func startConsensus(params []interface{}) map[string]interface{} {
	if err := ConsensusSrvStart(); err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpcSuccess
}

func stopConsensus(params []interface{}) map[string]interface{} {
	if err := ConsensusSrvHalt(); err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpcSuccess
}

func sendSampleTransaction(params []interface{}) map[string]interface{} {
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

func setDebugInfo(params []interface{}) map[string]interface{} {
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
