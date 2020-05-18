package rpc

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	. "github.com/dad-go/common"
	"github.com/dad-go/common/config"
	"github.com/dad-go/core/types"
	. "github.com/dad-go/errors"
	. "github.com/dad-go/http/base/actor"
	. "github.com/dad-go/http/base/common"
	"math/rand"
	"os"
)

func GetBestBlockHash(params []interface{}) map[string]interface{} {
	hash, err := CurrentBlockHash()
	if err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpc(ToHexString(hash.ToArray()))
}

// Input JSON string examples for getblock method as following:
//   {"jsonrpc": "2.0", "method": "getblock", "params": [1], "id": 0}
//   {"jsonrpc": "2.0", "method": "getblock", "params": ["aabbcc.."], "id": 0}
func GetBlock(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	var err error
	var hash Uint256
	switch (params[0]).(type) {
	// block height
	case float64:
		index := uint32(params[0].(float64))
		hash, err = GetBlockHashFromStore(index)
		if err != nil {
			return dad-goRpcUnknownBlock
		}
		if hash.CompareTo(Uint256{}) == 0{
			return dad-goRpcInvalidParameter
		}
		// block hash
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		if err := hash.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidTransaction
		}
	default:
		return dad-goRpcInvalidParameter
	}

	block, err := GetBlockFromStore(hash)
	if err != nil {
		return dad-goRpcUnknownBlock
	}
	if block.Header == nil{
		return dad-goRpcUnknownBlock
	}
	return dad-goRpc(GetBlockInfo(block))
}

func GetBlockCount(params []interface{}) map[string]interface{} {
	height, err := BlockHeight()
	if err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpc(height + 1)
}

// A JSON example for getblockhash method as following:
//   {"jsonrpc": "2.0", "method": "getblockhash", "params": [1], "id": 0}
func GetBlockHash(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	switch params[0].(type) {
	case float64:
		height := uint32(params[0].(float64))
		hash, err := GetBlockHashFromStore(height)
		if err != nil {
			return dad-goRpcUnknownBlock
		}
		return dad-goRpc(fmt.Sprintf("%016x", hash))
	default:
		return dad-goRpcInvalidParameter
	}
}

func GetConnectionCount(params []interface{}) map[string]interface{} {
	count, err := GetConnectionCnt()
	if err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpc(count)
}

func GetRawMemPool(params []interface{}) map[string]interface{} {
	txs := []*Transactions{}
	txpool, _ := GetTxsFromPool(false)
	for _, t := range txpool {
		txs = append(txs, TransArryByteToHexString(t))
	}
	if len(txs) == 0 {
		return dad-goRpcNil
	}
	return dad-goRpc(txs)
}
func GetMemPoolTxState(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		var hash Uint256
		err = hash.Deserialize(bytes.NewReader(hex))
		if err != nil {
			return dad-goRpcInvalidTransaction
		}
		txEntry, err := GetTxFromPool(hash)
		if err != nil {
			return dad-goRpcUnknownTransaction
		}
		tran := TransArryByteToHexString(txEntry.Tx)
		attrs := []TXNAttrInfo{}
		for _, t := range txEntry.Attrs {
			attrs = append(attrs, TXNAttrInfo{t.Height, int(t.Type), int(t.ErrCode)})
		}
		info := TXNEntryInfo{*tran, int64(txEntry.Fee), attrs}
		return dad-goRpc(info)
	default:
		return dad-goRpcInvalidParameter
	}
}

// A JSON example for getrawtransaction method as following:
//   {"jsonrpc": "2.0", "method": "getrawtransaction", "params": ["transactioin hash in hex"], "id": 0}
func GetRawTransaction(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		var hash Uint256
		err = hash.Deserialize(bytes.NewReader(hex))
		if err != nil {
			return dad-goRpcInvalidTransaction
		}
		tx, err := GetTransaction(hash) //ledger.DefaultLedger.Store.GetTransaction(hash)
		if err != nil {
			return dad-goRpcUnknownTransaction
		}
		tran := TransArryByteToHexString(tx)
		return dad-goRpc(tran)
	default:
		return dad-goRpcInvalidParameter
	}
}

//   {"jsonrpc": "2.0", "method": "getstorage", "params": ["code hash", "key"], "id": 0}
func GetStorage(params []interface{}) map[string]interface{} {
	if len(params) < 2 {
		return dad-goRpcNil
	}

	var codeHash Uint160
	var key []byte
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		if err := codeHash.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidHash
		}
	default:
		return dad-goRpcInvalidParameter
	}

	switch params[1].(type) {
	case string:
		str := params[1].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		key = hex
	default:
		return dad-goRpcInvalidParameter
	}
	value, err := GetStorageItem(codeHash, key)
	if err != nil {
		return dad-goRpcInternalError
	}
	return dad-goRpc(ToHexString(value))
}

// A JSON example for sendrawtransaction method as following:
//   {"jsonrpc": "2.0", "method": "sendrawtransaction", "params": ["raw transactioin in hex"], "id": 0}
func SendRawTransaction(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	var hash Uint256
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		var txn types.Transaction
		if err := txn.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidTransaction
		}
		hash = txn.Hash()
		if errCode := VerifyAndSendTx(&txn); errCode != ErrNoError {
			return dad-goRpc(errCode.Error())
		}
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpc(ToHexString(hash.ToArray()))
}

// A JSON example for submitblock method as following:
//   {"jsonrpc": "2.0", "method": "submitblock", "params": ["raw block in hex"], "id": 0}
func SubmitBlock(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, _ := hex.DecodeString(str)
		var block types.Block
		if err := block.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidBlock
		}
		if err := AddBlock(&block); err != nil {
			return dad-goRpcInvalidBlock
		}
		if err := Xmit(&block); err != nil {
			return dad-goRpcInternalError
		}
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpcSuccess
}

func GetNodeVersion(params []interface{}) map[string]interface{} {
	return dad-goRpc(config.Parameters.Version)
}

func GetSystemFee(params []interface{}) map[string]interface{} {
	return dad-goRpc(config.Parameters.SystemFee)
}

func GetContractState(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		var hash Uint160
		if err := hash.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidParameter
		}
		contract, err := GetContractStateFromStore(hash)
		if err != nil || contract == nil{
			return dad-goRpcInternalError
		}
		return dad-goRpc(TransPayloadToHex(contract))
	default:
		return dad-goRpcInvalidParameter
	}
}

func UploadDataFile(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}

	rbuf := make([]byte, 4)
	rand.Read(rbuf)
	tmpname := hex.EncodeToString(rbuf)

	str := params[0].(string)

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return dad-goRpcInvalidParameter
	}
	f, err := os.OpenFile(tmpname, os.O_WRONLY|os.O_CREATE, 0664)
	if err != nil {
		return dad-goRpcIOError
	}
	defer f.Close()
	f.Write(data)

	refpath, err := AddFileIPFS(tmpname, true)
	if err != nil {
		return dad-goRpcAPIError
	}

	return dad-goRpc(refpath)

}
func GetSmartCodeEvent(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}

	switch (params[0]).(type) {
	// block height
	case float64:
		height := uint32(params[0].(float64))
		//TODO resp
		return dad-goRpc(map[string]interface{}{"Height": height})
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpcInvalidParameter
}

func GetTxBlockHeight(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}

	switch (params[0]).(type) {
	// tx hash
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		var hash Uint256
		if err := hash.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidParameter
		}
		//TODO resp
		return dad-goRpc(map[string]interface{}{"Height": 0})
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpcInvalidParameter
}
func RegDataFile(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	var hash Uint256
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		var txn types.Transaction
		if err := txn.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidTransaction
		}

		hash = txn.Hash()
		if errCode := VerifyAndSendTx(&txn); errCode != ErrNoError {
			return dad-goRpcInternalError
		}
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpc(ToHexString(hash.ToArray()))
}

func CatDataRecord(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		b, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		var hash Uint256
		err = hash.Deserialize(bytes.NewReader(b))
		if err != nil {
			return dad-goRpcInvalidTransaction
		}
		tx, err := GetTransaction(hash)
		if err != nil {
			return dad-goRpcUnknownTransaction
		}
		tran := TransArryByteToHexString(tx)
		info := tran.Payload.(*DataFileInfo)
		//ref := string(record.RecordData[:])
		return dad-goRpc(info)
	default:
		return dad-goRpcInvalidParameter
	}
}

func GetDataFile(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		var hash Uint256
		err = hash.Deserialize(bytes.NewReader(hex))
		if err != nil {
			return dad-goRpcInvalidTransaction
		}
		tx, err := GetTransaction(hash)
		if err != nil {
			return dad-goRpcUnknownTransaction
		}

		tran := TransArryByteToHexString(tx)
		info := tran.Payload.(*DataFileInfo)

		err = GetFileIPFS(info.IPFSPath, info.Filename)
		if err != nil {
			return dad-goRpcAPIError
		}
		//TODO: shoud return download address
		return dad-goRpcSuccess
	default:
		return dad-goRpcInvalidParameter
	}
}
