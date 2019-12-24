package main

import (
	"os"
	"sort"

	_ "dad-go/cli"
	"dad-go/cli/asset"
	"dad-go/cli/bookkeeper"
	. "dad-go/cli/common"
	"dad-go/cli/data"
	"dad-go/cli/debug"
	"dad-go/cli/info"
	"dad-go/cli/privpayload"
	"dad-go/cli/test"
	"dad-go/cli/wallet"

	"github.com/urfave/cli"
)

var Version string

func main() {
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
		*debug.NewCommand(),
		*info.NewCommand(),
		*test.NewCommand(),
		*wallet.NewCommand(),
		*asset.NewCommand(),
		*privpayload.NewCommand(),
		*data.NewCommand(),
		*bookkeeper.NewCommand(),
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	app.Run(os.Args)
}
