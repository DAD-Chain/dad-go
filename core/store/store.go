package store

import(
	. "github.com/DAD-Chain/dad-go/core/ledger"
	. "github.com/DAD-Chain/dad-go/core/store/LevelDBStore"
)

func NewLedgerStore() ILedgerStore {
	// TODO: read config file decide which db to use.
	ldbs,_ := NewLevelDBStore("Chain")

	return ldbs
}


