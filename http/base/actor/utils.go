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

// Package actor privides communication with other actor
package actor

import (
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/smartcontract/service/native/utils"
)

func updateNativeSCAddr(hash common.Address) common.Address {
	if hash == utils.OntContractAddress {
		hash = common.AddressFromVmCode(utils.OntContractAddress[:])
	} else if hash == utils.OngContractAddress {
		hash = common.AddressFromVmCode(utils.OngContractAddress[:])
	} else if hash == utils.OntIDContractAddress {
		hash = common.AddressFromVmCode(utils.OntIDContractAddress[:])
	} else if hash == utils.ParamContractAddress {
		hash = common.AddressFromVmCode(utils.ParamContractAddress[:])
	} else if hash == utils.AuthContractAddress {
		hash = common.AddressFromVmCode(utils.AuthContractAddress[:])
	} else if hash == utils.GovernanceContractAddress {
		hash = common.AddressFromVmCode(utils.GovernanceContractAddress[:])
	}
	return hash
}
