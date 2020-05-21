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
	. "github.com/dad-go/common"
	"github.com/dad-go/common/serialization"
	"github.com/dad-go/crypto"
	"io"
	. "github.com/dad-go/errors"
)

const (
	MaxVoteKeys = 1024
)

type Vote struct {
	PubKeys []*crypto.PubKey // vote node list

	Account Address
}

func (self *Vote) Check() bool {
	if len(self.PubKeys) > MaxVoteKeys {
		return false
	}
	return true
}

func (self *Vote) Serialize(w io.Writer) error {
	if err := serialization.WriteUint32(w, uint32(len(self.PubKeys))); err != nil {
		return NewDetailErr(err, ErrNoCode, "Vote PubKeys length Serialize failed.")
	}
	for _, key := range self.PubKeys {
		if err := key.Serialize(w); err != nil {
			return NewDetailErr(err, ErrNoCode, "InvokeCode PubKeys Serialize failed.")
		}
	}
	if err := self.Account.Serialize(w); err != nil {
		return NewDetailErr(err, ErrNoCode, "InvokeCode Account Serialize failed.")
	}

	return nil
}

func (self *Vote) Deserialize(r io.Reader) error {
	length, err := serialization.ReadUint32(r)
	if err != nil {
		return err
	}
	self.PubKeys = make([]*crypto.PubKey, length)
	for i := 0; i < int(length); i++ {
		pubkey := new(crypto.PubKey)
		err := pubkey.DeSerialize(r)
		if err != nil {
			return err
		}
		self.PubKeys[i] = pubkey
	}

	err = self.Account.Deserialize(r)
	if err != nil {
		return err
	}

	return nil
}
