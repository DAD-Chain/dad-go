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
	"github.com/dad-go/common/serialization"
	. "github.com/dad-go/common"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
)

type VmType byte

const (
	Native = VmType(0xFF)
	NEOVM    = VmType(0x80)
	WASMVM     = VmType(0x90)
	// EVM = VmType(0x90)
)

type VmCode struct {
	VmType VmType
	Code     []byte
}

func (self *VmCode) Serialize(w io.Writer) error {
	w.Write([]byte{byte(self.VmType)})
	return serialization.WriteVarBytes(w, self.Code)

}

func (self *VmCode) Deserialize(r io.Reader) error {
	var b [1]byte
	r.Read(b[:])
	buf, err := serialization.ReadVarBytes(r)
	if err != nil {
		return err
	}
	self.VmType = VmType(b[0])
	self.Code = buf
	return nil
}

func (self *VmCode) AddressFromVmCode() Address {
	var addr Address
	temp := sha256.Sum256(self.Code)
	md := ripemd160.New()
	md.Write(temp[:])
	md.Sum(addr[:0])

	addr[0] = byte(self.VmType)
	return addr
}
