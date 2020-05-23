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

package payload

import (
	"github.com/dad-go/common"
	vmtypes "github.com/dad-go/vm/types"
	"io"
	. "github.com/dad-go/errors"
)

type InvokeCode struct {
	GasLimit common.Fixed64
	Code     vmtypes.VmCode
}

func (self *InvokeCode) Serialize(w io.Writer) error {
	var err error
	err = self.GasLimit.Serialize(w)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode GasLimit Serialize failed.")
	}
	err = self.Code.Serialize(w)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode Code Serialize failed.")
	}
	return err
}

func (self *InvokeCode) Deserialize(r io.Reader) error {
	var err error

	err = self.GasLimit.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode GasLimit Deserialize failed.")
	}
	err = self.Code.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode Code Deserialize failed.")
	}
	return nil
}