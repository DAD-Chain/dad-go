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

package common

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressFromBase58(t *testing.T) {
	var addr Address
	rand.Read(addr[:])

	base58 := addr.ToBase58()
	b1 := string(append([]byte{'X'}, []byte(base58)...))
	_, err := AddressFromBase58(b1)

	assert.NotNil(t, err)

	b2 := string([]byte(base58)[1:10])
	_, err = AddressFromBase58(b2)

	assert.NotNil(t, err)
}

func TestAddressParseFromBytes(t *testing.T) {
	var addr Address
	rand.Read(addr[:])

	addr2, _ := AddressParseFromBytes(addr[:])

	assert.Equal(t, addr, addr2)
}

func TestAddress_Serialize(t *testing.T) {
	var addr Address
	rand.Read(addr[:])

	sink := NewZeroCopySink(nil)
	addr.Serialization(sink)

	var addr2 Address
	source := NewZeroCopySource(sink.Bytes())
	addr2.Deserialization(source)
	assert.Equal(t, addr, addr2)
}
