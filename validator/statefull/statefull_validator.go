package statefull

import (
	"github.com/dad-go/common/log"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/core/types"
	"github.com/dad-go/errors"
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/validator/db"
	vatypes "github.com/dad-go/validator/types"
)

type Validator interface {
	Register(poolId *actor.PID)
	UnRegister(poolId *actor.PID)
	VerifyType() vatypes.VerifyType
}

type validator struct {
	pid       *actor.PID
	id        string
	bestBlock db.BestBlock
}

func NewValidator(id string) (Validator, error) {

	validator := &validator{id: id}
	props := actor.FromProducer(func() actor.Actor {
		return validator
	})

	pid, err := actor.SpawnNamed(props, id)
	validator.pid = pid
	return validator, err
}

func (self *validator) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Info("Validator started and be ready to receive txn")
	case *actor.Stopping:
		log.Info("Validator stopping")
	case *actor.Restarting:
		log.Info("Validator Restarting")
	case *vatypes.CheckTx:
		log.Info("Validator receive tx")
		sender := context.Sender()
		height := ledger.DefLedger.GetCurrentBlockHeight()

		errCode := errors.ErrNoError
		hash := msg.Tx.Hash()

		exist, err := ledger.DefLedger.IsContainTransaction(&hash)
		if err != nil {
			log.Warn("query db error:", err)
			errCode = errors.ErrUnknown
		} else if exist {
			errCode = errors.ErrDuplicatedTx
		}

		response := &vatypes.CheckResponse{
			WorkerId: msg.WorkerId,
			Type:     self.VerifyType(),
			Hash:     msg.Tx.Hash(),
			Height:   height,
			ErrCode:  errCode,
		}

		sender.Tell(response)
	case *vatypes.UnRegisterAck:
		context.Self().Stop()
	case *types.Block:

		//bestBlock, _ := self.db.GetBestBlock()
		//if bestBlock.Height+1 < msg.Header.Height {
		//	// add sync block request
		//} else if bestBlock.Height+1 == msg.Header.Height {
		//	self.db.PersistBlock(msg)
		//}

	default:
		log.Info("Unknown msg type", msg)
	}

}

func (self *validator) VerifyType() vatypes.VerifyType {
	return vatypes.Statefull
}

func (self *validator) Register(poolId *actor.PID) {
	poolId.Tell(&vatypes.RegisterValidator{
		Sender: self.pid,
		Type:   self.VerifyType(),
		Id:     self.id,
	})
}

func (self *validator) UnRegister(poolId *actor.PID) {
	poolId.Tell(&vatypes.UnRegisterValidator{
		Id: self.id,
	})
}
