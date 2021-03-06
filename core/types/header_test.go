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
	"testing"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/stretchr/testify/assert"
)

func TestHeader_Serialize(t *testing.T) {
	header := Header{}
	header.Height = 321
	header.Bookkeepers = make([]keypair.PublicKey, 0)
	header.SigData = make([][]byte, 0)
	buf := bytes.NewBuffer(nil)
	err := header.Serialize(buf)
	bs := buf.Bytes()
	assert.Nil(t, err)

	var h2 Header
	h2.Deserialize(buf)
	assert.Equal(t, fmt.Sprint(header), fmt.Sprint(h2))

	buf = bytes.NewBuffer(bs[:len(bs)-1])
	err = h2.Deserialize(buf)

	assert.NotNil(t, err)
}
