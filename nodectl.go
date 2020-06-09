/*
 * Copyright (C) 2018 The dad-go Authors
 * This file is part of The dad-go library.
 *
 * The dad-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The dad-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"os"
	"sort"

	"github.com/urfave/cli"

	_ "github.com/ontio/dad-go/cli"
	"github.com/ontio/dad-go/cli/common"
	"github.com/ontio/dad-go/cli/test"
	"github.com/ontio/dad-go/cli/transfer"
	"github.com/ontio/dad-go/cli/wallet"
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
		common.NewIpFlag(),
		common.NewPortFlag(),
	}
	//commands
	app.Commands = []cli.Command{
		*test.NewCommand(),
		*wallet.NewCommand(),
		*transfer.NewCommand(),
	}
	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))

	app.Run(os.Args)
}
