package event

import (
	"github.com/dad-go/common"
	"github.com/dad-go/vm/neovm/types"
)

type NotifyEventArgs struct {
	Container common.Uint256
	CodeHash common.Uint160
	States types.StackItemInterface
}

type NotifyEventInfo struct {
	Container common.Uint256
	CodeHash common.Uint160
	States interface{}
}

