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
	"crypto/sha256"
	"github.com/dad-go/crypto"
	. "github.com/dad-go/common"
	"golang.org/x/crypto/ripemd160"
	"errors"
	"github.com/dad-go/common/serialization"
	"sort"
)

func AddressFromPubKey(pubkey *crypto.PubKey) Address {
	buf := bytes.Buffer{}
	pubkey.Serialize(&buf)

	var addr Address
	temp := sha256.Sum256(buf.Bytes())
	md := ripemd160.New()
	md.Write(temp[:])
	md.Sum(addr[:0])

	addr[0] = 0x01

	return addr
}

func AddressFromMultiPubKeys(pubkeys []*crypto.PubKey, m int) (Address, error) {
	var addr Address
	n := len(pubkeys)
	if m <= 0 || m > n || n > 24 {
		return addr, errors.New("wrong multi-sig param")
	}
	sort.Sort(crypto.PubKeySlice(pubkeys))
	buf := bytes.Buffer{}
	serialization.WriteUint8(&buf, uint8(n))
	serialization.WriteUint8(&buf, uint8(m))
	for _, pubkey := range pubkeys {
		pubkey.Serialize(&buf)
	}

	temp := sha256.Sum256(buf.Bytes())
	md := ripemd160.New()
	md.Write(temp[:])
	md.Sum(addr[:0])
	addr[0] = 0x02

	return addr, nil
}

func AddressFromBookkeepers(bookkeepers []*crypto.PubKey) (Address, error) {
	return AddressFromMultiPubKeys(bookkeepers, len(bookkeepers)-(len(bookkeepers)-1)/3)
}

