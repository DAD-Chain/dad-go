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

package ledgerstore

import (
	"fmt"
	"github.com/dad-go/common"
	"github.com/dad-go/core/types"
	"github.com/hashicorp/golang-lru"
)

const (
	BlockCacheSize       = 1000
	TransactionCacheSize = 100000
)

type TransactionCacheaValue struct {
	Tx *types.Transaction
	Height uint32
}

type BlockCache struct {
	blockCache       *lru.ARCCache
	transactionCache *lru.ARCCache
}

func NewBlockCache() (*BlockCache, error) {
	blockCache, err := lru.NewARC(BlockCacheSize)
	if err != nil {
		return nil, fmt.Errorf("NewARC block error %s", err)
	}
	transactionCache, err := lru.NewARC(TransactionCacheSize)
	if err != nil {
		return nil, fmt.Errorf("NewARC header error %s", err)
	}
	return &BlockCache{
		blockCache:       blockCache,
		transactionCache: transactionCache,
	}, nil
}

func (this *BlockCache) AddBlock(block *types.Block) {
	blockHash := block.Hash()
	this.blockCache.Add(string(blockHash.ToArray()), block)
}

func (this *BlockCache) GetBlock(blockHash common.Uint256) *types.Block {
	block, ok := this.blockCache.Get(string(blockHash.ToArray()))
	if !ok {
		return nil
	}
	return block.(*types.Block)
}

func (this *BlockCache) ContainBlock(blockHash common.Uint256) bool{
	return this.blockCache.Contains(string(blockHash.ToArray()))
}

func (this *BlockCache) AddTransaction(tx *types.Transaction, height uint32) {
	txHash := tx.Hash()
	this.transactionCache.Add(string(txHash.ToArray()), &TransactionCacheaValue{
		Tx:tx,
		Height:height,
	})
}

func (this *BlockCache) GetTransaction(txHash common.Uint256) (*types.Transaction ,uint32){
	value, ok := this.transactionCache.Get(string(txHash.ToArray()))
	if !ok {
		return nil, 0
	}
	txValue := value.(*TransactionCacheaValue)
	return txValue.Tx, txValue.Height
}

func (this *BlockCache) ContainTransaction(txHash common.Uint256) bool{
	return this.transactionCache.Contains(string(txHash.ToArray()))
}
