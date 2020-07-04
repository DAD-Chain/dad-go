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

package neovm

func opCat(e *ExecutionEngine) (VMState, error) {
	b2 := PopByteArray(e)
	b1 := PopByteArray(e)
	r := Concat(b1, b2)
	PushData(e, r)
	return NONE, nil
}

func opSubStr(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	index := PopInt(e)
	arr := PopByteArray(e)
	b := arr[index : index+count]
	PushData(e, b)
	return NONE, nil
}

func opLeft(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	s := PopByteArray(e)
	b := s[:count]
	PushData(e, b)
	return NONE, nil
}

func opRight(e *ExecutionEngine) (VMState, error) {
	count := PopInt(e)
	arr := PopByteArray(e)
	b := arr[len(arr)-count:]
	PushData(e, b)
	return NONE, nil
}

func opSize(e *ExecutionEngine) (VMState, error) {
	x := PeekStackItem(e)
	PushData(e, len(x.GetByteArray()))
	return NONE, nil
}
