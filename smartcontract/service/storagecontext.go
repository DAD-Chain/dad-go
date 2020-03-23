package service

import (
	"github.com/dad-go/common"
	"github.com/dad-go/vm/neovm/interfaces"
)

type StorageContext struct {
	codeHash common.Uint160
}

func NewStorageContext(codeHash common.Uint160) *StorageContext {
	var storageContext StorageContext
	storageContext.codeHash = codeHash
	return &storageContext
}

func (sc *StorageContext) ToArray() []byte {
	return sc.codeHash.ToArray()
}

func (sc *StorageContext) Clone() interfaces.IInteropInterface {
	s := *sc
	return &s
}
