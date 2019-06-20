package signature

import (
	"dad-go/common"
	"dad-go/core/contract/program"
	"dad-go/vm"
)

//SignableData describe the data need be signed.
type SignableData interface {

	vm.ISignableObject

	//Get the the SignableData's program hashes
	GetProgramHashes() ([]common.Uint160, error)

	SetPrograms([]*program.Program)

	GetPrograms()  []*program.Program
}


func Sign(data SignableData,signer Signer) ([]byte, error){

	//TODO: implement Sign()
	return  []byte{},nil
}

func GetHashData(data SignableData) []byte {
	//TODO: implement GetHashData()

	return nil
}
