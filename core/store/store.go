package store

import(
	. "dad-go/core/ledger"
	. "dad-go/core/store/LevelDBStore"
)

func NewLedgerStore() ILedgerStore {
	// TODO: read config file decide which db to use.
	ldbs,_ := NewLevelDBStore("Chain")

	return ldbs
}


