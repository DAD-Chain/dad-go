package stateful

import (
	"bytes"
	"github.com/dad-go/common/log"
	"github.com/dad-go/core"
	"github.com/dad-go/core/genesis"
	tx "github.com/dad-go/core/types"
	"github.com/dad-go/crypto"
	"github.com/dad-go/eventbus/actor"
	tc "github.com/dad-go/txnpool/common"
	"github.com/dad-go/validator/db"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	crypto.SetAlg("")
	log.Init(log.Path, log.Stdout)
}

func TestVerifier(t *testing.T) {
	store, err := db.NewStore("temp.db")
	if assert.Nil(t, err) == false {
		return
	}

	verifier := NewDBVerifier(tc.StatefulV, store)

	props := actor.FromProducer(func() actor.Actor {
		return verifier
	})
	pid := actor.Spawn(props)
	verifier.SetPID(pid)

	_, issuer, _ := crypto.GenKeyPair()
	txn, _ := core.NewBookKeeperTransaction(&issuer, true, []byte{}, &issuer)

	block, _ := genesis.GenesisBlockInit([]*crypto.PubKey{&issuer})
	block.Transactions = append(block.Transactions, txn)
	block.RebuildMerkleRoot()

	buf := bytes.NewBuffer(nil)
	block.Serialize(buf)
	block.Deserialize(buf)

	pid.Tell(block)

	time.Sleep(time.Second * 1)

	req := &tc.VerifyReq{
		WorkerId: 1,
		Txn:      &tx.Transaction{},
	}

	for i := 0; i < 100; i++ {
		pid.Tell(req)
	}

	time.Sleep(time.Second * 2)
	verifier.Stop()
}
