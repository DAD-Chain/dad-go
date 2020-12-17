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

	"github.com/ontio/dad-go/common/serialization"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/p2pserver/common"
)

type Ping struct {
	Height uint64
}

//Serialize message payload
func (this Ping) Serialization() ([]byte, error) {
	p := bytes.NewBuffer([]byte{})
	err := serialization.WriteUint64(p, this.Height)
	if err != nil {
		return nil, errors.NewDetailErr(err, errors.ErrNetPackFail, fmt.Sprintf("write error. Height:%v", this.Height))
	}

	return p.Bytes(), nil
}

func (this *Ping) CmdType() string {
	return common.PING_TYPE
}

//Deserialize message payload
func (this *Ping) Deserialization(p []byte) error {
	buf := bytes.NewBuffer(p)

	height, err := serialization.ReadUint64(buf)
	if err != nil {
		return errors.NewDetailErr(err, errors.ErrNetUnPackFail, fmt.Sprintf("read Height error. buf:%v", buf))
	}
	this.Height = height
	return nil
}
