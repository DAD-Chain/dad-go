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
	states "github.com/dad-go/core/states"
	."github.com/dad-go/common"
	"github.com/dad-go/smartcontract/event"
)

type IIterator interface {
	Next() bool
	Prev() bool
	First() bool
	Last() bool
	Seek(key []byte) bool
	Key() []byte
	Value() []byte
	Release()
}

type IStore interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	NewBatch() error
	BatchPut(key []byte, value []byte) error
	BatchDelete(key []byte) error
	BatchCommit() error
	Close() error
	NewIterator(prefix []byte) IIterator
}

type IStateStore interface {
	TryAdd(prefix DataEntryPrefix, key []byte, value states.IStateValue, trie bool)
	TryGetOrAdd(prefix DataEntryPrefix, key []byte, value states.IStateValue, trie bool) error
	TryGet(prefix DataEntryPrefix, key []byte) (*StateItem, error)
	TryGetAndChange(prefix DataEntryPrefix, key []byte, trie bool) (states.IStateValue, error)
	TryDelete(prefix DataEntryPrefix, key []byte)
	Find(prefix DataEntryPrefix, key []byte) ([]*StateItem, error)
}

type IMemoryStore interface {
	Put(prefix byte, key []byte, value states.IStateValue, state ItemState, trie bool)
	Get(prefix byte, key []byte) *StateItem
	Delete(prefix byte, key []byte)
	GetChangeSet() map[string]*StateItem
	Change(prefix byte, key []byte, trie bool)
}

type IEventStore interface {
	SaveEventNotifyByTx(txHash Uint256, notifies []*event.NotifyEventInfo) error
	SaveEventNotifyByBlock(height uint32, txHashs []Uint256) error
	GetEventNotifyByTx(txHash Uint256) ([]*event.NotifyEventInfo, error)
	CommitTo() error
}

type ItemState byte

const (
	None ItemState = iota
	Changed
	Deleted
)

type StateItem struct {
	Key   string
	Value states.IStateValue
	State ItemState
	Trie  bool
}

func (e *StateItem) copy() *StateItem {
	c := *e; return &c
}
