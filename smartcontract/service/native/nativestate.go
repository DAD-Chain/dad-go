package native

import (
	"github.com/dad-go/smartcontract/storage"
	scommon "github.com/dad-go/core/store/common"
	"bytes"
	"github.com/dad-go/errors"
	"github.com/dad-go/common/serialization"
	"github.com/dad-go/core/types"
	"github.com/dad-go/smartcontract/event"
)

type (
	Handler func(native *NativeService) (bool, error)
)

type NativeService struct {
	CloneCache *storage.CloneCache
	ServiceMap  map[string]Handler
	Notifications []*event.NotifyEventInfo
	Input []byte
	Tx *types.Transaction
}

func NewNativeService(dbCache scommon.IStateStore, input []byte, tx *types.Transaction) *NativeService {
	var nativeService NativeService
	nativeService.CloneCache = storage.NewCloneCache(dbCache)
	nativeService.Input = input
	nativeService.Tx = tx
	nativeService.ServiceMap = make(map[string]Handler)
	nativeService.Register("Token.Common.Transfer", Transfer)
	nativeService.Register("Token.Ont.Init", OntInit)
	nativeService.Register("Token.Ong.Init", OngInit)
	return &nativeService
}

func(native *NativeService) Register(methodad-gome string, handler Handler) {
	native.ServiceMap[methodad-gome] = handler
}

func(native *NativeService) Invoke() (bool, error){
	bf := bytes.NewBuffer(native.Input)
	serviceName, err := serialization.ReadVarBytes(bf); if err != nil {
		return false, err
	}
	service, ok := native.ServiceMap[string(serviceName)]; if !ok {
		return false, errors.NewErr("Native does not support this service!")
	}
	native.Input = bf.Bytes()
	return service(native)
}








