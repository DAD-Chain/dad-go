package main

import (
	"dad-go/common/log"
	"dad-go/crypto"
	"dad-go/utility"
	"dad-go/utility/consensus"
	"dad-go/utility/info"
	"dad-go/utility/test"
	"os"
)

const (
	path string = "./Log"
)

func main() {
	crypto.SetAlg(crypto.P256R1)
	log.CreatePrintLog(path)

	cmds := map[string]*utility.Command{
		"info":      info.Command,
		"consensus": consensus.Command,
		"test":      test.Command,
	}

	err := utility.Start(cmds)
	if err != nil {
		os.Exit(1)
	}
}
