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
package states

import (
	"testing"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/common"
	"github.com/stretchr/testify/assert"
)

func TestValidatorState_Deserialize_Serialize(t *testing.T) {
	_, pubKey, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)

	vs := ValidatorState{
		StateBase: StateBase{(byte)(1)},
		PublicKey: pubKey,
	}

	sink := common.NewZeroCopySink(nil)
	vs.Serialization(sink)
	bs := sink.Bytes()

	var vs2 ValidatorState
	source := common.NewZeroCopySource(sink.Bytes())
	vs2.Deserialization(source)
	assert.Equal(t, vs, vs2)

	source = common.NewZeroCopySource(bs[:len(bs)-1])
	err := vs2.Deserialization(source)
	assert.NotNil(t, err)
}
