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
	"fmt"
	"io"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/serialization"
)

// InvokeCode is an implementation of transaction payload for invoke smartcontract
type InvokeCode struct {
	Code []byte
}

func (self *InvokeCode) Serialize(w io.Writer) error {
	if err := serialization.WriteVarBytes(w, self.Code); err != nil {
		return fmt.Errorf("InvokeCode Code Serialize failed: %s", err)
	}
	return nil
}

func (self *InvokeCode) Deserialize(r io.Reader) error {
	code, err := serialization.ReadVarBytes(r)
	if err != nil {
		return fmt.Errorf("InvokeCode Code Deserialize failed: %s", err)
	}
	self.Code = code
	return nil
}

//note: InvokeCode.Code has data reference of param source
func (self *InvokeCode) Deserialization(source *common.ZeroCopySource) error {
	code, _, irregular, eof := source.NextVarBytes()
	if eof {
		return io.ErrUnexpectedEOF
	}
	if irregular {
		return common.ErrIrregularData
	}

	self.Code = code
	return nil
}

func (self *InvokeCode) Serialization(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(self.Code)
}
