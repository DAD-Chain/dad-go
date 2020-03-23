package event

import (
	"github.com/dad-go/common"
)

type LogEventArgs struct {
	Container common.Uint256
	CodeHash  common.Uint160
	Message   string
}