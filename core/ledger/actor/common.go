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
	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/states"
	"github.com/dad-go/core/types"
)

type AddHeaderReq struct {
	Header *types.Header
}

type AddHeaderRsp struct {
	BlockHash common.Uint256
	Error     error
}

type AddHeadersReq struct {
	Headers []*types.Header
}

type AddHeadersRsp struct {
	BlockHashes []common.Uint256
	Error       error
}

type AddBlockReq struct {
	Block *types.Block
}

type AddBlockRsp struct {
	BlockHash common.Uint256
	Error     error
}

type GetTransactionReq struct {
	TxHash common.Uint256
}

type GetTransactionRsp struct {
	Tx    *types.Transaction
	Error error
}

type GetBlockByHashReq struct {
	BlockHash common.Uint256
}

type GetBlockByHashRsp struct {
	Block *types.Block
	Error error
}

type GetBlockByHeightReq struct {
	Height uint32
}

type GetBlockByHeightRsp struct {
	Block *types.Block
	Error error
}

type GetHeaderByHashReq struct {
	BlockHash common.Uint256
}

type GetHeaderByHashRsp struct {
	Header *types.Header
	Error  error
}

type GetHeaderByHeightReq struct {
	Height uint32
}

type GetHeaderByHeightRsp struct {
	Header *types.Header
	Error  error
}

type GetCurrentBlockHashReq struct{}

type GetCurrentBlockHashRsp struct {
	BlockHash common.Uint256
	Error     error
}

type GetCurrentBlockHeightReq struct{}

type GetCurrentBlockHeightRsp struct {
	Height uint32
	Error  error
}

type GetCurrentHeaderHeightReq struct{}

type GetCurrentHeaderHeightRsp struct {
	Height uint32
	Error  error
}

type GetCurrentHeaderHashReq struct{}

type GetCurrentHeaderHashRsp struct {
	BlockHash common.Uint256
	Error     error
}

type GetBlockHashReq struct {
	Height uint32
}

type GetBlockHashRsp struct {
	BlockHash common.Uint256
	Error     error
}

type IsContainBlockReq struct {
	BlockHash common.Uint256
}

type IsContainBlockRsp struct {
	IsContain bool
	Error     error
}

type GetBlockRootWithNewTxRootReq struct {
	TxRoot common.Uint256
}

type GetBlockRootWithNewTxRootRsp struct {
	NewTxRoot common.Uint256
	Error     error
}

type GetTransactionWithHeightReq struct {
	TxHash common.Uint256
}

type GetTransactionWithHeightRsp struct {
	Tx     *types.Transaction
	Height uint32
	Error  error
}

type IsContainTransactionReq struct {
	TxHash common.Uint256
}

type IsContainTransactionRsp struct {
	IsContain bool
	Error     error
}

type GetCurrentStateRootReq struct{}

type GetCurrentStateRootRsp struct {
	StateRoot common.Uint256
	Error     error
}

type GetBookKeeperStateReq struct{}

type GetBookKeeperStateRsp struct {
	BookKeepState *states.BookKeeperState
	Error         error
}

type GetStorageItemReq struct {
	CodeHash *common.Address
	Key      []byte
}

type GetStorageItemRsp struct {
	Value []byte
	Error error
}

type GetContractStateReq struct {
	ContractHash common.Address
}

type GetContractStateRsp struct {
	ContractState *payload.DeployCode
	Error         error
}

type PreExecuteContractReq struct {
	Tx *types.Transaction
}

type PreExecuteContractRsp struct {
	Result []interface{}
	Error  error
}
