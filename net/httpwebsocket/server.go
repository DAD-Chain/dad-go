package httpwebsocket

import (
	. "dad-go/common"
	. "dad-go/common/config"
	"dad-go/core/ledger"
	"dad-go/events"
	"dad-go/net/httprestful/common"
	Err "dad-go/net/httprestful/error"
	"dad-go/net/httpwebsocket/websocket"
	. "dad-go/net/protocol"
)

var ws *websocket.WsServer
var pushBlockFlag bool = false

func StartServer(n Noder) {
	common.SetNode(n)
	ledger.DefaultLedger.Blockchain.BCEvents.Subscribe(events.EventBlockPersistCompleted, SendBlock2WSclient)
	go func() {
		ws = websocket.InitWsServer(common.CheckAccessToken)
		ws.Start()
	}()
}
func SendBlock2WSclient(v interface{}) {
	if Parameters.HttpWsPort != 0 && pushBlockFlag {
		go func() {
			PushBlock(v)
		}()
	}
}
func Stop() {
	if ws != nil {
		ws.Stop()
	}
}
func ReStartServer() {
	if ws == nil {
		ws = websocket.InitWsServer(common.CheckAccessToken)
		ws.Start()
		return
	}
	ws.Restart()
}
func GetWsPushBlockFlag() bool {
	return pushBlockFlag
}
func SetWsPushBlockFlag(b bool) {
	pushBlockFlag = b
}
func SetTxHashMap(txhash string, sessionid string) {
	if ws != nil {
		ws.SetTxHashMap(txhash, sessionid)
	}
}
func PushSmartCodeInvokeResult(txHash Uint256, errcode int64, result interface{}) {
	if ws != nil {
		resp := common.ResponsePack(Err.SUCCESS)
		var Result = make(map[string]interface{})
		txHashStr := ToHexString(txHash.ToArray())
		Result["TxHash"] = txHashStr
		Result["ExecResult"] = result

		resp["Result"] = Result
		resp["Action"] = "sendsmartcodeinvoke"
		resp["Error"] = errcode
		resp["Desc"] = Err.ErrMap[errcode]
		ws.PushTxResult(txHashStr, resp)
	}
}
func PushBlock(v interface{}) {
	if ws != nil {
		resp := common.ResponsePack(Err.SUCCESS)
		if block, ok := v.(*ledger.Block); ok {
			resp["Result"] = common.GetBlockInfo(block)
			resp["Action"] = "sendrawblock"
			ws.PushResult(resp)
		}
	}
}
