package cli

import (
	"os"
	"sort"

	. "dad-go/cli/common"
	"dad-go/cli/consensus"
	"dad-go/cli/debug"
	"dad-go/cli/info"
	"dad-go/cli/test"
	"dad-go/cli/wallet"

	"github.com/urfave/cli"
)

func init() {
	app := cli.NewApp()
	app.Name = "nodectl"
	app.Version = "1.0.1"
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
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	app.Run(os.Args)
}
