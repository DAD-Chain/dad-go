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
	"bytes"
	"encoding/json"
	"fmt"
	cmdCom "github.com/ontio/dad-go/cmd/common"
	"github.com/ontio/dad-go/cmd/utils"
	"github.com/urfave/cli"
)

var InfoCommand = cli.Command{
	Name:  "info",
	Usage: "Display informations about the chain",
	Subcommands: []cli.Command{
		{
			Action: utils.MigrateFlags(txInfo),
			Name:   "tx",
			Usage:  "Display transaction information",
			Flags: []cli.Flag{
				utils.TransactionHashFlag,
			},
			OnUsageError: cmdCom.CommonCommandErrorHandler,
			Description:  ``,
		},
		{
			Action: utils.MigrateFlags(blockInfo),
			Name:   "block",
			Usage:  "Display block information",
			Flags: []cli.Flag{
				utils.BlockHashInfoFlag,
				utils.BlockHeightInfoFlag,
			},
			OnUsageError: cmdCom.CommonCommandErrorHandler,
			Description:  ``,
		},
	},
	Description: ``,
}

func blockInfo(ctx *cli.Context) error {
	if !ctx.IsSet(utils.BlockHashInfoFlag.Name) && !ctx.IsSet(utils.BlockHeightInfoFlag.Name) {
		return fmt.Errorf("Missing hash or height argument")
	}
	var data []byte
	var err error
	if ctx.IsSet(utils.BlockHeightInfoFlag.Name) {
		data, err = utils.GetBlock(ctx.Uint(utils.BlockHeightInfoFlag.Name))
	} else {
		data, err = utils.GetBlock(ctx.String(utils.BlockHashInfoFlag.Name))
	}
	if err != nil {
		return err
	}
	var out bytes.Buffer
	err = json.Indent(&out, data, "", "   ")
	if err != nil {
		return err
	}
	fmt.Println(out.String())
	return nil
}

func txInfo(ctx *cli.Context) error {
	if !ctx.IsSet(utils.TransactionHashFlag.Name) {
		return fmt.Errorf("Missing hash argument")
	}
	txHash := ctx.String(utils.TransactionHashFlag.Name)
	txInfo, err := utils.GetRawTransaction(txHash)
	if err != nil {
		return err
	}
	var out bytes.Buffer
	err = json.Indent(&out, txInfo, "", "   ")
	if err != nil {
		return err
	}
	fmt.Println(out.String())
	return nil
}
