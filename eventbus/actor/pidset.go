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

package actor

const pidSetSliceLen = 16

type PIDSet struct {
	s []string
	m map[string]struct{}
}

// NewPIDSet returns a new PIDSet with the given pids.
func NewPIDSet(pids ...*PID) *PIDSet {
	var s PIDSet
	for _, pid := range pids {
		s.Add(pid)
	}
	return &s
}

func (p *PIDSet) indexOf(v *PID) int {
	for i, pid := range p.s {
		if v.key() == pid {
			return i
		}
	}
	return -1
}

func (p *PIDSet) migrate() {
	p.m = make(map[string]struct{}, pidSetSliceLen)
	for _, v := range p.s {
		p.m[v] = struct{}{}
	}
	p.s = p.s[:0]
}

// Add adds the element v to the set
func (p *PIDSet) Add(v *PID) {
	if p.m == nil {
		if p.indexOf(v) > -1 {
			return
		}

		if len(p.s) < pidSetSliceLen {
			if p.s == nil {
				p.s = make([]string, 0, pidSetSliceLen)
			}
			p.s = append(p.s, v.key())
			return
		}
		p.migrate()
	}
	p.m[v.key()] = struct{}{}
}

// Remove removes v from the set and returns true if them element existed
func (p *PIDSet) Remove(v *PID) bool {
	if p.m == nil {
		i := p.indexOf(v)
		if i == -1 {
			return false
		}
		l := len(p.s) - 1
		p.s[i] = p.s[l]
		p.s = p.s[:l]
		return true
	}
	_, ok := p.m[v.key()]
	if !ok {
		return false
	}
	delete(p.m, v.key())
	return true
}

// Contains reports whether v is an element of the set
func (p *PIDSet) Contains(v *PID) bool {
	if p.m == nil {
		return p.indexOf(v) != -1
	}
	_, ok := p.m[v.key()]
	return ok
}

// Len returns the number of elements in the set
func (p *PIDSet) Len() int {
	if p.m == nil {
		return len(p.s)
	}
	return len(p.m)
}

// Clear removes all the elements in the set
func (p *PIDSet) Clear() {
	if p.m == nil {
		p.s = p.s[:0]
	} else {
		p.m = nil
	}
}

// Empty reports whether the set is empty
func (p *PIDSet) Empty() bool {
	return p.Len() == 0
}

// Values returns all the elements of the set as a slice
func (p *PIDSet) Values() []PID {
	if p.Len() == 0 {
		return nil
	}

	r := make([]PID, p.Len())
	if p.m == nil {
		for i, v := range p.s {
			pidFromKey(v, &r[i])
		}
	} else {
		i := 0
		for v := range p.m {
			pidFromKey(v, &r[i])
			i++
		}
	}
	return r
}

// ForEach invokes f for every element of the set
func (p *PIDSet) ForEach(f func(i int, pid PID)) {
	var pid PID
	if p.m == nil {
		for i, v := range p.s {
			pidFromKey(v, &pid)
			f(i, pid)
		}
	} else {
		i := 0
		for v := range p.m {
			pidFromKey(v, &pid)
			f(i, pid)
			i++
		}
	}
}
