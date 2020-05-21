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
	"errors"
	"time"
	"github.com/dad-go/common"
	"github.com/dad-go/common/config"
	"github.com/dad-go/core/types"
	"github.com/dad-go/core/utils"
	"github.com/dad-go/crypto"
	vmtypes "github.com/dad-go/vm/types"
)

const (
	BlockVersion      uint32 = 0
	GenesisNonce      uint64 = 2083236893
	DecrementInterval uint32 = 2000000

	OntRegisterAmount = 1000000000
	OngRegisterAmount = 1000000000
)

var (
	GenerationAmount = [17]uint32{80, 70, 60, 50, 40, 30, 20, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}

	OntContractAddress, _ = common.Uint160ParseFromBytes([]byte{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	OngContractAddress, _ = common.Uint160ParseFromBytes([]byte{0xff, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2})

	ONTToken   = NewGoverningToken()
	ONGToken   = NewUtilityToken()
	ONTTokenID = ONTToken.Hash()
	ONGTokenID = ONGToken.Hash()
)

var GenBlockTime = (config.DEFAULTGENBLOCKTIME * time.Second)

var GenesisBookKeepers []*crypto.PubKey

func GenesisBlockInit(defaultBookKeeper []*crypto.PubKey) (*types.Block, error) {
	//getBookKeeper
	GenesisBookKeepers = defaultBookKeeper
	nextBookKeeper, err := types.AddressFromBookKeepers(defaultBookKeeper)
	if err != nil {
		return nil, errors.New("[Block],GenesisBlockInit err with GetBookKeeperAddress")
	}
	//blockdata
	genesisHeader := &types.Header{
		Version:          BlockVersion,
		PrevBlockHash:    common.Uint256{},
		TransactionsRoot: common.Uint256{},
		Timestamp:        uint32(uint32(time.Date(2017, time.February, 23, 0, 0, 0, 0, time.UTC).Unix())),
		Height:           uint32(0),
		ConsensusData:    GenesisNonce,
		NextBookKeeper:   nextBookKeeper,

		BookKeepers: nil,
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
		},
	}
	return genesisBlock, nil
}

func NewGoverningToken() *types.Transaction {
	tx := utils.NewDeployTransaction([]byte("ONT Token"), "ONT", "1.0",
		"dad-go Team", "contact@ont.io", "dad-go Network ONT Token", vmtypes.NativeVM, true)
	return tx
}

func NewUtilityToken() *types.Transaction {
	tx := utils.NewDeployTransaction([]byte("ONT Token"), "ONG", "1.0",
		"dad-go Team", "contact@ont.io", "dad-go Network ONG Token", vmtypes.NativeVM, true)
	return tx
}

func NewGoverningInit() *types.Transaction {
	vmCode := vmtypes.VmCode{
		VmType: vmtypes.NativeVM,
		Code: []byte{14, 84, 111, 107, 101, 110, 46, 79, 110, 116, 46, 73, 110, 105, 116},
	}
	tx := utils.NewInvokeTransaction(vmCode)
	return tx
}
