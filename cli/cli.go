package cli

import (
	"math/rand"
	"os"
	"sort"
	"time"

	"dad-go/cli/asset"
	. "dad-go/cli/common"
	"dad-go/cli/consensus"
	"dad-go/cli/debug"
	"dad-go/cli/info"
	"dad-go/cli/test"
	"dad-go/cli/wallet"
	"dad-go/common/log"
	"dad-go/crypto"

	"github.com/urfave/cli"
)

var Version string

func init() {
	var path string = "./Log/"
	log.CreatePrintLog(path)
	crypto.SetAlg(crypto.P256R1)
	//seed transaction nonce
	rand.Seed(time.Now().UnixNano())

	app := cli.NewApp()
	app.Name = "nodectl"
	app.Version = Version
	app.HelpName = "nodectl"
	app.Usage = "command line tool for dad-go blockchain"
	app.UsageText = "nodectl [global options] command [command options] [args]"
	app.HideHelp = false
	app.HideVersion = false
	//global options
	app.Flags = []cli.Flag{
		NewIpFlag(),
		NewPortFlag(),
	}
	//commands
	app.Commands = []cli.Command{
		*consensus.NewCommand(),
		*debug.NewCommand(),
		*info.NewCommand(),
		*test.NewCommand(),
		*wallet.NewCommand(),
		*asset.NewCommand(),
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	app.Run(os.Args)
}
