package transaction

import (
. "github.com/DAD-Chain/dad-go/common"
)

// ILedgerStore provides func with store package.
type ILedgerStore interface {
	GetTransaction(hash Uint256) (*Transaction,error)

}
