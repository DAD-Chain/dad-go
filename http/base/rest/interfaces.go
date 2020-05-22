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

package rest

import (
	"bytes"
	. "github.com/dad-go/common"
	"github.com/dad-go/common/config"
	"github.com/dad-go/core/types"
	. "github.com/dad-go/errors"
	Err "github.com/dad-go/http/base/error"
	"strconv"
	"github.com/dad-go/core/payload"
	"github.com/dad-go/common/log"
	. "github.com/dad-go/http/base/common"
	. "github.com/dad-go/http/base/actor"
	"math/big"
	"github.com/dad-go/core/genesis"
)

const TlsPort int = 443

type ApiServer interface {
	Start() error
	Stop()
}

//Node
func GetGenerateBlockTime(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	resp["Result"] = config.DEFAULTGENBLOCKTIME
	return resp
}
func GetConnectionCount(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	count, err := GetConnectionCnt()
	if err != nil {
		return rspInternalError
	}
	resp["Result"] = count
	return resp
}

//Block
func GetBlockHeight(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	height,err := BlockHeight()
	if err != nil{
		return rspInternalError
	}
	resp["Result"] = height
	return resp
}
func GetBlockHash(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	param := cmd["Height"].(string)
	if len(param) == 0 {
		return rspInvalidParams
	}
	height, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return rspInvalidParams
	}
	hash, err := GetBlockHashFromStore(uint32(height))
	if err != nil {
		return rspInvalidParams
	}
	resp["Result"] = ToHexString(hash.ToArray())
	return resp
}

func GetBlockTransactions(block *types.Block) interface{} {
	trans := make([]string, len(block.Transactions))
	for i := 0; i < len(block.Transactions); i++ {
		h := block.Transactions[i].Hash()
		trans[i] = ToHexString(h.ToArray())
	}
	hash := block.Hash()
	type BlockTransactions struct {
		Hash         string
		Height       uint32
		Transactions []string
	}
	b := BlockTransactions{
		Hash:         ToHexString(hash.ToArray()),
		Height:       block.Header.Height,
		Transactions: trans,
	}
	return b
}
func getBlock(hash Uint256, getTxBytes bool) (interface{}, int64) {
	block, err := GetBlockFromStore(hash)
	if err != nil {
		return nil, Err.UNKNOWN_BLOCK
	}
	if block.Header == nil {
		return nil, Err.UNKNOWN_BLOCK
	}
	if getTxBytes {
		w := bytes.NewBuffer(nil)
		block.Serialize(w)
		return ToHexString(w.Bytes()), Err.SUCCESS
	}
	return GetBlockInfo(block), Err.SUCCESS
}
func GetBlockByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	param := cmd["Hash"].(string)
	if len(param) == 0 {
		return rspInvalidParams
	}
	var getTxBytes = false
	if raw, ok := cmd["Raw"].(string); ok && raw == "1" {
		getTxBytes = true
	}
	var hash Uint256
	hex, err := HexToBytes(param)
	if err != nil {
		return rspInvalidParams
	}
	if err := hash.Deserialize(bytes.NewReader(hex)); err != nil {
		return rspInvalidTx
	}
	resp["Result"], resp["Error"] = getBlock(hash, getTxBytes)
	return resp
}

func GetBlockHeightByTxHash(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	param := cmd["Hash"].(string)
	if len(param) == 0 {
		return rspInvalidParams
	}
	var hash Uint256
	hex, err := HexToBytes(param)
	if err != nil {
		return rspInvalidParams
	}
	if err := hash.Deserialize(bytes.NewReader(hex)); err != nil {
		return rspInvalidTx
	}
	height,err := GetBlockHeightByTxHashFromStore(hash)
	if err != nil {
		return rspInternalError
	}
	resp["Result"] = height
	return resp
}
func GetBlockTxsByHeight(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess

	param := cmd["Height"].(string)
	if len(param) == 0 {
		return rspInvalidParams
	}
	height, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return rspInvalidParams
	}
	index := uint32(height)
	hash, err := GetBlockHashFromStore(index)
	if err != nil {
		return rspUnkownBlock
	}
	if hash.CompareTo(Uint256{}) == 0{
		return rspInvalidParams
	}
	block, err := GetBlockFromStore(hash)
	if err != nil {
		return rspUnkownBlock
	}
	resp["Result"] = GetBlockTransactions(block)
	return resp
}
func GetBlockByHeight(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess

	param := cmd["Height"].(string)
	if len(param) == 0 {
		return rspInvalidParams
	}
	var getTxBytes bool = false
	if raw, ok := cmd["Raw"].(string); ok && raw == "1" {
		getTxBytes = true
	}
	height, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return rspInvalidParams
	}
	index := uint32(height)
	hash, err := GetBlockHashFromStore(index)
	if err != nil {
		return rspUnkownBlock
	}
	resp["Result"], resp["Error"] = getBlock(hash, getTxBytes)
	return resp
}


//Transaction
func GetTransactionByHash(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess

	str := cmd["Hash"].(string)
	bys, err := HexToBytes(str)
	if err != nil {
		return rspInvalidParams
	}
	var hash Uint256
	err = hash.Deserialize(bytes.NewReader(bys))
	if err != nil {
		return rspInvalidTx
	}
	tx, err := GetTransaction(hash)
	if err != nil {
		return rspUnkownTx
	}
	if tx == nil {
		return rspUnkownTx
	}
	if raw, ok := cmd["Raw"].(string); ok && raw == "1" {
		w := bytes.NewBuffer(nil)
		tx.Serialize(w)
		resp["Result"] = ToHexString(w.Bytes())
		return resp
	}
	tran := TransArryByteToHexString(tx)
	resp["Result"] = tran
	return resp
}
func SendRawTransaction(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess

	str, ok := cmd["Data"].(string)
	if !ok {
		return rspInvalidParams
	}
	bys, err := HexToBytes(str)
	if err != nil {
		return rspInvalidParams
	}
	var txn types.Transaction
	if err := txn.Deserialize(bytes.NewReader(bys)); err != nil {
		return rspInvalidTx
	}
	if txn.TxType == types.Invoke {
		if preExec, ok := cmd["PreExec"].(string); ok && preExec == "1" {
			log.Tracef("PreExec SMARTCODE")
			if _, ok := txn.Payload.(*payload.InvokeCode); ok {
				resp["Result"], err = PreExecuteContract(&txn)
				if err != nil {
					log.Error(err)
					return rspSmartCodeError
				}
				return resp
			}
		}
	}
	var hash Uint256
	hash = txn.Hash()
	if errCode := VerifyAndSendTx(&txn); errCode != ErrNoError {
		resp["Error"] = int64(errCode)
		return resp
	}
	resp["Result"] = ToHexString(hash.ToArray())

	if txn.TxType == types.Invoke {
		if userid, ok := cmd["Userid"].(string); ok && len(userid) > 0 {
			resp["Userid"] = userid
		}
	}
	return resp
}

func GetSmartCodeEventByHeight(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess

	param := cmd["Height"].(string)
	if len(param) == 0 {
		return rspInvalidParams
	}
	height, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return rspInvalidParams
	}
	index := uint32(height)
	txs, err := GetEventNotifyByHeight(index)
	if err != nil {
		return rspInvalidParams
	}
	var txhexs []string
	for _, v := range txs {
		txhexs = append(txhexs, ToHexString(v.ToArray()))
	}
	resp["Result"] = txhexs
	return resp
}

func GetSmartCodeEventByTxHash(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess

	str := cmd["Hash"].(string)
	bys, err := HexToBytes(str)
	if err != nil {
		return rspInvalidParams
	}
	var hash Uint256
	err = hash.Deserialize(bytes.NewReader(bys))
	if err != nil {
		return rspInvalidTx
	}
	eventInfos, err := GetEventNotifyByTxHash(hash)
	if err != nil {
		return rspInvalidParams
	}
	var evs []map[string]interface{}
	for _, v := range eventInfos {
		evs = append(evs, map[string]interface{}{"CodeHash": v.CodeHash,
			"States": v.States,
			"Container": v.Container})
	}
	resp["Result"] = evs
	return resp
}

func GetContractState(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	str := cmd["Hash"].(string)
	bys, err := HexToBytes(str)
	if err != nil {
		return rspInvalidParams
	}
	var hash Address
	err = hash.Deserialize(bytes.NewReader(bys))
	if err != nil {
		return rspInvalidParams
	}
	contract, err := GetContractStateFromStore(hash)
	if err != nil || contract == nil {
		return rspInternalError
	}
	resp["Result"] = TransPayloadToHex(contract)
	return resp
}

func GetStorage(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	str := cmd["Hash"].(string)
	bys, err := HexToBytes(str)
	if err != nil {
		return rspInvalidParams
	}
	var hash Address
	err = hash.Deserialize(bytes.NewReader(bys))
	if err != nil {
		return rspInvalidParams
	}
	key := cmd["Key"].(string)
	item, err := HexToBytes(key)
	if err != nil {
		return rspInvalidParams
	}
	log.Info("[GetStorage] ",str,key)
	value, err := GetStorageItem(hash,item)
	if err != nil || value == nil {
		return rspInternalError
	}
	resp["Result"] = ToHexString(value)
	return resp
}
func GetBalance(cmd map[string]interface{}) map[string]interface{} {
	resp := rspSuccess
	addrBase58 := cmd["Addr"].(string)
	address, err := AddressFromBase58(addrBase58)
	if err != nil {
		return rspInvalidParams
	}
	ont := new(big.Int)
	ong := new(big.Int)

	ontBalance, err := GetStorageItem(genesis.OntContractAddress, address.ToArray())
	if err != nil {
		log.Errorf("GetOntBalanceOf GetStorageItem ont address:%s error:%s", address, err)
		return rspInternalError
	}
	if ontBalance != nil {
		ont.SetBytes(ontBalance)
	}
	rsp := &BalanceOfRsp{
		Ont: ont.String(),
		Ong: ong.String(),
	}
	resp["Result"] = rsp
	return resp
}