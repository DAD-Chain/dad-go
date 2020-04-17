package core

import (
	"bytes"

	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/types"

	"testing"
)

func TestNewONTTransferTransaction(t *testing.T) {
	executor := transactionExecutor{}
}

type TransactionExector interface {
	Execute(tx *types.Transaction)
}

type transactionExecutor struct {
	db map[types.Address]int64
}

func (self *transactionExecutor) Execute(tx *types.Transaction) {
	switch pld := tx.Payload.(type) {
	case *payload.InvokeCode:
		vmcode := pld.Code
		if vmcode.CodeType == types.NativeVM {
			if bytes.Equal(vmcode.Code, []byte("ont")) {
				if bytes.Equal(pld.Params[:8], []byte("transfer")) {
					var from, to types.Address
					copy(from[:], pld.Params[8:28])
					copy(to[:], pld.Params[28:48])
					value := int64(10)
					if self.db[from] >= value {
						self.db[from] -= value
						self.db[to] += value
					}
				}
			}
		}
	}
}

func NewONTTransferTransaction(from, to types.Address) *types.Transaction {
	code := []byte("ont")
	params := append([]byte("transfer"), from[:]...)
	params = append(params, to[:]...)
	vmcode := types.VmCode{
		CodeType: types.NativeVM,
		Code:     code,
	}

	tx, _ := NewInvokeTransaction(vmcode, params)
	return tx
}

func TestONT(t *testing.T) {

}
