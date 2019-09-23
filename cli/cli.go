package cli

import (
	"math/rand"
	"time"

	"dad-go/common/config"
	"dad-go/common/log"
	"dad-go/crypto"
)

func init() {
	log.Init()
	crypto.SetAlg(config.Parameters.EncryptAlg)
	//seed transaction nonce
	rand.Seed(time.Now().UnixNano())
}
