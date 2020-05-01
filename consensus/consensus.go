package consensus

import (
	"strings"

	"github.com/dad-go/account"
	"github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	"github.com/dad-go/consensus/dbft"
	"github.com/dad-go/consensus/solo"
	"github.com/dad-go/eventbus/actor"
)

type ConsensusService interface {
	Start() error
	Halt() error
	GetPID() *actor.PID
}

const (
	CONSENSUS_DBFT = "dbft"
	CONSENSUS_SOLO = "solo"
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
	}
	log.Infof("ConsensusType:%s", consensusType)
	return consensus, err
}
