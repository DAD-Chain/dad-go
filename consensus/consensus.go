package consensus

import (
	"strings"

	"github.com/dad-go/account"
	"github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	"github.com/dad-go/consensus/dbft"
	"github.com/dad-go/consensus/solo"
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/net"
)

type ConsensusService interface {
	Start() error
	Halt() error
}

const (
	CONSENSUS_DBFT = "dbft"
	CONSENSUS_SOLO = "solo"
)

//func NewConsensusService(client cl.Client, localNet net.Neter) ConsensusService {
func NewConsensusService(account *account.Account, txpool *actor.PID, ledger *actor.PID, localNet net.Neter) (ConsensusService, error) {
	consensusType := strings.ToLower(config.Parameters.ConsensusType)
	if consensusType == "" {
		consensusType = CONSENSUS_DBFT
	}

	var consensus ConsensusService
	var err error
	switch consensusType {
	case CONSENSUS_DBFT:
		consensus = dbft.NewDbftService(account, "dbft", nil)
	case CONSENSUS_SOLO:
		consensus, err = solo.NewSoloService(account, nil, nil)
	}
	log.Infof("ConsensusType:%s", consensusType)
	return consensus, err
}
