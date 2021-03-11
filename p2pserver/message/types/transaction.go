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
	"bytes"
	"fmt"

	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/p2pserver/common"
)

// Transaction message
type Trn struct {
	Txn *types.Transaction
}

//Serialize message payload
func (this Trn) Serialization() ([]byte, error) {
	p := bytes.NewBuffer(nil)
	err := this.Txn.Serialize(p)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("serialize error. Txn:%v", this.Txn))
	}

	return p.Bytes(), nil
}

func (this *Trn) CmdType() string {
	return common.TX_TYPE
}

//Deserialize message payload
func (this *Trn) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)

	tx, err := types.TransactionFromRawBytes(p)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read txn error. buf:%v", buf))
	}

	this.Txn = tx
	return nil
}
