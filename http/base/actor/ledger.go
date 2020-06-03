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
	"time"
	"github.com/ontio/dad-go-eventbus/actor"
	"github.com/dad-go/core/types"
	lactor "github.com/dad-go/core/ledger/actor"
	"github.com/dad-go/common"
	"errors"
	"github.com/dad-go/smartcontract/event"
	"github.com/dad-go/core/payload"
	"github.com/dad-go/common/log"
)

const (
	REQ_TIMEOUT   = 5
	ERR_ACTOR_COMM = "[http] Actor comm error: %v"
)

var defLedgerPid *actor.PID

func SetLedgerPid(actr *actor.PID) {
	defLedgerPid = actr
}

func GetBlockHashFromStore(height uint32) (common.Uint256, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetBlockHashReq{height}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return common.Uint256{}, err
	}
	if rsp, ok := result.(*lactor.GetBlockHashRsp); !ok {
		return common.Uint256{}, errors.New("fail")
	} else {
		return rsp.BlockHash, rsp.Error
	}
}

func CurrentBlockHash() (common.Uint256, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetCurrentBlockHashReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return common.Uint256{}, err
	}
	if rsp, ok := result.(*lactor.GetCurrentBlockHashRsp); !ok {
		return common.Uint256{}, errors.New("fail")
	} else {
		return rsp.BlockHash, rsp.Error
	}
}

func GetBlockFromStore(hash common.Uint256) (*types.Block, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetBlockByHashReq{hash}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil, err
	}
	if rsp, ok := result.(*lactor.GetBlockByHashRsp); !ok {
		return nil, errors.New("fail")
	} else {
		return rsp.Block, rsp.Error
	}
}

func BlockHeight() (uint32, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetCurrentBlockHeightReq{}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	if rsp, ok := result.(*lactor.GetCurrentBlockHeightRsp); !ok {
		return 0, errors.New("fail")
	} else {
		return rsp.Height, rsp.Error
	}
}

func GetTransaction(hash common.Uint256) (*types.Transaction, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetTransactionReq{hash}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil, err
	}
	if rsp, ok := result.(*lactor.GetTransactionRsp); !ok {
		return nil, errors.New("fail")
	} else {
		return rsp.Tx, rsp.Error
	}
}

func GetStorageItem(codeHash common.Address, key []byte) ([]byte, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetStorageItemReq{CodeHash: codeHash, Key: key}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil, err
	}
	if rsp, ok := result.(*lactor.GetStorageItemRsp); !ok {
		return nil, errors.New("fail")
	} else {
		return rsp.Value, rsp.Error
	}
}

func GetContractStateFromStore(hash common.Address) (*payload.DeployCode, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetContractStateReq{hash}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil, err
	}
	if rsp, ok := result.(*lactor.GetContractStateRsp); !ok {
		log.Error(ERR_ACTOR_COMM,"GetContractStateRsp")
		return nil, errors.New("fail")
	} else {
		return rsp.ContractState, rsp.Error
	}
}

func GetBlockHeightByTxHashFromStore(hash common.Uint256) (uint32, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetTransactionWithHeightReq{hash}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return 0, err
	}
	if rsp, ok := result.(*lactor.GetTransactionWithHeightRsp); !ok {
		return 0, errors.New("fail")
	} else {
		return rsp.Height, rsp.Error
	}
}

func AddBlock(block *types.Block) error {
	future := defLedgerPid.RequestFuture(&lactor.AddBlockReq{block}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return err
	}
	if rsp, ok := result.(*lactor.AddBlockRsp); !ok {
		return errors.New("fail")
	} else {
		return rsp.Error
	}
}

func PreExecuteContract(tx *types.Transaction) ([]interface{}, error) {
	future := defLedgerPid.RequestFuture(&lactor.PreExecuteContractReq{tx}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil, err
	}
	if rsp, ok := result.(*lactor.PreExecuteContractRsp); !ok {
		return nil, errors.New("fail")
	} else {
		return rsp.Result, rsp.Error
	}
}

func GetEventNotifyByTxHash(txHash common.Uint256) ([]*event.NotifyEventInfo, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetEventNotifyByTxReq{txHash}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil, err
	}
	if rsp, ok := result.(*lactor.GetEventNotifyByTxRsp); !ok {
		return nil,errors.New("fail")
	}else {
		return rsp.Notifies,rsp.Error
	}
}

func GetEventNotifyByHeight(height uint32) ([]common.Uint256, error) {
	future := defLedgerPid.RequestFuture(&lactor.GetEventNotifyByBlockReq{height}, REQ_TIMEOUT*time.Second)
	result, err := future.Result()
	if err != nil {
		log.Errorf(ERR_ACTOR_COMM, err)
		return nil, err
	}
	if rsp, ok := result.(*lactor.GetEventNotifyByBlockRsp); !ok {
		return nil,errors.New("fail")
	}else {
		return rsp.TxHashes,rsp.Error
	}
}
