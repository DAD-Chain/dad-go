package main

import (
	"dad-go/utility"
	"dad-go/utility/consensus"
	"dad-go/utility/info"
	"dad-go/utility/test"
	"os"
)

func main() {
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
