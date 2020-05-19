package account

import (
	. "github.com/dad-go/common"
	ct "github.com/dad-go/core/contract"
)

type IClientStore interface {
	BuildDatabase(path string)

	SaveStoredData(name string, value []byte)

	LoadStoredData(name string) []byte

	LoadAccount() map[Address]*Account

	LoadContracts() map[Address]*ct.Contract
}
