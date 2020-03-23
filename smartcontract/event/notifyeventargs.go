package event

import (
	"github.com/dad-go/vm/neovm/interfaces"
	"github.com/dad-go/common"
	"github.com/dad-go/vm/neovm/types"
)

type NotifyEventArgs struct {
	container interfaces.ICodeContainer
	codeHash  common.Uint160
	state     types.StackItemInterface
}

