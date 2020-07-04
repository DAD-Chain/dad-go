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

package genesis

import (
	"bytes"
	"errors"
	"time"

	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/config"
	"github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/core/utils"
	"github.com/ontio/dad-go/smartcontract/states"
	stypes "github.com/ontio/dad-go/smartcontract/types"
	"github.com/ontio/dad-go-crypto/keypair"
)

const (
	BlockVersion uint32 = 0
	GenesisNonce uint64 = 2083236893
)

var (
	OntContractAddress, _ = common.AddressParseFromBytes([]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01})
	OngContractAddress, _ = common.AddressParseFromBytes([]byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02})

	ONTToken   = NewGoverningToken()
	ONGToken   = NewUtilityToken()
	ONTTokenID = ONTToken.Hash()
	ONGTokenID = ONGToken.Hash()
)

var GenBlockTime = (config.DEFAULT_GEN_BLOCK_TIME * time.Second)

var GenesisBookkeepers []keypair.PublicKey

func GenesisBlockInit(defaultBookkeeper []keypair.PublicKey) (*types.Block, error) {
	//getBookkeeper
	GenesisBookkeepers = defaultBookkeeper
	nextBookkeeper, err := types.AddressFromBookkeepers(defaultBookkeeper)
	if err != nil {
		return nil, errors.New("[Block],GenesisBlockInit err with GetBookkeeperAddress")
	}
	//blockdata
	genesisHeader := &types.Header{
		Version:          BlockVersion,
		PrevBlockHash:    common.Uint256{},
		TransactionsRoot: common.Uint256{},
		Timestamp:        uint32(uint32(time.Date(2017, time.February, 23, 0, 0, 0, 0, time.UTC).Unix())),
		Height:           uint32(0),
		ConsensusData:    GenesisNonce,
		NextBookkeeper:   nextBookkeeper,

		Bookkeepers: nil,
		SigData:     nil,
	}

	//block
	ont := NewGoverningToken()
	ong := NewUtilityToken()

	genesisBlock := &types.Block{
		Header: genesisHeader,
		Transactions: []*types.Transaction{
			ont,
			ong,
			NewGoverningInit(),
			NewUtilityInit(),
		},
	}
	genesisBlock.RebuildMerkleRoot()
	return genesisBlock, nil
}

func NewGoverningToken() *types.Transaction {
	tx := utils.NewDeployTransaction(stypes.VmCode{Code: OntContractAddress[:], VmType: stypes.Native}, "ONT", "1.0",
		"dad-go Team", "contact@ont.io", "dad-go Network ONT Token", true)
	return tx
}

func NewUtilityToken() *types.Transaction {
	tx := utils.NewDeployTransaction(stypes.VmCode{Code: OngContractAddress[:], VmType: stypes.Native}, "ONG", "1.0",
		"dad-go Team", "contact@ont.io", "dad-go Network ONG Token", true)
	return tx
}

func NewGoverningInit() *types.Transaction {
	init := states.Contract{
		Address: OntContractAddress,
		Method:  "init",
	}
	bf := new(bytes.Buffer)
	init.Serialize(bf)
	vmCode := stypes.VmCode{
		VmType: stypes.Native,
		Code:   bf.Bytes(),
	}
	tx := utils.NewInvokeTransaction(vmCode)
	return tx
}

func NewUtilityInit() *types.Transaction {
	init := states.Contract{
		Address: OngContractAddress,
		Method:  "init",
	}
	bf := new(bytes.Buffer)
	init.Serialize(bf)
	vmCode := stypes.VmCode{
		VmType: stypes.Native,
		Code:   bf.Bytes(),
	}
	tx := utils.NewInvokeTransaction(vmCode)
	return tx
}
