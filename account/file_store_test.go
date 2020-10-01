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
package account

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/core/types"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func genAccountData() (*AccountData, *keypair.ProtectedKey) {
	var acc = new(AccountData)
	prvkey, pubkey, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	ta := types.AddressFromPubKey(pubkey)
	address := ta.ToBase58()
	password := []byte("123456")
	prvSectet, _ := keypair.EncryptPrivateKey(prvkey, address, password)
	h := sha256.Sum256(password)
	acc.SetKeyPair(prvSectet)
	acc.SigSch = "SHA256withECDSA"
	acc.PubKey = hex.EncodeToString(keypair.SerializePublicKey(pubkey))
	acc.PassHash = hex.EncodeToString(h[:])
	return acc, prvSectet
}

func TestAccountData(t *testing.T) {
	acc, prvSectet := genAccountData()
	assert.NotNil(t, acc)
	assert.Equal(t, acc.Address, acc.ProtectedKey.Address)
	assert.Equal(t, prvSectet, acc.GetKeyPair())
	assert.True(t, acc.VerifyPassword([]byte("123456")))
}

func TestWalletStorage(t *testing.T) {
	defer func() {
		os.Remove(WALLET_FILENAME)
		os.RemoveAll("Log/")
	}()

	wallet := new(WalletData)
	wallet.Inititalize()
	wallet.Save(WALLET_FILENAME)
	walletReadFromFile := new(WalletData)
	walletReadFromFile.Load(WALLET_FILENAME)
	assert.Equal(t, walletReadFromFile, wallet)
	accout, _ := genAccountData()
	wallet.AddAccount(accout)
	wallet.AddAccount(accout)
	wallet.Save(WALLET_FILENAME)
	walletReadFromFile.Load(WALLET_FILENAME)
	assert.Equal(t, walletReadFromFile, wallet)
	wallet.DelAccount(2)
	assert.Equal(t, 1, len(wallet.Accounts))
	assert.Panics(t, func() { wallet.DelAccount(2) })
	wallet.Save(WALLET_FILENAME)
	walletReadFromFile.Load(WALLET_FILENAME)
	assert.Equal(t, walletReadFromFile, wallet)
	defaultAccount := wallet.GetDefaultAccount()
	assert.Equal(t, defaultAccount, accout)
}
