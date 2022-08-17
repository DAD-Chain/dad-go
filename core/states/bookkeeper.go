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
	"io"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/common"
)

type BookkeeperState struct {
	StateBase
	CurrBookkeeper []keypair.PublicKey
	NextBookkeeper []keypair.PublicKey
}

func (this *BookkeeperState) Serialization(sink *common.ZeroCopySink) {
	this.StateBase.Serialization(sink)
	sink.WriteUint32(uint32(len(this.CurrBookkeeper)))
	for _, v := range this.CurrBookkeeper {
		buf := keypair.SerializePublicKey(v)
		sink.WriteVarBytes(buf)
	}
	sink.WriteUint32(uint32(len(this.NextBookkeeper)))
	for _, v := range this.NextBookkeeper {
		buf := keypair.SerializePublicKey(v)
		sink.WriteVarBytes(buf)
	}
}

func (this *BookkeeperState) Deserialization(source *common.ZeroCopySource) error {
	err := this.StateBase.Deserialization(source)
	if err != nil {
		return err
	}
	n, eof := source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < int(n); i++ {
		buf, _, irregular, eof := source.NextVarBytes()
		if irregular {
			return common.ErrIrregularData
		}
		if eof {
			return io.ErrUnexpectedEOF
		}
		key, err := keypair.DeserializePublicKey(buf)
		if err != nil {
			return err
		}
		this.CurrBookkeeper = append(this.CurrBookkeeper, key)
	}
	n, eof = source.NextUint32()
	if eof {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < int(n); i++ {
		buf, _, irregular, eof := source.NextVarBytes()
		if irregular {
			return common.ErrIrregularData
		}
		if eof {
			return io.ErrUnexpectedEOF
		}
		key, err := keypair.DeserializePublicKey(buf)
		if err != nil {
			return err
		}
		this.NextBookkeeper = append(this.NextBookkeeper, key)
	}
	return nil
}

func (v *BookkeeperState) ToArray() []byte {
	return common.SerializeToBytes(v)
}
