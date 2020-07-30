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

package consensus

import (
	"strings"

	"github.com/ontio/dad-go/account"
	"github.com/ontio/dad-go/common/config"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/consensus/dbft"
	"github.com/ontio/dad-go/consensus/solo"
	"github.com/ontio/dad-go/eventbus/actor"
	"github.com/ontio/dad-go/consensus/vbft"
)

type ConsensusService interface {
	Start() error
	Halt() error
	GetPID() *actor.PID
}

const (
	CONSENSUS_DBFT = "dbft"
	CONSENSUS_SOLO = "solo"
	CONSENSUS_VBFT = "vbft"
)

func NewConsensusService(account *account.Account, txpool *actor.PID, ledger *actor.PID, p2p *actor.PID) (ConsensusService, error) {
	consensusType := strings.ToLower(config.Parameters.ConsensusType)
	if consensusType == "" {
		consensusType = CONSENSUS_DBFT
	}

	var consensus ConsensusService
	var err error
	switch consensusType {
	case CONSENSUS_DBFT:
		consensus, err = dbft.NewDbftService(account, txpool, p2p)
	case CONSENSUS_SOLO:
		consensus, err = solo.NewSoloService(account, txpool, ledger)
	case CONSENSUS_VBFT:
		consensus, err = vbft.NewVbftServer(account, txpool, ledger, p2p)
	}
	log.Infof("ConsensusType:%s", consensusType)
	return consensus, err
}
