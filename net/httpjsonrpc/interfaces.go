package httpjsonrpc

import (
	"github.com/dad-go/account"
	. "github.com/dad-go/common"
	"github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	"github.com/dad-go/core/ledger"
	tx "github.com/dad-go/core/transaction"
	. "github.com/dad-go/errors"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

const (
	RANDBYTELEN = 4
)

func TransArryByteToHexString(ptx *tx.Transaction) *Transactions {

	trans := new(Transactions)
	trans.TxType = ptx.TxType
	trans.PayloadVersion = ptx.PayloadVersion
	trans.Payload = TransPayloadToHex(ptx.Payload)

	n := 0
	trans.Attributes = make([]TxAttributeInfo, len(ptx.Attributes))
	for _, v := range ptx.Attributes {
		trans.Attributes[n].Usage = v.Usage
		trans.Attributes[n].Data = ToHexString(v.Data)
		n++
	}

	n = 0
	trans.UTXOInputs = make([]UTXOTxInputInfo, len(ptx.UTXOInputs))
	for _, v := range ptx.UTXOInputs {
		trans.UTXOInputs[n].ReferTxID = ToHexString(v.ReferTxID.ToArray())
		trans.UTXOInputs[n].ReferTxOutputIndex = v.ReferTxOutputIndex
		n++
	}

	n = 0
	trans.BalanceInputs = make([]BalanceTxInputInfo, len(ptx.BalanceInputs))
	for _, v := range ptx.BalanceInputs {
		trans.BalanceInputs[n].AssetID = ToHexString(v.AssetID.ToArray())
		trans.BalanceInputs[n].Value = v.Value
		trans.BalanceInputs[n].ProgramHash = ToHexString(v.ProgramHash.ToArray())
		n++
	}

	n = 0
	trans.Outputs = make([]TxoutputInfo, len(ptx.Outputs))
	for _, v := range ptx.Outputs {
		trans.Outputs[n].AssetID = ToHexString(v.AssetID.ToArray())
		trans.Outputs[n].Value = v.Value
		trans.Outputs[n].ProgramHash = ToHexString(v.ProgramHash.ToArray())
		n++
	}

	n = 0
	trans.Programs = make([]ProgramInfo, len(ptx.Programs))
	for _, v := range ptx.Programs {
		trans.Programs[n].Code = ToHexString(v.Code)
		trans.Programs[n].Parameter = ToHexString(v.Parameter)
		n++
	}

	n = 0
	trans.AssetOutputs = make([]TxoutputMap, len(ptx.AssetOutputs))
	for k, v := range ptx.AssetOutputs {
		trans.AssetOutputs[n].Key = k
		trans.AssetOutputs[n].Txout = make([]TxoutputInfo, len(v))
		for m := 0; m < len(v); m++ {
			trans.AssetOutputs[n].Txout[m].AssetID = ToHexString(v[m].AssetID.ToArray())
			trans.AssetOutputs[n].Txout[m].Value = v[m].Value
			trans.AssetOutputs[n].Txout[m].ProgramHash = ToHexString(v[m].ProgramHash.ToArray())
		}
		n += 1
	}

	n = 0
	trans.AssetInputAmount = make([]AmountMap, len(ptx.AssetInputAmount))
	for k, v := range ptx.AssetInputAmount {
		trans.AssetInputAmount[n].Key = k
		trans.AssetInputAmount[n].Value = v
		n += 1
	}

	n = 0
	trans.AssetOutputAmount = make([]AmountMap, len(ptx.AssetOutputAmount))
	for k, v := range ptx.AssetOutputAmount {
		trans.AssetInputAmount[n].Key = k
		trans.AssetInputAmount[n].Value = v
		n += 1
	}

	mhash := ptx.Hash()
	trans.Hash = ToHexString(mhash.ToArray())

	return trans
}
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
func getBestBlockHash(params []interface{}) map[string]interface{} {
	hash := ledger.DefaultLedger.Blockchain.CurrentBlockHash()
	return dad-goRpc(ToHexString(hash.ToArray()))
}

// Input JSON string examples for getblock method as following:
//   {"jsonrpc": "2.0", "method": "getblock", "params": [1], "id": 0}
//   {"jsonrpc": "2.0", "method": "getblock", "params": ["aabbcc.."], "id": 0}
func getBlock(params []interface{}) map[string]interface{} {
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
		Version:          block.Blockdata.Version,
		PrevBlockHash:    ToHexString(block.Blockdata.PrevBlockHash.ToArray()),
		TransactionsRoot: ToHexString(block.Blockdata.TransactionsRoot.ToArray()),
		BlockRoot:        ToHexString(block.Blockdata.BlockRoot.ToArray()),
		Timestamp:        block.Blockdata.Timestamp,
		Height:           block.Blockdata.Height,
		ConsensusData:    block.Blockdata.ConsensusData,
		NextBookKeeper:   ToHexString(block.Blockdata.NextBookKeeper.ToArray()),
		Program: ProgramInfo{
			Code:      ToHexString(block.Blockdata.Program.Code),
			Parameter: ToHexString(block.Blockdata.Program.Parameter),
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

func getBlockCount(params []interface{}) map[string]interface{} {
	return dad-goRpc(ledger.DefaultLedger.Blockchain.BlockHeight + 1)
}

// A JSON example for getblockhash method as following:
//   {"jsonrpc": "2.0", "method": "getblockhash", "params": [1], "id": 0}
func getBlockHash(params []interface{}) map[string]interface{} {
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

func getConnectionCount(params []interface{}) map[string]interface{} {
	return dad-goRpc(node.GetConnectionCnt())
}

func getRawMemPool(params []interface{}) map[string]interface{} {
	txs := []*Transactions{}
	txpool := node.GetTxnPool(false)
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
func getRawTransaction(params []interface{}) map[string]interface{} {
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

// A JSON example for sendrawtransaction method as following:
//   {"jsonrpc": "2.0", "method": "sendrawtransaction", "params": ["raw transactioin in hex"], "id": 0}
func sendRawTransaction(params []interface{}) map[string]interface{} {
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
		var txn tx.Transaction
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

func getUnspendOutput(params []interface{}) map[string]interface{} {
	if len(params) < 2 {
		return dad-goRpcNil
	}
	var programHash Uint160
	var assetHash Uint256
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, err := hex.DecodeString(str)
		if err != nil {
			return dad-goRpcInvalidParameter
		}
		if err := programHash.Deserialize(bytes.NewReader(hex)); err != nil {
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
		if err := assetHash.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidHash
		}
	default:
		return dad-goRpcInvalidParameter
	}
	type TxOutputInfo struct {
		AssetID     string
		Value       int64
		ProgramHash string
	}
	outputs := make(map[string]*TxOutputInfo)
	height := ledger.DefaultLedger.GetLocalBlockChainHeight()
	var i uint32
	// construct global UTXO table
	for i = 0; i <= height; i++ {
		block, err := ledger.DefaultLedger.GetBlockWithHeight(i)
		if err != nil {
			return dad-goRpcInternalError
		}
		// skip the bookkeeping transaction
		for _, t := range block.Transactions[1:] {
			// skip the register transaction
			if t.TxType == tx.RegisterAsset {
				continue
			}
			txHash := t.Hash()
			txHashHex := ToHexString(txHash.ToArray())
			for i, output := range t.Outputs {
				if output.AssetID.CompareTo(assetHash) == 0 &&
					output.ProgramHash.CompareTo(programHash) == 0 {
					key := txHashHex + ":" + strconv.Itoa(i)
					asset := ToHexString(output.AssetID.ToArray())
					pHash := ToHexString(output.ProgramHash.ToArray())
					value := int64(output.Value)
					info := &TxOutputInfo{
						asset,
						value,
						pHash,
					}
					outputs[key] = info
				}
			}
		}
	}
	// delete spent output from global UTXO table
	height = ledger.DefaultLedger.GetLocalBlockChainHeight()
	for i = 0; i <= height; i++ {
		block, err := ledger.DefaultLedger.GetBlockWithHeight(i)
		if err != nil {
			return dad-goRpcInternalError
		}
		// skip the bookkeeping transaction
		for _, t := range block.Transactions[1:] {
			// skip the register transaction
			if t.TxType == tx.RegisterAsset {
				continue
			}
			for _, input := range t.UTXOInputs {
				refer := ToHexString(input.ReferTxID.ToArray())
				index := strconv.Itoa(int(input.ReferTxOutputIndex))
				key := refer + ":" + index
				delete(outputs, key)
			}
		}
	}
	return dad-goRpc(outputs)
}

func getTxout(params []interface{}) map[string]interface{} {
	//TODO
	return dad-goRpcUnsupported
}

// A JSON example for submitblock method as following:
//   {"jsonrpc": "2.0", "method": "submitblock", "params": ["raw block in hex"], "id": 0}
func submitBlock(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	switch params[0].(type) {
	case string:
		str := params[0].(string)
		hex, _ := hex.DecodeString(str)
		var block ledger.Block
		if err := block.Deserialize(bytes.NewReader(hex)); err != nil {
			return dad-goRpcInvalidBlock
		}
		if err := ledger.DefaultLedger.Blockchain.AddBlock(&block); err != nil {
			return dad-goRpcInvalidBlock
		}
		if err := node.Xmit(&block); err != nil {
			return dad-goRpcInternalError
		}
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpcSuccess
}

func getNeighbor(params []interface{}) map[string]interface{} {
	addr, _ := node.GetNeighborAddrs()
	return dad-goRpc(addr)
}

func getNodeState(params []interface{}) map[string]interface{} {
	n := NodeInfo{
		State:    uint(node.GetState()),
		Time:     node.GetTime(),
		Port:     node.GetPort(),
		ID:       node.GetID(),
		Version:  node.Version(),
		Services: node.Services(),
		Relay:    node.GetRelay(),
		Height:   node.GetHeight(),
		TxnCnt:   node.GetTxnCnt(),
		RxTxnCnt: node.GetRxTxnCnt(),
	}
	return dad-goRpc(n)
}

func startConsensus(params []interface{}) map[string]interface{} {
	if err := dBFT.Start(); err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpcSuccess
}

func stopConsensus(params []interface{}) map[string]interface{} {
	if err := dBFT.Halt(); err != nil {
		return dad-goRpcFailed
	}
	return dad-goRpcSuccess
}

func sendSampleTransaction(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcNil
	}
	var txType string
	switch params[0].(type) {
	case string:
		txType = params[0].(string)
	default:
		return dad-goRpcInvalidParameter
	}

	issuer, err := account.NewAccount()
	if err != nil {
		return dad-goRpc("Failed to create account")
	}
	admin := issuer

	rbuf := make([]byte, RANDBYTELEN)
	rand.Read(rbuf)
	switch string(txType) {
	case "perf":
		num := 1
		if len(params) == 2 {
			switch params[1].(type) {
			case float64:
				num = int(params[1].(float64))
			}
		}
		for i := 0; i < num; i++ {
			regTx := NewRegTx(ToHexString(rbuf), i, admin, issuer)
			SignTx(admin, regTx)
			VerifyAndSendTx(regTx)
		}
		return dad-goRpc(fmt.Sprintf("%d transaction(s) was sent", num))
	default:
		return dad-goRpc("Invalid transacion type")
	}
}

func setDebugInfo(params []interface{}) map[string]interface{} {
	if len(params) < 1 {
		return dad-goRpcInvalidParameter
	}
	switch params[0].(type) {
	case float64:
		level := params[0].(float64)
		if err := log.Log.SetDebugLevel(int(level)); err != nil {
			return dad-goRpcInvalidParameter
		}
	default:
		return dad-goRpcInvalidParameter
	}
	return dad-goRpcSuccess
}

func getVersion(params []interface{}) map[string]interface{} {
	return dad-goRpc(config.Version)
}

func uploadDataFile(params []interface{}) map[string]interface{} {
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
	f, err := os.OpenFile(tmpname, os.O_WRONLY | os.O_CREATE, 0664)
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

func regDataFile(params []interface{}) map[string]interface{} {
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
		var txn tx.Transaction
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

func catDataRecord(params []interface{}) map[string]interface{} {
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

func getDataFile(params []interface{}) map[string]interface{} {
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
