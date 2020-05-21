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

package validation

import (
	"errors"
	"fmt"

	"github.com/dad-go/common/log"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/types"
	"github.com/dad-go/crypto"
	. "github.com/dad-go/errors"
	"github.com/dad-go/common"
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
	address := make(map[common.Address]bool, len(tx.Sigs))
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

			address[types.AddressFromPubKey(sig.PubKeys[0])] = true
		} else {
			if err := crypto.VerifyMultiSignature(hash[:], sig.PubKeys, m, sig.SigData); err != nil {
				return err
			}

			addr, _ := types.AddressFromMultiPubKeys(sig.PubKeys, m)
			address[addr] = true
		}
	}

	// check all payers in address
	for _, fee := range tx.Fee {
		if address[fee.Payer] == false {
			return errors.New("signature missing for payer: " + common.ToHexString(fee.Payer.ToArray()))
		}
	}

	return nil
}

func checkTransactionPayload(tx *types.Transaction) error {

	switch pld := tx.Payload.(type) {
	case *payload.DeployCode:
		return nil
	case *payload.InvokeCode:
		return nil
	case *payload.BookKeeping:
		return nil
	default:
		return errors.New(fmt.Sprint("[txValidator], unimplemented transaction payload type.", pld))
	}
	return nil
}
