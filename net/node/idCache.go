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

package node

import (
	"github.com/dad-go/common"
	"sync"
)

type idCache struct {
	sync.RWMutex
	list map[common.Uint256]bool
}

func (c *idCache) init() {
}

func (c *idCache) add() {
}

func (c *idCache) del() {
}

func (c *idCache) ExistedID(id common.Uint256) bool {
	// TODO
	return false
}
