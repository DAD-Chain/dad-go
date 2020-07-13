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

package types

import (
	"github.com/ontio/dad-go-eventbus/actor"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/errors"
)

// message
type RegisterValidator struct {
	Sender *actor.PID
	Type   VerifyType
	Id     string
}

type UnRegisterValidator struct {
	Id   string
	Type VerifyType
}

type UnRegisterAck struct {
	Id   string
	Type VerifyType
}

type CheckTx struct {
	WorkerId uint8
	Tx       types.Transaction
}

type CheckResponse struct {
	WorkerId uint8
	Type     VerifyType
	Hash     common.Uint256
	Height   uint32
	ErrCode  errors.ErrCode
}

// VerifyType of validator
type VerifyType uint8

const (
	Stateless VerifyType = iota
	Statefull VerifyType = iota
)
