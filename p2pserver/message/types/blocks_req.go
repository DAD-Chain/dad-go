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
	"io"

	comm "github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/p2pserver/common"
)

type BlocksReq struct {
	HeaderHashCount uint8
	HashStart       comm.Uint256
	HashStop        comm.Uint256
}

//Serialize message payload
func (this *BlocksReq) Serialization(sink *comm.ZeroCopySink) {
	sink.WriteUint8(this.HeaderHashCount)
	sink.WriteHash(this.HashStart)
	sink.WriteHash(this.HashStop)
}

func (this *BlocksReq) CmdType() string {
	return common.GET_BLOCKS_TYPE
}

//Deserialize message payload
func (this *BlocksReq) Deserialization(source *comm.ZeroCopySource) error {
	var eof bool
	this.HeaderHashCount, eof = source.NextUint8()
	this.HashStart, eof = source.NextHash()
	this.HashStop, eof = source.NextHash()

	if eof {
		return io.ErrUnexpectedEOF
	}
	return nil
}
