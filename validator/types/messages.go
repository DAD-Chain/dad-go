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
	WorkerId uint8
	Tx       types.Transaction
}

type CheckResponse struct {
	WorkerId uint8
	Type     VerifyType
	Hash     common.Uint256
	Height   uint32
	ErrCode  errors.ErrCode
}

type VerifyType uint8

const (
	Stateless VerifyType = iota
	Statefull VerifyType = iota
)
