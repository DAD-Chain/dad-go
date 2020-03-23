package event

import (
	"github.com/dad-go/vm/neovm/interfaces"
	"github.com/dad-go/common"
)

type LogEventArgs struct {
	container interfaces.ICodeContainer
	codeHash  common.Uint160
	message   string
}