package contract

import (
	"dad-go/crypto"
	. "dad-go/common"
)

//create a Single Singature contract for owner  。
func CreateSignatureContract(ownerPubKey crypto.PubKey) (*Contract,error){
	//TODO: implement func CreateSignatureContract
	return nil,nil
}

//create a Multi Singature contract for owner  。
func CreateMultiSigContract(publicKeyHash Uint160,m int, publicKeys ...[]*crypto.PubKey) (*Contract,error){
	//TODO: implement func CreateSignatureContract
	return nil,nil
}
