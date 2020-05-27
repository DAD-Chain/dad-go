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
	"errors"
	"io"

	. "github.com/dad-go/common"
	"github.com/dad-go/common/serialization"
	. "github.com/dad-go/errors"
	"github.com/ontio/dad-go-crypto/keypair"
)

type Header struct {
	Version          uint32
	PrevBlockHash    Uint256
	TransactionsRoot Uint256
	BlockRoot        Uint256
	Timestamp        uint32
	Height           uint32
	ConsensusData    uint64
	NextBookkeeper   Address

	//Program *program.Program
	Bookkeepers []keypair.PublicKey
	SigData     [][]byte

	hash Uint256
}

//Serialize the blockheader
func (bd *Header) Serialize(w io.Writer) error {
	bd.SerializeUnsigned(w)

	err := serialization.WriteVarUint(w, uint64(len(bd.Bookkeepers)))
	if err != nil {
		return errors.New("serialize sig pubkey length failed")
	}
	for _, pubkey := range bd.Bookkeepers {
		err := serialization.WriteVarBytes(w, keypair.SerializePublicKey(pubkey))
		if err != nil {
			return err
		}
	}

	err = serialization.WriteVarUint(w, uint64(len(bd.SigData)))
	if err != nil {
		return errors.New("serialize sig pubkey length failed")
	}

	for _, sig := range bd.SigData {
		err = serialization.WriteVarBytes(w, sig)
		if err != nil {
			return err
		}
	}

	return nil
}

//Serialize the blockheader data without program
func (bd *Header) SerializeUnsigned(w io.Writer) error {
	serialization.WriteUint32(w, bd.Version)
	bd.PrevBlockHash.Serialize(w)
	bd.TransactionsRoot.Serialize(w)
	bd.BlockRoot.Serialize(w)
	serialization.WriteUint32(w, bd.Timestamp)
	serialization.WriteUint32(w, bd.Height)
	serialization.WriteUint64(w, bd.ConsensusData)
	bd.NextBookkeeper.Serialize(w)
	return nil
}

func (bd *Header) Deserialize(r io.Reader) error {
	err := bd.DeserializeUnsigned(r)
	if err != nil {
		return err
	}

	n, err := serialization.ReadVarUint(r, 0)
	if err != nil {
		return err
	}

	bd.Bookkeepers = make([]keypair.PublicKey, n)
	for i := 0; i < int(n); i++ {
		buf, err := serialization.ReadVarBytes(r)
		if err != nil {
			return err
		}
		bd.Bookkeepers[i], err = keypair.DeserializePublicKey(buf)
		if err != nil {
			return err
		}
	}

	m, err := serialization.ReadVarUint(r, 0)
	if err != nil {
		return err
	}

	bd.SigData = make([][]byte, m)
	for i := 0; i < int(m); i++ {
		sig, err := serialization.ReadVarBytes(r)
		if err != nil {
			return err
		}
		bd.SigData[i] = sig
	}

	return nil
}

func (bd *Header) DeserializeUnsigned(r io.Reader) error {
	//Version
	temp, err := serialization.ReadUint32(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Header item Version Deserialize failed.")
	}
	bd.Version = temp

	//PrevBlockHash
	err = bd.PrevBlockHash.Deserialize(r)
	if err != nil {
		return NewDetailErr(err, ErrNoCode, "Header item preBlock Deserialize failed.")
	}

	//TransactionsRoot
	err = bd.TransactionsRoot.Deserialize(r)
	if err != nil {
		return err
	}

	err = bd.BlockRoot.Deserialize(r)
	if err != nil {
		return err
	}

	//Timestamp
	temp, _ = serialization.ReadUint32(r)
	bd.Timestamp = uint32(temp)

	//Height
	temp, _ = serialization.ReadUint32(r)
	bd.Height = uint32(temp)

	//consensusData
	bd.ConsensusData, _ = serialization.ReadUint64(r)

	//NextBookkeeper
	err = bd.NextBookkeeper.Deserialize(r)

	return err
}

func (bd *Header) Hash() Uint256 {
	buf := new(bytes.Buffer)
	bd.SerializeUnsigned(buf)
	temp := sha256.Sum256(buf.Bytes())
	hash := sha256.Sum256(temp[:])
	return hash
}

func (bd *Header) GetMessage() []byte {
	bf := new(bytes.Buffer)
	bd.SerializeUnsigned(bf)
	return bf.Bytes()
}

func (bd *Header) ToArray() []byte {
	bf := new(bytes.Buffer)
	bd.Serialize(bf)
	return bf.Bytes()
}
