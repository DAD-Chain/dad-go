package httpjsonrpc

/*
import (
	. "github.com/dad-go/account"
	. "github.com/dad-go/common"
	"github.com/dad-go/common/log"
	. "github.com/dad-go/core/asset"
	"github.com/dad-go/core/contract"
	"github.com/dad-go/core/signature"
	"github.com/dad-go/core/types"
	"strconv"
)

const (
	ASSETPREFIX = "dad-go"
)

func SignTx(admin *Account, tx *types.Transaction) {
	signdate, err := signature.SignBySigner(tx, admin)
	if err != nil {
		log.Error(err, "signdate SignBySigner failed")
	}
	transactionContract, _ := contract.CreateSignatureContract(admin.PublicKey)
	transactionContractContext := contract.NewContractContext(tx)
	transactionContractContext.AddContract(transactionContract, admin.PublicKey, signdate)
	tx.SetPrograms(transactionContractContext.GetPrograms())
}
*/
