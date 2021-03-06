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

package feeSplit

import (
	"bytes"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/genesis"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/smartcontract/service/native"
	"github.com/ontio/dad-go/smartcontract/service/native/ont"
)

func AppCallApproveOng(native *native.NativeService, from common.Address, to common.Address, amount uint64) error {
	buf := bytes.NewBuffer(nil)
	sts := &ont.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	err := sts.Serialize(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[appCallApproveOng] transfers.Serialize error!")
	}

	if _, err := native.ContextRef.AppCall(genesis.OngContractAddress, "approve", []byte{}, buf.Bytes()); err != nil {
		return errors.NewDetailErr(err, errors.ErrNoCode, "[appCallApproveOng] appCall error!")
	}
	return nil
}

func splitCurve(pos uint64, avg uint64) uint64 {
	xi := PRECISE * YITA * 2 * pos / (avg * 10)
	index := xi / (PRECISE / 10)
	s := ((Yi[index+1]-Yi[index])*xi + Yi[index]*Xi[index+1] - Yi[index+1]*Xi[index]) / (Xi[index+1] - Xi[index])
	return s
}
