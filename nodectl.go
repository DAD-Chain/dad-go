package main

import (
	"os"

	"github.com/DAD-Chain/dad-go/common/log"
	"github.com/DAD-Chain/dad-go/crypto"
	"github.com/DAD-Chain/dad-go/utility"
	"github.com/DAD-Chain/dad-go/utility/consensus"
	"github.com/DAD-Chain/dad-go/utility/info"
	"github.com/DAD-Chain/dad-go/utility/test"
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
