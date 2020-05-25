/*
 * Copyright (C) 2018 The dad-go Authors
 * This file is part of The dad-go library.
 *
 * The dad-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The dad-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */

package utils

import (
	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/types"
	"github.com/dad-go/crypto"
	vmtypes "github.com/dad-go/vm/types"
)

//initial a new transaction with asset registration payload
func NewBookkeeperTransaction(pubKey *crypto.PubKey, isAdd bool, cert []byte, issuer *crypto.PubKey) (*types.Transaction, error) {
	bookkeeperPayload := &payload.Bookkeeper{
		PubKey: pubKey,
		Action: payload.BookkeeperAction_SUB,
		Cert:   cert,
		Issuer: issuer,
	}

	if isAdd {
		bookkeeperPayload.Action = payload.BookkeeperAction_ADD
	}

	return &types.Transaction{
		TxType:     types.Bookkeeper,
		Payload:    bookkeeperPayload,
		Attributes: nil,
	}, nil
}

func NewDeployTransaction(code *vmtypes.VmCode, name, version, author, email, desp string, needStorage bool) *types.Transaction {
	//TODO: check arguments
	DeployCodePayload := &payload.DeployCode{
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
