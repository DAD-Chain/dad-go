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

package cmd

import (
	"io"
	"strconv"
	"strings"

	sdk "github.com/ontio/dad-go-go-sdk"
	"github.com/ontio/dad-go/cmd/utils"
	"github.com/ontio/dad-go/common/config"
	"github.com/urfave/cli"
)

// AppHelpTemplate is the test template for the default, global app help topic.
var AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .VisibleCommands}}

COMMANDS:{{range .Commands}}
   {{join .Names ", "}}{{ "\t" }}{{.Usage}}{{end}}{{end}}{{if .Version}}

VERSION:
   {{.Version}}{{end}}{{if .Copyright }}

COPYRIGHT:
   {{.Copyright}}{{end}}
`

var ontSdk *sdk.dad-goSdk

func localRpcAddress() string {
	return "http://localhost:" + strconv.Itoa(config.Parameters.HttpJsonPort)
}

func rpcAddress() string {
	return "http://localhost:" + strconv.Itoa(config.Parameters.HttpJsonPort)
}

func restfulAddr() string {
	return "http://localhost:" + strconv.Itoa(config.Parameters.HttpRestPort)
}

func init() {
	ontSdk = sdk.Newdad-goSdk()
	ontSdk.Rpc.SetAddress(rpcAddress())
	//cli.AppHelpTemplate = AppHelpTemplate
}

// flagGroup is a collection of flags belonging to a single topic.
type flagGroup struct {
	Name  string
	Flags []cli.Flag
}

// AppHelpFlagGroups is the application flags, grouped by functionality.
var AppHelpFlagGroups = []flagGroup{
	{
		Name: "dad-go INFO BLOCK",
		Flags: []cli.Flag{
			utils.HashInfoFlag,
			utils.HeightInfoFlag,
		},
	},

	{
		Name: "dad-go INFO TRANSACTION",
		Flags: []cli.Flag{
			utils.HashInfoFlag,
		},
	},

	{
		Name: "dad-go INFO VERSION",
		Flags: []cli.Flag{
			utils.NonOptionFlag,
		},
	},

	{
		Name: "dad-go INFO BLOCK HEIGHT",
		Flags: []cli.Flag{
			utils.NonOptionFlag,
		},
	},

	{
		Name: "dad-go ASSET TRANSFER",
		Flags: []cli.Flag{
			utils.TransactionFromFlag,
			utils.TransactionToFlag,
			utils.TransactionValueFlag,
			utils.ContractAddrFlag,
			utils.AccountPassFlag,
		},
	},

	{
		Name: "dad-go SET DEBUG",
		Flags: []cli.Flag{
			utils.DebugLevelFlag,
		},
	},

	{
		Name: "dad-go SET CONSENSUS",
		Flags: []cli.Flag{
			utils.ConsensusFlag,
		},
	},

	{
		Name: "dad-go CONTRACT DEPLOY",
		Flags: []cli.Flag{
			utils.ContractVmTypeFlag,
			utils.ContractStorageFlag,
			utils.ContractCodeFlag,
			utils.ContractNameFlag,
			utils.ContractVersionFlag,
			utils.ContractAuthorFlag,
			utils.ContractDescFlag,
			utils.ContractEmailFlag,
		},
	},

	{
		Name: "dad-go CONTRACT INVOKE",
		Flags: []cli.Flag{
			utils.ContractAddrFlag,
			utils.ContractParamsFlag,
		},
	},

	{
		Name: "MISC",
	},
}

type byCategory []flagGroup

func (a byCategory) Len() int      { return len(a) }
func (a byCategory) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byCategory) Less(i, j int) bool {
	iCat, jCat := a[i].Name, a[j].Name
	iIdx, jIdx := len(AppHelpFlagGroups), len(AppHelpFlagGroups) // ensure non categorized flags come last

	for i, group := range AppHelpFlagGroups {
		if iCat == group.Name {
			iIdx = i
		}
		if jCat == group.Name {
			jIdx = i
		}
	}

	return iIdx < jIdx
}

func HelpUsage() {
	cli.AppHelpTemplate = AppHelpTemplate

	// Define a one shot struct to pass to the usage template
	type helpData struct {
		App        interface{}
		FlagGroups []flagGroup
	}

	// Override the default app help printer, but only for the global app help
	originalHelpPrinter := cli.HelpPrinter
	cli.HelpPrinter = func(w io.Writer, tmpl string, data interface{}) {
		if tmpl == AppHelpTemplate {
			// Iterate over all the flags and add any uncategorized ones
			categorized := make(map[string]struct{})
			for _, group := range AppHelpFlagGroups {
				for _, flag := range group.Flags {
					categorized[flag.String()] = struct{}{}
				}
			}
			uncategorized := []cli.Flag{}
			for _, flag := range data.(*cli.App).Flags {
				if _, ok := categorized[flag.String()]; !ok {
					if strings.HasPrefix(flag.GetName(), "dashboard") {
						continue
					}
					uncategorized = append(uncategorized, flag)
				}
			}
			if len(uncategorized) > 0 {
				// Append all ungategorized options to the misc group
				miscs := len(AppHelpFlagGroups[len(AppHelpFlagGroups)-1].Flags)
				AppHelpFlagGroups[len(AppHelpFlagGroups)-1].Flags = append(AppHelpFlagGroups[len(AppHelpFlagGroups)-1].Flags, uncategorized...)

				// Make sure they are removed afterwards
				defer func() {
					AppHelpFlagGroups[len(AppHelpFlagGroups)-1].Flags = AppHelpFlagGroups[len(AppHelpFlagGroups)-1].Flags[:miscs]
				}()
			}
			// Render out custom usage screen
			originalHelpPrinter(w, tmpl, helpData{data, AppHelpFlagGroups})
		}
	}
}
