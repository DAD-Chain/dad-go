package event

import (
	. "github.com/dad-go/common"
	"github.com/dad-go/events"
	"github.com/dad-go/events/message"
	"github.com/dad-go/core/types"
)

func PushSmartCodeEvent(txHash Uint256, errcode int64, action string, result interface{}) {
	smartCodeEvt := &types.SmartCodeEvent{
		TxHash: ToHexString(txHash.ToArray()),
		Action: action,
		Result: result,
		Error:  errcode,
	}
	events.DefActorPublisher.Publish(message.TopicSmartCodeEvent, smartCodeEvt)
}
