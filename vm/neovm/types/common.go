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
	"math/big"

	"github.com/ontio/dad-go/common"
)

func ConvertBigIntegerToBytes(data *big.Int) []byte {
	if data.Int64() == 0 {
		return []byte{}
	}

	bs := data.Bytes()
	b := bs[0]
	if data.Sign() < 0 {
		for i, b := range bs {
			bs[i] = ^b
		}
		temp := big.NewInt(0)
		temp.SetBytes(bs)
		temp2 := big.NewInt(0)
		temp2.Add(temp, big.NewInt(1))
		bs = temp2.Bytes()
		common.BytesReverse(bs)
		if b>>7 == 1 {
			bs = append(bs, 255)
		}
	} else {
		common.BytesReverse(bs)
		if b>>7 == 1 {
			bs = append(bs, 0)
		}
	}
	return bs
}

func ConvertBytesToBigInteger(ba []byte) *big.Int {
	res := big.NewInt(0)
	l := len(ba)
	if l == 0 {
		return res
	}

	bytes := make([]byte, 0, l)
	bytes = append(bytes, ba...)
	common.BytesReverse(bytes)

	if bytes[0]>>7 == 1 {
		for i, b := range bytes {
			bytes[i] = ^b
		}

		temp := big.NewInt(0)
		temp.SetBytes(bytes)
		temp2 := big.NewInt(0)
		temp2.Add(temp, big.NewInt(1))
		bytes = temp2.Bytes()
		res.SetBytes(bytes)
		return res.Neg(res)
	}

	res.SetBytes(bytes)
	return res
}
