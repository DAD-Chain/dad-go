package validation

import (
	"errors"
	"fmt"

	"github.com/dad-go/common/log"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/types"
	"github.com/dad-go/core/utils"
	"github.com/dad-go/crypto"
	. "github.com/dad-go/errors"
)

// VerifyTransaction verifys received single transaction
func VerifyTransaction(tx *types.Transaction) ErrCode {
	if err := checkTransactionSignatures(tx); err != nil {
		log.Info("transaction verify error:", err)
		return ErrTransactionContracts
	}

	if err := checkTransactionPayload(tx); err != nil {
		log.Warn("[VerifyTransaction],", err)
		return ErrTransactionPayload
	}

	return ErrNoError
}

func VerifyTransactionWithLedger(tx *types.Transaction, ledger *ledger.Ledger) ErrCode {
	//TODO: replay check
	return ErrNoError
}


func checkTransactionSignatures(tx *types.Transaction) error {
	hash := tx.Hash()
	address := make(map[types.Address]bool, len(tx.Sigs))
	for _, sig := range tx.Sigs {
		m := int(sig.M)
		n := len(sig.PubKeys)
		s := len(sig.SigData)

		if n > 24 || s < m || m > n {
			return errors.New("wrong tx sig param length")
		}

		if n == 1 {
			err := crypto.Verify(*sig.PubKeys[0], hash[:], sig.SigData[0])
			if err != nil {
				return err
			}

			address[utils.AddressFromPubKey(sig.PubKeys[0])] = true
		} else {
			if err := crypto.VerifyMultiSignature(hash[:], sig.PubKeys, m, sig.SigData); err != nil {
				return err
			}

			addr, _ := utils.AddressFromMultiPubKeys(sig.PubKeys, m)
			address[addr] = true
		}
	}

	// check all payers in address
	for _, fee := range tx.Fee {
		if address[fee.Payer] == false {
			return errors.New("signature missing for payer: " + fee.Payer.ToHexString())
		}
	}

	return nil
}

func checkTransactionPayload(tx *types.Transaction) error {

	switch pld := tx.Payload.(type) {
	case *payload.InvokeCode:
		return nil
	default:
		return errors.New(fmt.Sprint("[txValidator], unimplemented transaction payload type.", pld))
	}
	return nil
}
