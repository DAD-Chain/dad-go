package client

import (
	ct "dad-go/core/contract"
	. "dad-go/common"
)

type IClientStore interface {
	BuildDatabase(path string)

	SaveStoredData(name string,value []byte)

	LoadStoredData(name string) []byte

	LoadAccount()  map[Uint160]*Account

	LoadContracts() map[Uint160]*ct.Contract
}
