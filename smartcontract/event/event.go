package event

import (
	"github.com/dad-go/events"
	"github.com/dad-go/core/ledger"
	. "github.com/dad-go/common"
)

func PushSmartCodeEvent(txHash Uint256, errcode int64, action string, result interface{}) {
	resp := map[string]interface{}{
		"TxHash": ToHexString(txHash.ToArray()),
		"Action": action,
		"Result": result,
		"Error":  errcode,
	}
	ledger.DefaultLedger.Blockchain.BCEvents.Notify(events.EventSmartCode, resp)
}
