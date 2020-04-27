package statefull

import (
	"github.com/dad-go/common/log"
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
	db        db.TransactionProvider
	bestBlock db.BestBlock
}

func NewValidator(id string, db db.TransactionProvider) (Validator, error) {
	bestBlock, err := db.GetBestBlock()
	if err != nil {
		return nil, err
	}

	validator := &validator{id: id, db: db, bestBlock: bestBlock}
	props := actor.FromProducer(func() actor.Actor {
		return validator
	})

	validator.pid, err = actor.SpawnNamed(props, id)
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
		bestBlock, _ := self.db.GetBestBlock()

		errCode := errors.ErrNoError
		if exist := self.db.ContainTransaction(msg.Tx.Hash()); exist {
			errCode = errors.ErrDuplicatedTx
		}

		response := &vatypes.CheckResponse{
			WorkerId: msg.WorkerId,
			Type:     self.VerifyType(),
			Hash:     msg.Tx.Hash(),
			Height:   bestBlock.Height,
			ErrCode:  errCode,
		}

		sender.Tell(response)
	case *vatypes.UnRegisterAck:
		context.Self().Stop()
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
