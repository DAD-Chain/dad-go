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

package vbft

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go-crypto/vrf"
	"github.com/ontio/dad-go/account"
	"github.com/ontio/dad-go/common"
	"github.com/ontio/dad-go/common/config"
	"github.com/ontio/dad-go/consensus/vbft/config"
	"github.com/ontio/dad-go/core/ledger"
	"github.com/ontio/dad-go/core/signature"
	"github.com/ontio/dad-go/core/states"
	scommon "github.com/ontio/dad-go/core/store/common"
	gov "github.com/ontio/dad-go/smartcontract/service/native/governance"
	nutils "github.com/ontio/dad-go/smartcontract/service/native/utils"
)

func SignMsg(account *account.Account, msg ConsensusMsg) ([]byte, error) {

	data, err := msg.Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal msg when signing: %s", err)
	}

	return signature.Sign(account, data)
}

func hashData(data []byte) common.Uint256 {
	t := sha256.Sum256(data)
	f := sha256.Sum256(t[:])
	return common.Uint256(f)
}

func HashMsg(msg ConsensusMsg) (common.Uint256, error) {

	// FIXME: has to do marshal on each call

	data, err := SerializeVbftMsg(msg)
	if err != nil {
		return common.Uint256{}, fmt.Errorf("failed to marshal block: %s", err)
	}

	return hashData(data), nil
}

type seedData struct {
	BlockNum          uint32         `json:"block_num"`
	PrevBlockProposer uint32         `json:"prev_block_proposer"` // TODO: change to NodeID
	BlockRoot         common.Uint256 `json:"block_root"`
	VrfValue          []byte         `json:"vrf_value"`
}

func getParticipantSelectionSeed(block *Block) vconfig.VRFValue {

	data, err := json.Marshal(&seedData{
		BlockNum:          block.getBlockNum() + 1,
		PrevBlockProposer: block.getProposer(),
		BlockRoot:         block.Block.Header.BlockRoot,
		VrfValue:          block.getVrfValue(),
	})
	if err != nil {
		return vconfig.VRFValue{}
	}

	t := sha512.Sum512(data)
	f := sha512.Sum512(t[:])
	return vconfig.VRFValue(f)
}

type vrfData struct {
	BlockNum uint32 `json:"block_num"`
	PrevVrf  []byte `json:"prev_vrf"`
}

func computeVrf(sk keypair.PrivateKey, blkNum uint32, prevVrf []byte) ([]byte, []byte, error) {
	data, err := json.Marshal(&vrfData{
		BlockNum: blkNum,
		PrevVrf:  prevVrf,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("computeVrf failed to marshal vrfData: %s", err)
	}

	return vrf.Vrf(sk, data)
}

func verifyVrf(pk keypair.PublicKey, blkNum uint32, prevVrf, newVrf, proof []byte) error {
	data, err := json.Marshal(&vrfData{
		BlockNum: blkNum,
		PrevVrf:  prevVrf,
	})
	if err != nil {
		return fmt.Errorf("verifyVrf failed to marshal vrfData: %s", err)
	}

	result, err := vrf.Verify(pk, data, newVrf, proof)
	if err != nil {
		return fmt.Errorf("verifyVrf failed: %s", err)
	}
	if !result {
		return fmt.Errorf("verifyVrf failed")
	}
	return nil
}
func GetVbftConfigInfo() (*config.VBFTConfig, error) {
	//get governance view
	goveranceview, err := GetGovernanceView()
	if err != nil {
		return nil, err
	}

	//get preConfig
	storageKey := &states.StorageKey{
		ContractAddress: nutils.GovernanceContractAddress,
		Key:             append([]byte(gov.PRE_CONFIG)),
	}
	preCfg := new(gov.PreConfig)
	data, err := ledger.DefLedger.GetStorageItem(storageKey.ContractAddress, storageKey.Key)
	if err != nil && err != scommon.ErrNotFound {
		return nil, err
	}
	if data != nil {
		err = preCfg.Deserialize(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}
	}

	chainconfig := new(config.VBFTConfig)
	if preCfg.SetView == goveranceview.View {
		chainconfig = &config.VBFTConfig{
			N:                    uint32(preCfg.Configuration.N),
			C:                    uint32(preCfg.Configuration.C),
			K:                    uint32(preCfg.Configuration.K),
			L:                    uint32(preCfg.Configuration.L),
			BlockMsgDelay:        uint32(preCfg.Configuration.BlockMsgDelay),
			HashMsgDelay:         uint32(preCfg.Configuration.HashMsgDelay),
			PeerHandshakeTimeout: uint32(preCfg.Configuration.PeerHandshakeTimeout),
			MaxBlockChangeView:   uint32(preCfg.Configuration.MaxBlockChangeView),
		}
	} else {
		storageKey := &states.StorageKey{
			ContractAddress: nutils.GovernanceContractAddress,
			Key:             append([]byte(gov.VBFT_CONFIG)),
		}
		data, err := ledger.DefLedger.GetStorageItem(storageKey.ContractAddress, storageKey.Key)
		if err != nil {
			return nil, err
		}
		cfg := new(gov.Configuration)
		err = cfg.Deserialize(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}
		chainconfig = &config.VBFTConfig{
			N:                    uint32(cfg.N),
			C:                    uint32(cfg.C),
			K:                    uint32(cfg.K),
			L:                    uint32(cfg.L),
			BlockMsgDelay:        uint32(cfg.BlockMsgDelay),
			HashMsgDelay:         uint32(cfg.HashMsgDelay),
			PeerHandshakeTimeout: uint32(cfg.PeerHandshakeTimeout),
			MaxBlockChangeView:   uint32(cfg.MaxBlockChangeView),
		}
	}
	return chainconfig, nil
}

func GetPeersConfig() ([]*config.VBFTPeerStakeInfo, error) {
	goveranceview, err := GetGovernanceView()
	if err != nil {
		return nil, err
	}
	viewBytes, err := gov.GetUint32Bytes(goveranceview.View)
	if err != nil {
		return nil, err
	}
	storageKey := &states.StorageKey{
		ContractAddress: nutils.GovernanceContractAddress,
		Key:             append([]byte(gov.PEER_POOL), viewBytes...),
	}
	data, err := ledger.DefLedger.GetStorageItem(storageKey.ContractAddress, storageKey.Key)
	if err != nil {
		return nil, err
	}
	peerMap := &gov.PeerPoolMap{
		PeerPoolMap: make(map[string]*gov.PeerPoolItem),
	}
	err = peerMap.Deserialize(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	var peerstakes []*config.VBFTPeerStakeInfo
	for _, id := range peerMap.PeerPoolMap {
		if id.Status == gov.CandidateStatus || id.Status == gov.ConsensusStatus {
			config := &config.VBFTPeerStakeInfo{
				Index:      uint32(id.Index),
				PeerPubkey: id.PeerPubkey,
				InitPos:    id.InitPos + id.TotalPos,
			}
			peerstakes = append(peerstakes, config)
		}
	}
	return peerstakes, nil
}

func isUpdate(view uint32) (bool, error) {
	goveranceview, err := GetGovernanceView()
	if err != nil {
		return false, err
	}
	if goveranceview.View > view {
		return true, nil
	}
	return false, nil
}

func GetGovernanceView() (*gov.GovernanceView, error) {
	storageKey := &states.StorageKey{
		ContractAddress: nutils.GovernanceContractAddress,
		Key:             append([]byte(gov.GOVERNANCE_VIEW)),
	}
	data, err := ledger.DefLedger.GetStorageItem(storageKey.ContractAddress, storageKey.Key)
	if err != nil {
		return nil, err
	}
	governanceView := new(gov.GovernanceView)
	err = governanceView.Deserialize(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return governanceView, nil
}

func getChainConfig(blkNum uint32) (*vconfig.ChainConfig, error) {
	config, err := GetVbftConfigInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get chainconfig from leveldb: %s", err)
	}

	peersinfo, err := GetPeersConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get peersinfo from leveldb: %s", err)
	}
	goverview, err := GetGovernanceView()
	if err != nil {
		return nil, fmt.Errorf("failed to get governanceview failed:%s", err)
	}

	cfg, err := vconfig.GenesisChainConfig(config, peersinfo, goverview.TxHash, blkNum)
	if err != nil {
		return nil, fmt.Errorf("GenesisChainConfig failed: %s", err)
	}
	cfg.View = goverview.View
	return cfg, err
}
