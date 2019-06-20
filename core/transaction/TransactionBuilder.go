package transaction


import (
	"dad-go/core/asset"
	"dad-go/common"
	"dad-go/crypto"
	"dad-go/core/transaction/payload"
	"dad-go/core/contract/program"
)

//initial a new transaction with asset registration payload
func NewAssetRegistrationTransaction(asset *asset.Asset,amount *common.Fixed64,precision byte,issuer crypto.PubKey,conroller *common.Uint160) (*Transaction, error){

	//TODO: check arguments

	assetRegPayload := &payload.AssetRegistration {
		Asset: asset,
		Amount: amount,
		Precision: precision,
		Issuer: issuer,
		Controller: conroller,
	}

	return &Transaction{
		//nonce uint64 //TODO: genenrate nonce
		UTXOInputs: []*UTXOTxInput{},
		BalanceInputs: []*BalanceTxInput{},
		Attributes: []*TxAttribute{},
		TxType: RegisterAsset,
		Payload: assetRegPayload,
		Programs: []*program.Program{},
	}, nil
}


