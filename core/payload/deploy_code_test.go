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
	"testing"

	"bytes"
	"github.com/ontio/dad-go/smartcontract/types"
	"github.com/stretchr/testify/assert"
)

func TestDeployCode_Serialize(t *testing.T) {
	deploy := DeployCode{
		Code: types.VmCode{
			VmType: types.NEOVM,
			Code:   []byte{1, 2, 3},
		},
	}

	buf := bytes.NewBuffer(nil)
	deploy.Serialize(buf)
	bs := buf.Bytes()
	var deploy2 DeployCode
	deploy2.Deserialize(buf)
	assert.Equal(t, deploy2, deploy)

	buf = bytes.NewBuffer(bs[:len(bs)-1])
	err := deploy2.Deserialize(buf)
	assert.NotNil(t, err)
}
