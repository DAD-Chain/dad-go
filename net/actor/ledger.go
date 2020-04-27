package actor

import (
	"github.com/ONTID/eventbus/actor"
	"github.com/dad-go/core/types"
	"github.com/dad-go/common"
	"github.com/dad-go/common/log"
	"github.com/dad-go/errors"
	"time"
)

var DefLedgerPid actor.PID

type AddHeaderReq struct{
	Header *[]types.Header
}
type AddHeaderRsp struct{
	BlockHash *common.Uint256
	Error error
}

type AddBlockReq struct{
	Block *types.Block
}
type AddBlockRsp struct{
	BlockHash *common.Uint256
	Error error
}

type GetTransactionReq struct{
	TxHash common.Uint256
}
type GetTransactionRsp struct{
	Tx *types.Transaction
	Error error
}

type GetBlockByHashReq struct{
	BlockHash common.Uint256
}
type GetBlockByHashRsp struct{
	Block *types.Block
	Error error
}

type GetBlockByHeightReq struct{
	Height uint32
}
type GetBlockByHeightRsp struct{
	Block *types.Block
	Error error
}

type GetHeaderByHashReq struct{
	BlockHash common.Uint256
}
type GetHeaderByHashRsp struct{
	Header *types.Header
	Error error
}

type GetHeaderByHeightReq struct{
	Height uint32
}
type GetHeaderByHeightRsp struct{
	Header *types.Header
	Error error
}

type GetCurrentBlockHashReq struct{}
type GetCurrentBlockHashRsp struct{
	BlockHash common.Uint256
	Error error
}

type GetCurrentBlockHeightReq struct{}
type GetCurrentBlockHeightRsp struct{
	Height uint32
	Error error
}

type GetCurrentHeaderHeightReq struct{}
type GetCurrentHeaderHeightRsp struct{
	Height uint32
	Error error
}

type GetBlockHashReq struct{
	Height uint32
}
type GetBlockHashRsp struct{
	BlockHash common.Uint256
	Error error
}

type IsContainBlockReq struct{
	BlockHash common.Uint256
}
type IsContainBlockRsp struct{
	IsContain bool
	Error error
}

//------------------------------------------------------------------------------------
func AddHeader(header *[]types.Header){
	DefLedgerPid.Tell(&AddHeaderReq{Header: header})
}

func AddBlock(block *types.Block){
	DefLedgerPid.Tell(&AddBlockReq{Block: block})
}

func GetTxnFromLedger(hash common.Uint256)(*types.Transaction, error){
	future := DefLedgerPid.RequestFuture(&GetTransactionReq{hash}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetTransactionRsp).Tx, result.(GetTransactionRsp).Error
}

func GetCurrentBlockHash() (common.Uint256, error) {
	future := DefLedgerPid.RequestFuture(&GetCurrentBlockHashReq{}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetCurrentBlockHashRsp).BlockHash, result.(GetCurrentBlockHashRsp).Error
}

func GetBlockHashByHeight(height uint32) (common.Uint256, error) {
	future := DefLedgerPid.RequestFuture(&GetBlockHashReq{height}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetBlockHashRsp).BlockHash, result.(GetBlockHashRsp).Error
}

func GetHeaderByHeight(height uint32) (*types.Header, error){
	future := DefLedgerPid.RequestFuture(&GetHeaderByHeightReq{height}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetHeaderByHeightRsp).Header, result.(GetHeaderByHeightRsp).Error
}

func GetBlockByHeight(height uint32)(*types.Block, error){
	future := DefLedgerPid.RequestFuture(&GetBlockByHeightReq{height}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetBlockByHeightRsp).Block, result.(GetBlockByHeightRsp).Error
}

func GetHeaderByHash(hash common.Uint256)(*types.Header, error){
	future := DefLedgerPid.RequestFuture(&GetHeaderByHashReq{hash}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetHeaderByHashRsp).Header, result.(GetHeaderByHashRsp).Error
}

func GetBlockByHash(hash common.Uint256)(*types.Block, error){
	future := DefLedgerPid.RequestFuture(&GetBlockByHashReq{hash}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetBlockByHashRsp).Block, result.(GetBlockByHashRsp).Error
}

func GetCurrentHeaderHeight() (uint32, error) {
	future := DefLedgerPid.RequestFuture(&GetCurrentHeaderHeightReq{}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetCurrentHeaderHeightRsp).Height, result.(GetCurrentHeaderHeightRsp).Error
}

func GetCurrentBlockHeight() (uint32, error) {
	future := DefLedgerPid.RequestFuture(&GetCurrentBlockHeightReq{}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(GetCurrentBlockHeightRsp).Height, result.(GetCurrentBlockHeightRsp).Error
}

func IsContainBlock(hash common.Uint256) (bool, error) {
	future := DefLedgerPid.RequestFuture(&IsContainBlockReq{hash}, 5*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Error(errors.NewErr("ERROR: "), err)
	}
	return result.(IsContainBlockRsp).IsContain, result.(IsContainBlockRsp).Error
}
