package types

import (
	"math/big"
	"github.com/dad-go/vm/neovm/interfaces"
)

type StackItemInterface interface {
	Equals(other StackItemInterface) bool
	GetBigInteger() *big.Int
	GetBoolean() bool
	GetByteArray() []byte
	GetInterface() interfaces.IInteropInterface
	GetArray() []StackItemInterface
	GetStruct() []StackItemInterface
}

