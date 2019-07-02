package code

import (
	."github.com/DAD-Chain/dad-go/common"
	."github.com/DAD-Chain/dad-go/core/contract"
)
//ICode is the abstract interface of smart contract code.
type ICode interface {

	GetCode() []byte

	GetParameterTypes() []ContractParameterType

	GetReturnTypes() []ContractParameterType

	CodeHash() Uint160

}

