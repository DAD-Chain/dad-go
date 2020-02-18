package cli

import (
	"math/rand"
	"time"

	"github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	"github.com/dad-go/crypto"
)

func init() {
	log.Init()
	crypto.SetAlg(config.Parameters.EncryptAlg)
	//seed transaction nonce
	rand.Seed(time.Now().UnixNano())
}
