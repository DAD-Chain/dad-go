package main

import (
	"os"
	"sort"

	_ "github.com/dad-go/cli"
	"github.com/dad-go/cli/bookkeeper"
	. "github.com/dad-go/cli/common"
	"github.com/dad-go/cli/data"
	"github.com/dad-go/cli/debug"
	"github.com/dad-go/cli/info"
	"github.com/dad-go/cli/privpayload"
	"github.com/dad-go/cli/test"
	"github.com/dad-go/cli/wallet"

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
