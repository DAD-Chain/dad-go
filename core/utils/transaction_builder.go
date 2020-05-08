package utils

import (
	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/types"
	"github.com/dad-go/crypto"
	vmtypes "github.com/dad-go/vm/types"
)

//initial a new transaction with asset registration payload
func NewBookKeeperTransaction(pubKey *crypto.PubKey, isAdd bool, cert []byte, issuer *crypto.PubKey) (*types.Transaction, error) {
	bookKeeperPayload := &payload.BookKeeper{
		PubKey: pubKey,
		Action: payload.BookKeeperAction_SUB,
		Cert:   cert,
		Issuer: issuer,
	}

	if isAdd {
		bookKeeperPayload.Action = payload.BookKeeperAction_ADD
	}

	return &types.Transaction{
		TxType:     types.BookKeeper,
		Payload:    bookKeeperPayload,
		Attributes: nil,
	}, nil
}

func NewDeployTransaction(code []byte, name, version, author, email, desp string, vmType vmtypes.VmType, needStorage bool) *types.Transaction {
	//TODO: check arguments
	DeployCodePayload := &payload.DeployCode{
		VmType: vmType,
		Code:        code,
		NeedStorage: needStorage,
		Name:        name,
		Version:     version,
		Author:      author,
		Email:       email,
		Description: desp,
	}

	return &types.Transaction{
		TxType:     types.Deploy,
		Payload:    DeployCodePayload,
		Attributes: nil,
	}
}

func NewInvokeTransaction(vmcode vmtypes.VmCode) *types.Transaction {
	//TODO: check arguments
	invokeCodePayload := &payload.InvokeCode{
		Code:   vmcode,
	}

	return &types.Transaction{
		TxType:     types.Invoke,
		Payload:    invokeCodePayload,
		Attributes: nil,
	}
}
