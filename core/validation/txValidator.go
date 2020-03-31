package validation

import (
	"errors"
	"fmt"
	"github.com/dad-go/common"
	"github.com/dad-go/common/log"
	"github.com/dad-go/core/asset"
	"github.com/dad-go/core/ledger"
	tx "github.com/dad-go/core/transaction"
	"github.com/dad-go/core/transaction/payload"
	"github.com/dad-go/core/transaction/utxo"
	"github.com/dad-go/crypto"
	. "github.com/dad-go/errors"
	"math"
)

// VerifyTransaction verifys received single transaction
func VerifyTransaction(Tx *tx.Transaction) ErrCode {

	if err := CheckDuplicateInput(Tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return ErrDuplicateInput
	}

	if err := CheckAssetPrecision(Tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return ErrAssetPrecision
	}

	if err := CheckTransactionBalance(Tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return ErrTransactionBalance
	}

	if err := CheckAttributeProgram(Tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return ErrAttributeProgram
	}

	if err := CheckTransactionContracts(Tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return ErrTransactionContracts
	}

	if err := CheckTransactionPayload(Tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return ErrTransactionPayload
	}

	return ErrNoError
}

// VerifyTransactionWithBlock verifys a transaction with current transaction pool in memory
func VerifyTransactionWithBlock(TxPool []*tx.Transaction) error {
	//initial
	txnlist := make(map[common.Uint256]*tx.Transaction, 0)
	var txPoolInputs []string
	//sum all inputs in TxPool
	for _, Tx := range TxPool {
		for _, UTXOinput := range Tx.UTXOInputs {
			txPoolInputs = append(txPoolInputs, UTXOinput.ToString())
		}
	}
	//start check
	for _, txn := range TxPool {
		//1.check weather have duplicate transaction.
		if _, exist := txnlist[txn.Hash()]; exist {
			return errors.New("[VerifyTransactionWithBlock], duplicate transaction exist in block.")
		} else {
			txnlist[txn.Hash()] = txn
		}
		//2.check Duplicate Utxo input
		if err := CheckDuplicateUtxoInBlock(txn, txPoolInputs); err != nil {
			return err
		}
		//3.check issue amount
		switch txn.TxType {
		case tx.IssueAsset:
			//TODO: use delta mode to improve performance
			results := txn.GetMergedAssetIDValueFromOutputs()
			for k, _ := range results {
				//Get the Asset amount when RegisterAsseted.
				trx, err := tx.TxStore.GetTransaction(k)
				if trx.TxType != tx.RegisterAsset {
					return errors.New("[VerifyTransaction], TxType is illegal.")
				}
				AssetReg := trx.Payload.(*payload.RegisterAsset)

				//Get the amount has been issued of this assetID
				var quantity_issued common.Fixed64
				if AssetReg.Amount < common.Fixed64(0) {
					continue
				} else {
					quantity_issued, err = tx.TxStore.GetQuantityIssued(k)
					if err != nil {
						return errors.New("[VerifyTransaction], GetQuantityIssued failed.")
					}
				}

				//calc the amounts in txPool which are also IssueAsset
				var txPoolAmounts common.Fixed64
				for _, t := range TxPool {
					if t.TxType == tx.IssueAsset {
						outputResult := t.GetMergedAssetIDValueFromOutputs()
						for txidInPool, txValueInPool := range outputResult {
							if txidInPool == k {
								txPoolAmounts = txPoolAmounts + txValueInPool
							}
						}
					}
				}

				//calc weather out off the amount when Registed.
				//AssetReg.Amount : amount when RegisterAsset of this assedID
				//quantity_issued : amount has been issued of this assedID
				//txPoolAmounts   : amount in transactionPool of this assedID of issue transaction.
				if AssetReg.Amount-quantity_issued < txPoolAmounts {
					return errors.New("[VerifyTransaction], Amount check error.")
				}
			}
		}

	}

	return nil
}

// VerifyTransactionWithLedger verifys a transaction with history transaction in ledger
func VerifyTransactionWithLedger(Tx *tx.Transaction, ledger *ledger.Ledger) ErrCode {
	if IsDoubleSpend(Tx, ledger) {
		log.Info("[VerifyTransactionWithLedger] IsDoubleSpend check faild.")
		return ErrDoubleSpend
	}
	if exist := ledger.Store.IsTxHashDuplicate(Tx.Hash()); exist {
		log.Info("[VerifyTransactionWithLedger] duplicate transaction check faild.")
		return ErrTxHashDuplicate
	}
	return ErrNoError
}

//validate the transaction of duplicate UTXO input
func CheckDuplicateInput(tx *tx.Transaction) error {
	if len(tx.UTXOInputs) == 0 {
		return nil
	}
	for i, utxoin := range tx.UTXOInputs {
		for j := 0; j < i; j++ {
			if utxoin.ReferTxID == tx.UTXOInputs[j].ReferTxID && utxoin.ReferTxOutputIndex == tx.UTXOInputs[j].ReferTxOutputIndex {
				return errors.New("invalid transaction")
			}
		}
	}
	return nil
}

func CheckDuplicateUtxoInBlock(tx *tx.Transaction, txPoolInputs []string) error {
	var txInputs []string
	for _, t := range tx.UTXOInputs {
		txInputs = append(txInputs, t.ToString())
	}
	for _, i := range txInputs {
		for _, j := range txPoolInputs {
			if i == j {
				return errors.New("Duplicated UTXO inputs found in tx pool")
			}
		}
	}
	return nil
}

func IsDoubleSpend(tx *tx.Transaction, ledger *ledger.Ledger) bool {
	return ledger.IsDoubleSpend(tx)
}

func CheckAssetPrecision(Tx *tx.Transaction) error {
	if len(Tx.Outputs) == 0 {
		return nil
	}
	assetOutputs := make(map[common.Uint256][]*utxo.TxOutput, len(Tx.Outputs))

	for _, v := range Tx.Outputs {
		assetOutputs[v.AssetID] = append(assetOutputs[v.AssetID], v)
	}
	for k, outputs := range assetOutputs {
		asset, err := ledger.DefaultLedger.GetAsset(k)
		if err != nil {
			log.Debugf("The asset not exist in local blockchain. %x", k)
			return errors.New("The asset not exist in local blockchain.")
		}
		precision := asset.Precision
		for _, output := range outputs {
			if checkAmountPrecise(output.Value, precision) {
				log.Debugf("output.Value", output.Value, "precision", precision)
				return errors.New("The precision of asset is incorrect.")
			}
		}
	}
	return nil
}

func CheckTransactionBalance(Tx *tx.Transaction) error {
	if Tx.SystemFee < 0 {
		return errors.New("[CheckTransactionBalance] Invalide transaction SystemFee.")
	}
	for _, v := range Tx.Outputs {
		if v.Value <= common.Fixed64(0) {
			return errors.New("[CheckTransactionBalance] Invalide transaction UTXO output.")
		}
	}
	networkfee, err := Tx.GetNetworkFee()
	if err != nil {
		return errors.New(fmt.Sprintf("[CheckTransactionBalance] GetNetworkFee failed. with err=",err))
	}
	results, err := Tx.GetTransactionResults()
	if err != nil {
		return err
	}
	switch Tx.TxType {
	case tx.IssueAsset:
		for k, v := range results {
			if k.CompareTo(tx.ONGTokenID) == 0 {
				if Tx.GetSysFee().GetData() != Tx.SystemFee.GetData() {
					return errors.New(fmt.Sprintf("AssetID %x in Transfer transactions %x ,SystemFee/NetworkFee Not equal.", k, Tx.Hash()))
				}
				if v.GetData() != Tx.GetSysFee().GetData()+ networkfee.GetData() {
					return errors.New(fmt.Sprintf("AssetID %x in Transfer transactions %x ,SystemFee/NetworkFee Not equal.", k, Tx.Hash()))
				}
			}
		}
	case tx.Claim:
		return nil
	default:
		for k, v := range results {
			if k.CompareTo(tx.ONGTokenID) == 0 {
				if Tx.GetSysFee().GetData() != Tx.SystemFee.GetData() {
					return errors.New(fmt.Sprintf("AssetID %x in Transfer transactions %x ,SystemFee/NetworkFee Not equal.", k, Tx.Hash()))
				}
				if v.GetData() != Tx.GetSysFee().GetData() + networkfee.GetData() {
					log.Debug(fmt.Sprintf("AssetID %x in Transfer transactions %x ,SystemFee Not equal.", k, Tx.Hash()))
					return errors.New(fmt.Sprintf("AssetID %x in Transfer transactions %x ,SystemFee Not equal.", k, Tx.Hash()))
				}
			} else {
				if v != 0 {
					log.Debug(fmt.Sprintf("AssetID %x in Transfer transactions %x , Input/output UTXO not equal.", k, Tx.Hash()))
					return errors.New(fmt.Sprintf("AssetID %x in Transfer transactions %x , Input/output UTXO not equal.", k, Tx.Hash()))
				}
			}
		}
	}
	return nil
}

func CheckAttributeProgram(Tx *tx.Transaction) error {
	//TODO: implement CheckAttributeProgram
	return nil
}

func CheckTransactionContracts(Tx *tx.Transaction) error {
	err := VerifySignableDataSignature(Tx)
	if err == nil {
		return nil
	}

	err = VerifySignableDataProgramHashes(Tx)
	return err
}

func checkAmountPrecise(amount common.Fixed64, precision byte) bool {
	return amount.GetData()%int64(math.Pow(10, 8-float64(precision))) != 0
}

func CheckTransactionPayload(Tx *tx.Transaction) error {

	switch pld := Tx.Payload.(type) {
	case *payload.BookKeeper:
		//Todo: validate bookKeeper Cert
		_ = pld.Cert
		bookKeepers, _, _ := ledger.DefaultLedger.Store.GetBookKeeperList()

		index := crypto.ContainPubKey(pld.Issuer, bookKeepers)
		if index == -1 {
			return errors.New("The issuer isn't bookekeeper, can't add other in bookkeepers list.")
		}
		return nil
	case *payload.RegisterAsset:
		if pld.Asset.Precision < asset.MinPrecision || pld.Asset.Precision > asset.MaxPrecision {
			return errors.New("Invalide asset Precision.")
		}
		if checkAmountPrecise(pld.Amount, pld.Asset.Precision) {
			return errors.New("Invalide asset value,out of precise.")
		}
	case *payload.IssueAsset:
	case *payload.TransferAsset:
	case *payload.BookKeeping:
	case *payload.PrivacyPayload:
	case *payload.Record:
	case *payload.DeployCode:
	case *payload.InvokeCode:
	case *payload.DataFile:
	case *payload.Claim:
		claims := Tx.Payload.(*payload.Claim).Claims
		if isDoubleClaim(claims) {
			return errors.New("[CheckTransactionPayload], Invalid double claim")
		}
		result, err := Tx.GetTransactionResults()
		if err != nil {
			return errors.New("[CheckTransactionPayload], Invalid transaction results")
		}
		var claimAmount common.Fixed64
		for k, v := range result {
			if k.CompareTo(tx.ONGTokenID) == 0 {
				claimAmount += v
			}
		}
		amount, err := ledger.CalculateBouns(claims, false)
		if err != nil {
			return errors.New(fmt.Sprintf("[CheckTransactionPayload], CalculateBouns error:%v", err))
		}
		if -amount.GetData() != claimAmount.GetData() {
			return errors.New(fmt.Sprintf("[CheckTransactionPayload], claims amount error claimed amount =%d, actual=%d", claimAmount.GetData(), amount.GetData()))
		}
	default:
		return errors.New("[txValidator],invalidate transaction payload type.")
	}
	return nil
}

func isDoubleClaim(claims []*utxo.UTXOTxInput) bool {
	for i := 0; i < len(claims); i++ {
		for j := 0; i < i; j++ {
			if claims[i] == claims[j] {
				return true
			}
		}
	}
	return false
}
