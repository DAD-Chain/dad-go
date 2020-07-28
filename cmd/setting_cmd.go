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
	"encoding/json"
	"fmt"

	"github.com/ontio/dad-go/account"
	"github.com/ontio/dad-go/cmd/utils"
	jrpc "github.com/ontio/dad-go/http/base/rpc"
	"github.com/urfave/cli"
)

var (
	SettingCommand = cli.Command{
		Name:      "set",
		Action:    utils.MigrateFlags(settingCommand),
		Usage:     "Set node configurations",
		ArgsUsage: " ",
		Flags: []cli.Flag{
			utils.ConsensusFlag,
			utils.DebugLevelFlag,
		},
		Description: ``,
	}
)

func settingCommand(ctx *cli.Context) error {
	if ctx.IsSet(utils.DebugLevelFlag.Name) {
		level := ctx.Uint(utils.DebugLevelFlag.Name)
		resp, err := jrpc.Call(localRpcAddress(), "setdebuginfo", 0, []interface{}{level})
		if nil != err {
			return err
		}
		r := make(map[string]interface{})
		json.Unmarshal(resp, &r)
		fmt.Printf("%v\n", r)
		return nil
	} else if ctx.IsSet(utils.ConsensusFlag.Name) {
		consensusSwitch := ctx.String(utils.ConsensusFlag.Name)
		wallet := ctx.GlobalString("file")
		client := account.Open(wallet, nil)
		if client == nil {
			fmt.Println("Can't get local account.")
		}
		var resp []byte
		var err error
		fmt.Println("consensusSwitch:", consensusSwitch)
		switch consensusSwitch {
		case "on":
			resp, err = jrpc.Call(localRpcAddress(), "startconsensus", 0, []interface{}{1})
		case "off":
			resp, err = jrpc.Call(localRpcAddress(), "stopconsensus", 0, []interface{}{0})
		default:
			fmt.Println("Start:1; Stop:0; Pls enter valid value between 0 and 1.")
		}
		if nil != err {
			return err
		}
		r := make(map[string]interface{})
		json.Unmarshal(resp, &r)
		fmt.Printf("%v\n", r)
		return nil
	}
	cli.ShowSubcommandHelp(ctx)
	return nil
}
