package client

import (
	. "dad-go/common"
	ct "dad-go/core/contract"
)

type IClientStore interface {
	BuildDatabase(path string)

	SaveStoredData(name string, value []byte)

	LoadStoredData(name string) []byte

	LoadAccount() map[Uint160]*Account

	LoadContracts() map[Uint160]*ct.Contract
}
