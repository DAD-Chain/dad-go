package jsonrpc

import (
	"bytes"
	"encoding/hex"
	. "github.com/dad-go/common"
	"github.com/dad-go/core/types"
	. "github.com/dad-go/errors"
	. "github.com/dad-go/http/base/common"
	. "github.com/dad-go/http/base/rpc"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/core/states"
	"github.com/dad-go/common/config"
	"math/rand"
	"fmt"
	"encoding/base64"
	"os"
)

func GetBestBlockHash(params []interface{}) map[string]interface{} {
	hash := ledger.DefaultLedger.Blockchain.CurrentBlockHash()
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
		hash, err = ledger.DefaultLedger.Store.GetBlockHash(index)
		if err != nil {
			return dad-goRpcUnknownBlock
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

	block, err := ledger.DefaultLedger.Store.GetBlock(hash)
	if err != nil {
		return dad-goRpcUnknownBlock
	}

	blockHead := &BlockHead{
		Version:          block.Header.Version,
		PrevBlockHash:    ToHexString(block.Header.PrevBlockHash.ToArray()),
		TransactionsRoot: ToHexString(block.Header.TransactionsRoot.ToArray()),
		BlockRoot:        ToHexString(block.Header.BlockRoot.ToArray()),
		StateRoot:        ToHexString(block.Header.StateRoot.ToArray()),
		Timestamp:        block.Header.Timestamp,
		Height:           block.Header.Height,
		ConsensusData:    block.Header.ConsensusData,
		NextBookKeeper:   ToHexString(block.Header.NextBookKeeper[:]),
		Program: ProgramInfo{
			Code:      ToHexString(block.Header.Program.Code),
			Parameter: ToHexString(block.Header.Program.Parameter),
		},
		Hash: ToHexString(hash.ToArray()),
	}

	trans := make([]*Transactions, len(block.Transactions))
	for i := 0; i < len(block.Transactions); i++ {
		trans[i] = TransArryByteToHexString(block.Transactions[i])
	}

	b := BlockInfo{
		Hash:         ToHexString(hash.ToArray()),
		BlockData:    blockHead,
		Transactions: trans,
	}
	return dad-goRpc(b)
}

func GetBlockCount(params []interface{}) map[string]interface{} {
	return dad-goRpc(ledger.DefaultLedger.Blockchain.BlockHeight + 1)
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
		hash, err := ledger.DefaultLedger.Store.GetBlockHash(height)
		if err != nil {
			return dad-goRpcUnknownBlock
		}
		return dad-goRpc(fmt.Sprintf("%016x", hash))
	default:
		return dad-goRpcInvalidParameter
	}
}

func GetConnectionCount(params []interface{}) map[string]interface{} {
	return dad-goRpc(CNoder.GetConnectionCnt())
}

func GetRawMemPool(params []interface{}) map[string]interface{} {
	txs := []*Transactions{}
	txpool, _ := CNoder.GetTxnPool(false)
	for _, t := range txpool {
		txs = append(txs, TransArryByteToHexString(t))
	}
	if len(txs) == 0 {
		return dad-goRpcNil
	}
	return dad-goRpc(txs)
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
		tx, err := ledger.DefaultLedger.Store.GetTransaction(hash)
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
	item, err := ledger.DefaultLedger.Store.GetStorageItem(&states.StorageKey{CodeHash: codeHash, Key: key})
	if err != nil {
		return dad-goRpcInternalError
	}
	return dad-goRpc(ToHexString(item.Value))
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

func GetBalance(params []interface{}) map[string]interface{} {
	if len(params) < 2 {
		return dad-goRpcNil
	}

	addr, ok := params[0].(string)
	if !ok {
		return dad-goRpcInvalidParameter
	}
	assetId, ok := params[1].(string)
	if !ok {
		return dad-goRpcInvalidParameter
	}

	programHash, err := ToScriptHash(addr)
	if err != nil {
		return dad-goRpcInvalidParameter
	}
	account, err := ledger.DefaultLedger.Store.GetAccount(programHash)
	if err != nil {
		return dad-goRpcAccountNotFound
	}
	c, err := HexToBytes(assetId)
	if err != nil {
		return dad-goRpcInvalidParameter
	}
	ass, err := Uint256ParseFromBytes(c)
	if err != nil {
		return dad-goRpcInvalidParameter
	}

	for _, v := range account.Balances {
		if v.AssetId.CompareTo(ass) == 0 {
			return dad-goRpc(v.Amount.GetData())
		}
	}

	return dad-goRpcNil
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
		if err := ledger.DefaultLedger.Blockchain.AddBlock(&block); err != nil {
			return dad-goRpcInvalidBlock
		}
		if err := CNoder.Xmit(&block); err != nil {
			return dad-goRpcInternalError
		}
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpcSuccess
}

func GetVersion(params []interface{}) map[string]interface{} {
	return dad-goRpc(config.Version)
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
		tx, err := ledger.DefaultLedger.Store.GetTransaction(hash)
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
		tx, err := ledger.DefaultLedger.Store.GetTransaction(hash)
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
