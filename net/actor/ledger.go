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

import (
	"github.com/dad-go/common"
	"github.com/dad-go/common/log"
	ledger "github.com/dad-go/core/ledger/actor"
	"github.com/dad-go/core/types"
	"github.com/dad-go/errors"
	"github.com/ontio/dad-go-eventbus/actor"
	"time"
)

var defLedgerPid *actor.PID

func SetLedgerPid(ledgePid *actor.PID) {
	defLedgerPid = ledgePid
}

func AddHeader(header *types.Header) {
	defLedgerPid.Tell(&ledger.AddHeaderReq{Header: header})
}

func AddHeaders(headers []*types.Header) {
	defLedgerPid.Tell(&ledger.AddHeadersReq{Headers: headers})
}

func AddBlock(block *types.Block) {
	defLedgerPid.Tell(&ledger.AddBlockReq{Block: block})
}

func GetTxnFromLedger(hash common.Uint256) (*types.Transaction, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetTransactionReq{TxHash: hash}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return nil, err
	}
	return result.(*ledger.GetTransactionRsp).Tx, result.(*ledger.GetTransactionRsp).Error
}

func GetCurrentBlockHash() (common.Uint256, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetCurrentBlockHashReq{}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return common.Uint256{}, err
	}
	return result.(*ledger.GetCurrentBlockHashRsp).BlockHash, result.(*ledger.GetCurrentBlockHashRsp).Error
}

func GetCurrentHeaderHash() (common.Uint256, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetCurrentHeaderHashReq{}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return common.Uint256{}, err
	}
	return result.(*ledger.GetCurrentHeaderHashRsp).BlockHash, result.(*ledger.GetCurrentHeaderHashRsp).Error
}

func GetBlockHashByHeight(height uint32) (common.Uint256, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetBlockHashReq{Height: height}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return common.Uint256{}, err
	}
	return result.(*ledger.GetBlockHashRsp).BlockHash, result.(*ledger.GetBlockHashRsp).Error
}

func GetHeaderByHeight(height uint32) (*types.Header, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetHeaderByHeightReq{Height: height}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return nil, err
	}
	return result.(*ledger.GetHeaderByHeightRsp).Header, result.(*ledger.GetHeaderByHeightRsp).Error
}

func GetBlockByHeight(height uint32) (*types.Block, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetBlockByHeightReq{Height: height}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return nil, err
	}
	return result.(*ledger.GetBlockByHeightRsp).Block, result.(*ledger.GetBlockByHeightRsp).Error
}

func GetHeaderByHash(hash common.Uint256) (*types.Header, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetHeaderByHashReq{BlockHash: hash}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return nil, err
	}
	return result.(*ledger.GetHeaderByHashRsp).Header, result.(*ledger.GetHeaderByHashRsp).Error
}

func GetBlockByHash(hash common.Uint256) (*types.Block, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetBlockByHashReq{BlockHash: hash}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return nil, err
	}
	return result.(*ledger.GetBlockByHashRsp).Block, result.(*ledger.GetBlockByHashRsp).Error
}

func GetCurrentHeaderHeight() (uint32, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetCurrentHeaderHeightReq{}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return 0, err
	}
	return result.(*ledger.GetCurrentHeaderHeightRsp).Height, result.(*ledger.GetCurrentHeaderHeightRsp).Error
}

func GetCurrentBlockHeight() (uint32, error) {
	future := defLedgerPid.RequestFuture(&ledger.GetCurrentBlockHeightReq{}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return 0, err
	}
	return result.(*ledger.GetCurrentBlockHeightRsp).Height, result.(*ledger.GetCurrentBlockHeightRsp).Error
}

func IsContainBlock(hash common.Uint256) (bool, error) {
	future := defLedgerPid.RequestFuture(&ledger.IsContainBlockReq{BlockHash: hash}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
		return false, err
	}
	return result.(*ledger.IsContainBlockRsp).IsContain, result.(*ledger.IsContainBlockRsp).Error
}
