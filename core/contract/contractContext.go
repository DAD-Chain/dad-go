package contract

 import (
	"dad-go/common"
	sig "dad-go/core/signature"
)

 type ContractContext struct {
	//TODO: define ContractContext。
	Data sig.SignableData
	ProgramHashes []common.Uint160
	Programs [][]byte
	Parameters [][][]byte
}


 func NewContractContext(data sig.SignableData) *ContractContext {
	//TODO: implement NewContractContext
	return nil
} 
