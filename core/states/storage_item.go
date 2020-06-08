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
	"bytes"
	"io"

	"github.com/dad-go/common/serialization"
)

type StorageItem struct {
	StateBase
	Value []byte
}

func (this *StorageItem) Serialize(w io.Writer) error {
	this.StateBase.Serialize(w)
	serialization.WriteVarBytes(w, this.Value)
	return nil
}

func (this *StorageItem) Deserialize(r io.Reader) error {
	if this == nil {
		this = new(StorageItem)
	}
	err := this.StateBase.Deserialize(r)
	if err != nil {
		return err
	}
	value, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	this.Value = value
	return nil
}

func (storageItem *StorageItem) ToArray() []byte {
	b := new(bytes.Buffer)
	storageItem.Serialize(b)
	return b.Bytes()
}
