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
package storage

import (
	"github.com/ontio/dad-go/core/store/common"
	"github.com/ontio/dad-go/core/store/leveldbstore"
	"github.com/ontio/dad-go/core/store/overlaydb"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func genRandKeyVal() (string, string) {
	p := make([]byte, 100)
	rand.Read(p)
	key := string(p)
	rand.Read(p)
	val := string(p)
	return key, val
}

func TestCacheDB(t *testing.T) {
	N := 10000
	mem := make(map[string]string)
	memback, _ := leveldbstore.NewMemLevelDBStore()
	overlay := overlaydb.NewOverlayDB(memback)

	cache := NewCacheDB(overlay)
	for i := 0; i < N; i++ {
		key, val := genRandKeyVal()
		cache.Put([]byte(key), []byte(val))
		mem[key] = val
	}

	for key := range mem {
		op := rand.Int() % 2
		if op == 0 {
			//delete
			delete(mem, key)
			cache.Delete([]byte(key))
		} else if op == 1 {
			//update
			_, val := genRandKeyVal()
			mem[key] = val
			cache.Put([]byte(key), []byte(val))
		}
	}

	for key, val := range mem {
		value, err := cache.Get([]byte(key))
		assert.Nil(t, err)
		assert.NotNil(t, value)
		assert.Equal(t, []byte(val), value)
	}
	cache.Commit()

	prefix := common.ST_STORAGE
	for key, val := range mem {
		pkey := make([]byte, 1+len(key))
		pkey[0] = byte(prefix)
		copy(pkey[1:], key)
		raw, err := overlay.Get(pkey)
		assert.Nil(t, err)
		assert.NotNil(t, raw)
		assert.Equal(t, []byte(val), raw)
	}

}
