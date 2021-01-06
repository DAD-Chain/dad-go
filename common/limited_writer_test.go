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
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLimitedWriter_Write(t *testing.T) {
	bf := bytes.NewBuffer(nil)
	writer := NewLimitedWriter(bf, 5)
	_, err := writer.Write([]byte{1, 2, 3})
	assert.Nil(t, err)
	assert.Equal(t, bf.Bytes(), []byte{1, 2, 3})
	_, err = writer.Write([]byte{4, 5})
	assert.Nil(t, err)

	_, err = writer.Write([]byte{6})
	assert.Equal(t, err, ErrWriteExceedLimitedCount)
}
