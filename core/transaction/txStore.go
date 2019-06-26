package transaction

import (
. "dad-go/common"
)

// ILedgerStore provides func with store package.
type ILedgerStore interface {
	GetTransaction(hash Uint256) (*Transaction,error)

}
