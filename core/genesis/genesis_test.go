package genesis

import (
	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/ontio/dad-go/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenesisBlockInit(t *testing.T) {
	_, pub, _ := keypair.GenerateKeyPair(keypair.PK_ECDSA, keypair.P256)
	block, err := GenesisBlockInit([]keypair.PublicKey{pub})
	assert.Nil(t, err)
	assert.NotNil(t, block)
	assert.NotEqual(t, block.Header.TransactionsRoot, common.UINT256_EMPTY)
}
