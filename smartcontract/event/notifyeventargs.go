package event

import (
	"github.com/dad-go/common"
	"github.com/dad-go/vm/neovm/types"
)

type NotifyEventArgs struct {
	Container common.Uint256
	CodeHash common.Address
	States types.StackItemInterface
}

type NotifyEventInfo struct {
	Container common.Uint256
	CodeHash common.Address
	States interface{}
}

