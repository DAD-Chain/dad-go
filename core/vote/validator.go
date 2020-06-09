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

package vote

import (
	"sort"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/core/genesis"
	"github.com/ontio/dad-go/core/states"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go-crypto/keypair"
)

func GetValidators(txs []*types.Transaction) ([]keypair.PublicKey, error) {
	// TODO implement vote
	return genesis.GenesisBookkeepers, nil
}

func weightedAverage(votes []*states.VoteState) int64 {
	var sumWeight, sumValue int64
	for _, v := range votes {
		sumWeight += v.Count.GetData()
		sumValue += v.Count.GetData() * int64(len(v.PublicKeys))
	}
	if sumValue == 0 {
		return 0
	}
	return sumValue / sumWeight
}

type Pair struct {
	Key   string
	Value int64
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int      { return len(p) }
func (p PairList) Less(i, j int) bool {
	if p[j].Value < p[i].Value {
		return true
	} else if p[j].Value > p[i].Value {
		return false
	} else {
		return p[j].Key < p[i].Key
	}
}

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]common.Fixed64) []string {
	p := make(PairList, 0, len(m))
	for k, v := range m {
		p = append(p, Pair{k, v.GetData()})
	}
	sort.Sort(p)
	keys := make([]string, len(m))
	for i, k := range p {
		keys[i] = k.Key
	}
	return keys
}
