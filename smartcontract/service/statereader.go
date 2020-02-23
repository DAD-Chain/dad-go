package service

import (
	"github.com/dad-go/common"
	"github.com/dad-go/core/contract"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/core/signature"
	tx "github.com/dad-go/core/transaction"
	"github.com/dad-go/crypto"
	"github.com/dad-go/errors"
	vm "github.com/dad-go/vm/neovm"
	"github.com/dad-go/vm/neovm/types"
	"math/big"
	"github.com/dad-go/core/states"
	"github.com/dad-go/core/transaction/utxo"
)

type StateReader struct {
	serviceMap map[string]func(*vm.ExecutionEngine) (bool, error)
}

func NewStateReader() *StateReader {
	var stateReader StateReader
	stateReader.serviceMap = make(map[string]func(*vm.ExecutionEngine) (bool, error), 0)
	stateReader.Register("Neo.Runtime.GetTrigger", stateReader.RuntimeGetTrigger)
	stateReader.Register("Neo.Runtime.CheckWitness", stateReader.RuntimeCheckWitness)
	stateReader.Register("Neo.Runtime.Notify", stateReader.RuntimeNotify)
	stateReader.Register("Neo.Runtime.Log", stateReader.RuntimeLog)

	stateReader.Register("Neo.Blockchain.GetHeight", stateReader.BlockChainGetHeight)
	stateReader.Register("Neo.Blockchain.GetHeader", stateReader.BlockChainGetHeader)
	stateReader.Register("Neo.Blockchain.GetBlock", stateReader.BlockChainGetBlock)
	stateReader.Register("Neo.Blockchain.GetTransaction", stateReader.BlockChainGetTransaction)
	stateReader.Register("Neo.Blockchain.GetAccount", stateReader.BlockChainGetAccount)
	stateReader.Register("Neo.Blockchain.GetValidators", stateReader.BlockChainGetValidators)
	stateReader.Register("Neo.Blockchain.GetAsset", stateReader.BlockChainGetAsset)
	stateReader.Register("Neo.Blockchain.GetContract", stateReader.GetContract)

	stateReader.Register("Neo.Header.GetHash", stateReader.HeaderGetHash)
	stateReader.Register("Neo.Header.GetVersion", stateReader.HeaderGetVersion)
	stateReader.Register("Neo.Header.GetPrevHash", stateReader.HeaderGetPrevHash)
	stateReader.Register("Neo.Header.GetMerkleRoot", stateReader.HeaderGetMerkleRoot)
	stateReader.Register("Neo.Header.GetTimestamp", stateReader.HeaderGetTimestamp)
	stateReader.Register("Neo.Header.GetConsensusData", stateReader.HeaderGetConsensusData)
	stateReader.Register("Neo.Header.GetNextConsensus", stateReader.HeaderGetNextConsensus)

	stateReader.Register("Neo.Block.GetTransactionCount", stateReader.BlockGetTransactionCount)
	stateReader.Register("Neo.Block.GetTransactions", stateReader.BlockGetTransactions)
	stateReader.Register("Neo.Block.GetTransaction", stateReader.BlockGetTransaction)

	stateReader.Register("Neo.Transaction.GetHash", stateReader.TransactionGetHash)
	stateReader.Register("Neo.Transaction.GetType", stateReader.TransactionGetType)
	stateReader.Register("Neo.Transaction.GetAttributes", stateReader.TransactionGetAttributes)
	stateReader.Register("Neo.Transaction.GetInputs", stateReader.TransactionGetInputs)
	stateReader.Register("Neo.Transaction.GetOutputs", stateReader.TransactionGetOutputs)
	stateReader.Register("Neo.Transaction.GetReferences", stateReader.TransactionGetReferences)

	stateReader.Register("Neo.Attribute.GetUsage", stateReader.AttributeGetUsage)
	stateReader.Register("Neo.Attribute.GetData", stateReader.AttributeGetData)

	stateReader.Register("Neo.Input.GetHash", stateReader.InputGetHash)
	stateReader.Register("Neo.Input.GetIndex", stateReader.InputGetIndex)

	stateReader.Register("Neo.Output.GetAssetId", stateReader.OutputGetAssetId)
	stateReader.Register("Neo.Output.GetValue", stateReader.OutputGetValue)
	stateReader.Register("Neo.Output.GetScriptHash", stateReader.OutputGetCodeHash)

	stateReader.Register("Neo.Account.GetScriptHash", stateReader.AccountGetCodeHash)
	stateReader.Register("Neo.Account.GetBalance", stateReader.AccountGetBalance)

	stateReader.Register("Neo.Asset.GetAssetId", stateReader.AssetGetAssetId)
	stateReader.Register("Neo.Asset.GetAssetType", stateReader.AssetGetAssetType)
	stateReader.Register("Neo.Asset.GetAmount", stateReader.AssetGetAmount)
	stateReader.Register("Neo.Asset.GetAvailable", stateReader.AssetGetAvailable)
	stateReader.Register("Neo.Asset.GetPrecision", stateReader.AssetGetPrecision)
	stateReader.Register("Neo.Asset.GetOwner", stateReader.AssetGetOwner)
	stateReader.Register("Neo.Asset.GetAdmin", stateReader.AssetGetAdmin)

	return &stateReader
}

func (s *StateReader) Register(methodad-gome string, handler func(*vm.ExecutionEngine) (bool, error)) bool {
	if _, ok := s.serviceMap[methodad-gome]; ok {
		return false
	}
	s.serviceMap[methodad-gome] = handler
	return true
}

func (s *StateReader) GetServiceMap() map[string]func(*vm.ExecutionEngine) (bool, error) {
	return s.serviceMap
}

func (s *StateReader) RuntimeGetTrigger(e *vm.ExecutionEngine) (bool, error) {
	return true, nil
}

func (s *StateReader) RuntimeNotify(e *vm.ExecutionEngine) (bool, error) {
	vm.PopStackItem(e)
	return true, nil
}

func (s *StateReader) RuntimeLog(e *vm.ExecutionEngine) (bool, error) {
	return true, nil
}

func (s *StateReader) CheckWitnessHash(engine *vm.ExecutionEngine, programHash common.Uint160) (bool, error) {
	hashForVerifying, err := engine.GetCodeContainer().(signature.SignableData).GetProgramHashes()
	if err != nil {
		return false, err
	}
	return contains(hashForVerifying, programHash), nil
}

func (s *StateReader) CheckWitnessPublicKey(engine *vm.ExecutionEngine, publicKey *crypto.PubKey) (bool, error) {
	c, err := contract.CreateSignatureRedeemScript(publicKey)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "[CheckWitnessPublicKey] CreateSignatureRedeemScript error!")
	}
	h, err := common.ToCodeHash(c)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "[CheckWitnessPublicKey] ToCodeHash error!")
	}
	return s.CheckWitnessHash(engine, h)
}

func (s *StateReader) RuntimeCheckWitness(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[RuntimeCheckWitness] Too few input parameters ")
	}
	data := vm.PopByteArray(e)
	var (
		result bool
		err error
	)
	if len(data) == 20 {
		program, err := common.Uint160ParseFromBytes(data)
		if err != nil {
			return false, err
		}
		result, err = s.CheckWitnessHash(e, program)
	} else if len(data) == 33 {
		publicKey, err := crypto.DecodePoint(data)
		if err != nil {
			return false, err
		}
		result, err = s.CheckWitnessPublicKey(e, publicKey)
	} else {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "[RuntimeCheckWitness] data invalid.")
	}
	if err != nil {
		return false, err
	}
	vm.PushData(e, result)
	return true, nil
}

func (s *StateReader) BlockChainGetHeight(e *vm.ExecutionEngine) (bool, error) {
	var i uint32
	if ledger.DefaultLedger == nil {
		i = 0
	} else {
		i = ledger.DefaultLedger.Store.GetHeight()
	}
	vm.PushData(e, i)
	return true, nil
}

func (s *StateReader) BlockChainGetHeader(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[BlockChainGetHeader] Too few input parameters ")
	}
	data := vm.PopByteArray(e)
	var (
		header *ledger.Header
		err error
	)
	l := len(data)
	if l <= 5 {
		b := new(big.Int)
		height := uint32(b.SetBytes(common.BytesReverse(data)).Int64())
		if ledger.DefaultLedger != nil {
			hash, err := ledger.DefaultLedger.Store.GetBlockHash(height)
			if err != nil {
				return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetHeader] GetBlockHash error!.")
			}
			header, err = ledger.DefaultLedger.Store.GetHeader(hash)
			if err != nil {
				return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetHeader] GetHeader error!.")
			}
		}
	} else if l == 32 {
		hash, _ := common.Uint256ParseFromBytes(data)
		if ledger.DefaultLedger != nil {
			header, err = ledger.DefaultLedger.Store.GetHeader(hash)
			if err != nil {
				return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetHeader] GetHeader error!.")
			}
		}
	} else {
		return false, errors.NewErr("[BlockChainGetHeader] data invalid.")
	}
	vm.PushData(e, header)
	return true, nil
}

func (s *StateReader) BlockChainGetBlock(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[BlockChainGetBlock] Too few input parameters ")
	}
	data := vm.PopByteArray(e)
	var (
		block *ledger.Block
	)
	l := len(data)
	if l <= 5 {
		b := new(big.Int)
		height := uint32(b.SetBytes(common.BytesReverse(data)).Int64())
		if ledger.DefaultLedger != nil {
			hash, err := ledger.DefaultLedger.Store.GetBlockHash(height)
			if err != nil {
				return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetBlock] GetBlockHash error!.")
			}
			block, err = ledger.DefaultLedger.Store.GetBlock(hash)
			if err != nil {
				return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetBlock] GetBlock error!.")
			}
		}
	} else if l == 32 {
		hash, err := common.Uint256ParseFromBytes(data)
		if err != nil {
			return false, err
		}
		block, err = ledger.DefaultLedger.Store.GetBlock(hash)
		if err != nil {
			return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetBlock] GetBlock error!.")
		}
	} else {
		return false, errors.NewErr("[BlockChainGetBlock] data invalid.")
	}
	vm.PushData(e, block)
	return true, nil
}

func (s *StateReader) BlockChainGetTransaction(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[BlockChainGetTransaction] Too few input parameters ")
	}
	d := vm.PopByteArray(e)
	hash, err := common.Uint256ParseFromBytes(d)
	if err != nil {
		return false, err
	}
	t, err := ledger.DefaultLedger.Store.GetTransaction(hash)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetTransaction] GetTransaction error!")
	}

	vm.PushData(e, t)
	return true, nil
}

func (s *StateReader) BlockChainGetAccount(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[BlockChainGetAccount] Too few input parameters ")
	}
	d := vm.PopByteArray(e)
	hash, err := common.Uint160ParseFromBytes(d)
	if err != nil {
		return false, err
	}
	account, err := ledger.DefaultLedger.Store.GetAccount(hash)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetAccount] BlockChainGetAccount error!")
	}
	vm.PushData(e, account)
	return true, nil
}

func (s *StateReader) BlockChainGetAsset(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[BlockChainGetAsset] Too few input parameters ")
	}
	d := vm.PopByteArray(e)
	hash, err := common.Uint256ParseFromBytes(d)
	if err != nil {
		return false, err
	}
	assetState, err := ledger.DefaultLedger.Store.GetAsset(hash)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "[BlockChainGetAsset] GetAsset error!")
	}
	vm.PushData(e, assetState)
	return true, nil
}

func (s *StateReader) GetContract(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[GetContract] Too few input parameters ")
	}
	hashByte := vm.PopByteArray(e)
	hash, err := common.Uint160ParseFromBytes(hashByte)
	if err != nil {
		return false, err
	}
	item, err := ledger.DefaultLedger.Store.GetContract(hash)
	if err != nil {
		return false, errors.NewDetailErr(err, errors.ErrNoCode, "[GetContract] GetAsset error!")
	}
	vm.PushData(e, item)
	return true, nil
}

func (s *StateReader) BlockChainGetValidators(e *vm.ExecutionEngine) (bool, error) {
	bookKeeperList, err := ledger.GetValidators([]*tx.Transaction{})
	if err != nil {
		return false, err
	}
	var pkList []types.StackItemInterface
	for _, v := range bookKeeperList {
		pk, err := v.EncodePoint(true)
		if err != nil {
			return false, err
		}
		pkList = append(pkList, types.NewByteArray(pk))
	}
	vm.PushData(e, pkList)
	return true, nil
}

func (s *StateReader) HeaderGetHash(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[HeaderGetHash] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[HeaderGetHash] Pop blockdata nil!")
	}
	h := d.(*ledger.Header).Blockdata.Hash()
	vm.PushData(e, h.ToArray())
	return true, nil
}

func (s *StateReader) HeaderGetVersion(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[HeaderGetVersion] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[HeaderGetVersion] Pop blockdata nil!")
	}
	version := d.(*ledger.Header).Blockdata.Version
	vm.PushData(e, version)
	return true, nil
}

func (s *StateReader) HeaderGetPrevHash(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[HeaderGetPrevHash] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[HeaderGetPrevHash] Pop blockdata nil!")
	}
	preHash := d.(*ledger.Header).Blockdata.PrevBlockHash
	vm.PushData(e, preHash.ToArray())
	return true, nil
}

func (s *StateReader) HeaderGetMerkleRoot(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[HeaderGetMerkleRoot] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[HeaderGetMerkleRoot] Pop blockdata nil!")
	}
	root := d.(*ledger.Header).Blockdata.TransactionsRoot
	vm.PushData(e, root.ToArray())
	return true, nil
}

func (s *StateReader) HeaderGetTimestamp(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[HeaderGetTimestamp] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[HeaderGetTimestamp] Pop blockdata nil!")
	}
	timeStamp := d.(*ledger.Header).Blockdata.Timestamp
	vm.PushData(e, timeStamp)
	return true, nil
}

func (s *StateReader) HeaderGetConsensusData(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[HeaderGetConsensusData] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[HeaderGetConsensusData] Pop blockdata nil!")
	}
	consensusData := d.(*ledger.Header).Blockdata.ConsensusData
	vm.PushData(e, consensusData)
	return true, nil
}

func (s *StateReader) HeaderGetNextConsensus(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[HeaderGetNextConsensus] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[HeaderGetNextConsensus] Pop blockdata nil!")
	}
	nextBookKeeper := d.(*ledger.Header).Blockdata.NextBookKeeper
	vm.PushData(e, nextBookKeeper.ToArray())
	return true, nil
}

func (s *StateReader) BlockGetTransactionCount(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[BlockGetTransactionCount] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[BlockGetTransactionCount] Pop blockdata nil!")
	}
	transactions := d.(*ledger.Block).Transactions
	vm.PushData(e, len(transactions))
	return true, nil
}

func (s *StateReader) BlockGetTransactions(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[BlockGetTransactions] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[BlockGetTransactions] Pop blockdata nil!")
	}
	transactions := d.(*ledger.Block).Transactions
	transactionList := make([]types.StackItemInterface, 0)
	for _, v := range transactions {
		transactionList = append(transactionList, types.NewInteropInterface(v))
	}
	vm.PushData(e, transactionList)
	return true, nil
}

func (s *StateReader) BlockGetTransaction(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 2 {
		return false, errors.NewErr("[BlockGetTransaction] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[BlockGetTransaction] Pop transactions nil!")
	}
	index := vm.PopInt(e)
	if index < 0 {
		return false, errors.NewErr("[BlockGetTransaction] Pop index invalid!")
	}
	transactions := d.(*ledger.Block).Transactions
	if index >= len(transactions) {
		return false, errors.NewErr("[BlockGetTransaction] index invalid!")
	}
	vm.PushData(e, transactions[index])
	return true, nil
}

func (s *StateReader) TransactionGetHash(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[TransactionGetHash] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[TransactionGetHash] Pop transaction nil!")
	}
	txHash := d.(*tx.Transaction).Hash()
	vm.PushData(e, txHash.ToArray())
	return true, nil
}

func (s *StateReader) TransactionGetType(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[TransactionGetType] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[TransactionGetType] Pop transaction nil!")
	}
	txType := d.(*tx.Transaction).TxType
	vm.PushData(e, int(txType))
	return true, nil
}

func (s *StateReader) TransactionGetAttributes(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[TransactionGetAttributes] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[TransactionGetAttributes] Pop transaction nil!")
	}
	attributes := d.(*tx.Transaction).Attributes
	attributList := make([]types.StackItemInterface, 0)
	for _, v := range attributes {
		attributList = append(attributList, types.NewInteropInterface(v))
	}
	vm.PushData(e, attributList)
	return true, nil
}

func (s *StateReader) TransactionGetInputs(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[TransactionGetInputs] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[TransactionGetInputs] Pop transaction nil!")
	}
	inputs := d.(*tx.Transaction).UTXOInputs
	inputList := make([]types.StackItemInterface, 0)
	for _, v := range inputs {
		inputList = append(inputList, types.NewInteropInterface(v))
	}
	vm.PushData(e, inputList)
	return true, nil
}

func (s *StateReader) TransactionGetOutputs(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[TransactionGetOutputs] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[TransactionGetOutputs] Pop transaction nil!")
	}
	outputs := d.(*tx.Transaction).Outputs
	outputList := make([]types.StackItemInterface, 0)
	for _, v := range outputs {
		outputList = append(outputList, types.NewInteropInterface(v))
	}
	vm.PushData(e, outputList)
	return true, nil
}

func (s *StateReader) TransactionGetReferences(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[TransactionGetReferences] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[TransactionGetReferences] Pop transaction nil!")
	}
	references, err := d.(*tx.Transaction).GetReference()
	referenceList := make([]types.StackItemInterface, 0)
	for _, v := range references {
		referenceList = append(referenceList, types.NewInteropInterface(v))
	}
	vm.PushData(e, referenceList)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *StateReader) AttributeGetUsage(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AttributeGetUsage] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AttributeGetUsage] Pop txAttribute nil!")
	}
	attribute := d.(*tx.TxAttribute)
	vm.PushData(e, int(attribute.Usage))
	return true, nil
}

func (s *StateReader) AttributeGetData(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AttributeGetData] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AttributeGetData] Pop txAttribute nil!")
	}
	attribute := d.(*tx.TxAttribute)
	vm.PushData(e, attribute.Data)
	return true, nil
}

func (s *StateReader) InputGetHash(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[InputGetHash] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[InputGetHash] Pop utxoTxInput nil!")
	}
	input := d.(*utxo.UTXOTxInput)
	vm.PushData(e, input.ReferTxID.ToArray())
	return true, nil
}

func (s *StateReader) InputGetIndex(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[InputGetIndex] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[InputGetIndex] Pop utxoTxInput nil!")
	}
	input := d.(*utxo.UTXOTxInput)
	vm.PushData(e, input.ReferTxOutputIndex)
	return true, nil
}

func (s *StateReader) OutputGetAssetId(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[OutputGetAssetId] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[OutputGetAssetId] Pop txOutput nil!")
	}
	output := d.(*utxo.TxOutput)
	vm.PushData(e, output.AssetID.ToArray())
	return true, nil
}

func (s *StateReader) OutputGetValue(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[OutputGetValue] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[OutputGetValue] Pop txOutput nil!")
	}
	output := d.(*utxo.TxOutput)
	vm.PushData(e, output.Value.GetData())
	return true, nil
}

func (s *StateReader) OutputGetCodeHash(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[OutputGetCodeHash] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[OutputGetCodeHash] Pop txOutput nil!")
	}
	output := d.(*utxo.TxOutput)
	vm.PushData(e, output.ProgramHash.ToArray())
	return true, nil
}

func (s *StateReader) AccountGetCodeHash(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AccountGetCodeHash] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AccountGetCodeHash] Pop accountState nil!")
	}
	accountState := d.(*states.AccountState).ProgramHash
	vm.PushData(e, accountState.ToArray())
	return true, nil
}

func (s *StateReader) AccountGetBalance(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 2 {
		return false, errors.NewErr("[AccountGetBalance] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AccountGetBalance] Pop accountState nil!")
	}
	accountState := d.(*states.AccountState)
	assetIdByte := vm.PopByteArray(e)
	assetId, err := common.Uint256ParseFromBytes(assetIdByte)
	if err != nil {
		return false, err
	}
	balance := common.Fixed64(0)
	if v, ok := accountState.Balances[assetId]; ok {
		balance = v
	}
	vm.PushData(e, balance.GetData())
	return true, nil
}

func (s *StateReader) AssetGetAssetId(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AssetGetAssetId] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AssetGetAssetId] Pop assetState nil!")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(e, assetState.AssetId.ToArray())
	return true, nil
}

func (s *StateReader) AssetGetAssetType(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AssetGetAssetType] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AssetGetAssetType] Pop assetState nil!")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(e, int(assetState.AssetType))
	return true, nil
}

func (s *StateReader) AssetGetAmount(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AssetGetAmount] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AssetGetAmount] Pop assetState nil!")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(e, assetState.Amount.GetData())
	return true, nil
}

func (s *StateReader) AssetGetAvailable(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AssetGetAvailable] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AssetGetAvailable] Pop assetState nil!")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(e, assetState.Available.GetData())
	return true, nil
}

func (s *StateReader) AssetGetPrecision(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AssetGetPrecision] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AssetGetPrecision] Pop assetState nil!")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(e, int(assetState.Precision))
	return true, nil
}

func (s *StateReader) AssetGetOwner(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AssetGetOwner] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AssetGetOwner] Pop assetState nil!")
	}
	assetState := d.(*states.AssetState)
	owner, err := assetState.Owner.EncodePoint(true)
	if err != nil {
		return false, err
	}
	vm.PushData(e, owner)
	return true, nil
}

func (s *StateReader) AssetGetAdmin(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[AssetGetAdmin] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[AssetGetAdmin] Pop assetState nil!")
	}
	assetState := d.(*states.AssetState)
	vm.PushData(e, assetState.Admin.ToArray())
	return true, nil
}

func (s *StateReader) ContractGetCode(e *vm.ExecutionEngine) (bool, error) {
	if vm.EvaluationStackCount(e) < 1 {
		return false, errors.NewErr("[ContractGetCode] Too few input parameters ")
	}
	d := vm.PopInteropInterface(e)
	if d == nil {
		return false, errors.NewErr("[ContractGetCode] Pop contractState nil!")
	}
	assetState := d.(*states.ContractState)
	vm.PushData(e, assetState.Code.Code)
	return true, nil
}

func (s *StateReader) StorageGetContext(e *vm.ExecutionEngine) (bool, error) {
	codeHash, err := common.Uint160ParseFromBytes(e.CurrentContext().GetCodeHash())
	if err != nil {
		return false, err
	}
	vm.PushData(e, NewStorageContext(codeHash))
	return true, nil
}
