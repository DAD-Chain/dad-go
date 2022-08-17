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

package states

import (
	"github.com/ontio/dad-go/common"
	"io"
)

type StorageKey struct {
	ContractAddress common.Address
	Key             []byte
}

func (this *StorageKey) Serialization(sink *common.ZeroCopySink) {
	this.ContractAddress.Serialization(sink)
	sink.WriteVarBytes(this.Key)
}

func (this *StorageKey) Deserialization(source *common.ZeroCopySource) error {
	if err := this.ContractAddress.Deserialization(source); err != nil {
		return err
	}
	key, _, irregular, eof := source.NextVarBytes()
	if irregular {
		return common.ErrIrregularData
	}
	if eof {
		return io.ErrUnexpectedEOF
	}
	this.Key = key
	return nil
}

func (this *StorageKey) ToArray() []byte {
	return common.SerializeToBytes(this)
}
