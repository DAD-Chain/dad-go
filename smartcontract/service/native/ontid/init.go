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
package ontid

import (
	"github.com/ontio/dad-go/core/genesis"
	"github.com/ontio/dad-go/smartcontract/service/native"
)

func init() {
	native.Contracts[genesis.OntIDContractAddress] = RegisterIDContract
}

func RegisterIDContract(srvc *native.NativeService) {
	srvc.Register("regIDWithPublicKey", regIdWithPublicKey)
	srvc.Register("addKey", addKey)
	srvc.Register("removeKey", removeKey)
	srvc.Register("addRecovery", addRecovery)
	srvc.Register("changeRecovery", changeRecovery)
	srvc.Register("regIDWithAttributes", regIdWithAttributes)
	srvc.Register("addAttribute", addAttribute)
	srvc.Register("removeAttribute", removeAttribute)
	srvc.Register("verifySignature", verifySignature)
	srvc.Register("getPublicKeys", GetPublicKeys)
	srvc.Register("getKeyState", GetKeyState)
	srvc.Register("getAttributes", GetAttributes)
	srvc.Register("getDDO", GetDDO)
	return
}
