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
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/config"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go-crypto/keypair"
	s "github.com/ontio/dad-go-crypto/signature"
)

type Account struct {
	PrivateKey keypair.PrivateKey
	PublicKey  keypair.PublicKey
	Address    common.Address
}

func NewAccount(encrypt string) *Account {
	// Determine the public key algorithm and parameters according to
	// the config file.
	// FIXME: better to decouple from config file by inputing as arguments.
	var pkAlgorithm keypair.KeyType
	var params interface{}
	var scheme s.SignatureScheme
	var err error
	if "" != encrypt {
		scheme, err = s.GetScheme(encrypt)
	} else {
		scheme, err = s.GetScheme(config.Parameters.SignatureScheme)
	}
	if err != nil {
		log.Warn("unknown signature scheme, use SHA256withECDSA as default.")
		scheme = s.SHA256withECDSA
	}
	switch scheme {
	case s.SHA224withECDSA, s.SHA3_224withECDSA:
		pkAlgorithm = keypair.PK_ECDSA
		params = keypair.P224
	case s.SHA256withECDSA, s.SHA3_256withECDSA, s.RIPEMD160withECDSA:
		pkAlgorithm = keypair.PK_ECDSA
		params = keypair.P256
	case s.SHA384withECDSA, s.SHA3_384withECDSA:
		pkAlgorithm = keypair.PK_ECDSA
		params = keypair.P384
	case s.SHA512withECDSA, s.SHA3_512withECDSA:
		pkAlgorithm = keypair.PK_ECDSA
		params = keypair.P521
	case s.SM3withSM2:
		pkAlgorithm = keypair.PK_SM2
		params = keypair.SM2P256V1
	case s.SHA512withEDDSA:
		pkAlgorithm = keypair.PK_EDDSA
		params = keypair.ED25519
	}

	pri, pub, _ := keypair.GenerateKeyPair(pkAlgorithm, params)
	address := types.AddressFromPubKey(pub)
	return &Account{
		PrivateKey: pri,
		PublicKey:  pub,
		Address:    address,
	}
}

func NewAccountWithPrivatekey(privateKey []byte) (*Account, error) {
	pri, err := keypair.DeserializePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	pub := pri.Public()
	address := types.AddressFromPubKey(pub)
	return &Account{
		PrivateKey: pri,
		PublicKey:  pub,
		Address:    address,
	}, nil
}

func (ac *Account) PrivKey() keypair.PrivateKey {
	return ac.PrivateKey
}

func (ac *Account) PubKey() keypair.PublicKey {
	return ac.PublicKey
}
