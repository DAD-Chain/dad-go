package context

import (
	"github.com/dad-go/common"
	vmtypes "github.com/dad-go/vm/types"
)

type ContextRef interface {
	LoadContext(context *Context)
	CurrentContext() *Context
	CallingContext() *Context
	EntryContext() *Context
	Execute() error
}


type Context struct {
	ContractAddress common.Address
	Code vmtypes.VmCode
}
