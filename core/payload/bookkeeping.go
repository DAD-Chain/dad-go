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
	"io"

	"github.com/ontio/dad-go/common/serialization"
)

// Bookkeeping is an implementation of transaction payload for bookkeeper rewards
type Bookkeeping struct {
	Nonce uint64
}

func (a *Bookkeeping) Serialize(w io.Writer) error {
	err := serialization.WriteUint64(w, a.Nonce)
	return err
}

func (a *Bookkeeping) Deserialize(r io.Reader) error {
	var err error
	a.Nonce, err = serialization.ReadUint64(r)
	return err
}
