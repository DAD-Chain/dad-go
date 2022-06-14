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
package wasmvm

import (
	"errors"

	"github.com/go-interpreter/wagon/exec"
)

func ReadWasmMemory(proc *exec.Process, ptr uint32, len uint32) ([]byte, error) {
	if uint64(proc.MemSize()) < uint64(ptr)+uint64(len) {
		return nil, errors.New("contract create len is greater than memory size")
	}
	keybytes := make([]byte, len)
	_, err := proc.ReadAt(keybytes, int64(ptr))
	if err != nil {
		return nil, err
	}

	return keybytes, nil
}
