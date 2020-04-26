package actor

import (
	"time"
	"github.com/dad-go/core/types"
	"github.com/dad-go/eventbus/actor"
	. "github.com/dad-go/errors"
	. "github.com/dad-go/common"
)

var txPid *actor.PID

func SetTxActor(actr *actor.PID) {
	txPid = actr
}

func AppendTxnPool(txn *types.Transaction) ErrCode {
	future := txPid.RequestFuture(txn, 10*time.Second)
	result, err := future.Result()
	if err != nil {
		return ErrUnknown
	}
	if errCode, ok := result.(ErrCode); !ok {
		return errCode
	} else {
		return ErrUnknown
	}
}

func GetTxnPool(byCount bool) (map[Uint256]*types.Transaction, Fixed64) {
	future := txPid.RequestFuture(byCount, 10*time.Second)
	_, err := future.Result()
	if err != nil {
		return nil, 0
	}
	return nil, 0
}
