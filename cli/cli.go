package cli

import (
	"math/rand"
	"time"

	"dad-go/common/config"
	"dad-go/common/log"
	"dad-go/crypto"
)

func init() {
	var path string = "./Log/"
	log.CreatePrintLog(path)
	crypto.SetAlg(config.Parameters.EncryptAlg)
	//seed transaction nonce
	rand.Seed(time.Now().UnixNano())
}
