package code

import (
	. "dad-go/common"
	. "dad-go/core/contract"
)
//ICode is the abstract interface of smart contract code.
type ICode interface {

	GetCode() []byte

	GetParameterTypes() []ContractParameterType

	GetReturnTypes() []ContractParameterType

	CodeHash() Uint160

}

