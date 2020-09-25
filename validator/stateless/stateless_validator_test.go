package stateless

import (
	"testing"

	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go-eventbus/actor"
	"github.com/ontio/dad-go/account"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/core/signature"
	ctypes "github.com/ontio/dad-go/core/types"
	"github.com/ontio/dad-go/core/utils"
	"github.com/ontio/dad-go/errors"
	"github.com/ontio/dad-go/smartcontract/types"
	types2 "github.com/ontio/dad-go/validator/types"
	"github.com/stretchr/testify/assert"
	"time"
)

func signTransaction(signer *account.Account, tx *ctypes.Transaction) error {
	hash := tx.Hash()
	sign, _ := signature.Sign(signer, hash[:])
	tx.Sigs = append(tx.Sigs, &ctypes.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sign},
	})
	return nil
}

func TestStatelessValidator(t *testing.T) {
	log.Init(log.PATH, log.Stdout)
	acc := account.NewAccount("")

	code := types.VmCode{
		VmType: types.NEOVM,
		Code:   []byte{1, 2, 3},
	}
	tx := utils.NewDeployTransaction(code, "test", "1", "author", "author@123.com", "test desp", false)

	tx.Payer = acc.Address

	signTransaction(acc, tx)

	validator := &validator{id: "test"}
	props := actor.FromProducer(func() actor.Actor {
		return validator
	})

	pid, err := actor.SpawnNamed(props, validator.id)
	assert.Nil(t, err)

	msg := &types2.CheckTx{WorkerId: 1, Tx: *tx}
	fut := pid.RequestFuture(msg, time.Second)

	res, err := fut.Result()
	assert.Nil(t, err)

	result := res.(*types2.CheckResponse)
	assert.Equal(t, result.ErrCode, errors.ErrNoError)
	assert.Equal(t, tx.Hash(), result.Hash)
}
