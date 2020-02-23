package neovm

import (
	. "github.com/dad-go/vm/neovm/errors"
	"github.com/dad-go/common/log"
)

type IInteropService interface {
	Register(method string, handler func(*ExecutionEngine) (bool, error)) bool
	GetServiceMap() map[string]func(*ExecutionEngine) (bool, error)
}

type InteropService struct {
	serviceMap map[string]func(*ExecutionEngine) (bool, error)
}

func NewInteropService() *InteropService {
	var i InteropService
	i.serviceMap = make(map[string]func(*ExecutionEngine) (bool, error), 0)
	i.Register("System.ExecutionEngine.GetScriptContainer", i.GetCodeContainer)
	i.Register("System.ExecutionEngine.GetExecutingScriptHash", i.GetExecutingCodeHash)
	i.Register("System.ExecutionEngine.GetCallingScriptHash", i.GetCallingCodeHash)
	i.Register("System.ExecutionEngine.GetEntryScriptHash", i.GetEntryCodeHash)
	return &i
}

func (is *InteropService) Register(methodad-gome string, handler func(*ExecutionEngine) (bool, error)) bool {
	if _, ok := is.serviceMap[methodad-gome]; ok {
		return false
	}
	is.serviceMap[methodad-gome] = handler
	return true
}

func (i *InteropService) MergeMap(dictionary map[string]func(*ExecutionEngine) (bool, error)) {
	for k, v := range dictionary {
		if _, ok := i.serviceMap[k]; !ok {
			i.serviceMap[k] = v
		}
	}
}

func (i *InteropService) GetServiceMap() map[string]func(*ExecutionEngine) (bool, error) {
	return i.serviceMap
}

func (i *InteropService) Invoke(methodad-gome string, engine *ExecutionEngine) (bool, error) {
	if v, ok := i.serviceMap[methodad-gome]; ok {
		log.Error("Invoke Methodad-gome:", methodad-gome)
		return v(engine)
	}
	return false, ErrNotSupportService
}

func (i *InteropService) GetCodeContainer(engine *ExecutionEngine) (bool, error) {
	PushData(engine, engine.codeContainer)
	return true, nil
}

func (i *InteropService) GetExecutingCodeHash(engine *ExecutionEngine) (bool, error) {
	PushData(engine, engine.crypto.Hash160(engine.ExecutingCode()))
	return true, nil
}

func (i *InteropService) GetCallingCodeHash(engine *ExecutionEngine) (bool, error) {
	PushData(engine, engine.crypto.Hash160(engine.CallingContext().GetCodeHash()))
	return true, nil
}
func (i *InteropService) GetEntryCodeHash(engine *ExecutionEngine) (bool, error) {
	PushData(engine, engine.crypto.Hash160(engine.EntryContext().GetCodeHash()))
	return true, nil
}
