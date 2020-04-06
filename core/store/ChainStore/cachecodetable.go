package ChainStore

import (
	"github.com/dad-go/core/states"
	"github.com/dad-go/core/store"
	"github.com/dad-go/errors"

	"fmt"
)

type CacheCodeTable struct {
	store store.IStateStore
}

func NewCacheCodeTable(store store.IStateStore) *CacheCodeTable{
	return &CacheCodeTable{
		store: store,
	}
}

func (table *CacheCodeTable) GetCode(codeHash []byte) ([]byte, error) {
	value, _ := table.store.TryGet(store.ST_Contract, codeHash)
	if value == nil {
		return nil, errors.NewErr(fmt.Sprintf("[GetCode] TryGet contract error! codeHash:%x", codeHash))
	}

	return value.Value.(*states.ContractState).Code.Code, nil
}
