package code

import (
	. "github.com/dad-go/common"
	"github.com/dad-go/common/serialization"
	. "github.com/dad-go/core/contract"
	. "github.com/dad-go/errors"
	"io"
)

type FunctionCode struct {
	// Contract Code
	Code []byte

	// Contract parameter type list
	ParameterTypes []ContractParameterType

	// Contract return type
	ReturnType ContractParameterType

	codeHash Uint160
}

// method of SerializableData
func (fc *FunctionCode) Serialize(w io.Writer) error {
	var err error
	err = serialization.WriteVarBytes(w, fc.Code)
	if err != nil {
		return err
	}

	err = serialization.WriteVarBytes(w, ContractParameterTypeToByte(fc.ParameterTypes))
	if err != nil {
		return err
	}

	err = serialization.WriteByte(w, byte(fc.ReturnType))
	if err != nil {
		return err
	}

	return nil
}

// method of SerializableData
func (fc *FunctionCode) Deserialize(r io.Reader) error {
	var err error

	fc.Code, err = serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction FunctionCode Code Deserialize failed.")
	}

	p, err := serialization.ReadVarBytes(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction FunctionCode ParameterTypes Deserialize failed.")
	}
	fc.ParameterTypes = ByteToContractParameterType(p)

	returnType, err := serialization.ReadByte(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Transaction FunctionCode returnType Deserialize failed.")
	}
	fc.ReturnType = ContractParameterType(returnType)
	return nil
}

// method of ICode
// Get the hash of the smart contract
func (fc *FunctionCode) CodeHash() Uint160 {
	u160 := Uint160{}
	if fc.codeHash == u160 {
		u160 := ToCodeHash(fc.Code)
		fc.codeHash = u160
	}
	return fc.codeHash
}
