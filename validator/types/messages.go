package types

import (
	"github.com/dad-go/common"
	"github.com/dad-go/core/types"
	"github.com/dad-go/errors"
	"github.com/dad-go/eventbus/actor"
)

// message
type RegisterValidator struct {
	Sender *actor.PID
	Type   VerifyType
	Id     string
}

type UnRegisterValidator struct {
	Id string
}

type UnRegisterAck struct {
	Id string
}

type CheckTx struct {
	Tx types.Transaction
}

type StatelessCheckResponse struct {
	ErrCode errors.ErrCode
	Hash    common.Uint256
}

type StatefullCheckResponse struct {
	ErrCode errors.ErrCode
	Hash    common.Uint256
	Height  int32
}

type VerifyType uint8

const (
	Stateless VerifyType = iota
	Statefull VerifyType = iota
)
