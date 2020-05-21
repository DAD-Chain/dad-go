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

package info

//import (
//	"fmt"
//	"os"
//
//	. "github.com/dad-go/cli/common"
//	"github.com/dad-go/http/httpjsonrpc"
//
//	"github.com/urfave/cli"
//)
//
//func infoAction(c *cli.Context) (err error) {
//	if c.NumFlags() == 0 {
//		cli.ShowSubcommandHelp(c)
//		return nil
//	}
//	blockhash := c.String("blockhash")
//	txhash := c.String("txhash")
//	bestblockhash := c.Bool("bestblockhash")
//	height := c.Int("height")
//	blockcount := c.Bool("blockcount")
//	connections := c.Bool("connections")
//	neighbor := c.Bool("neighbor")
//	state := c.Bool("state")
//	version := c.Bool("nodeversion")
//
//	var resp []byte
//	var output [][]byte
//	if height != -1 {
//		resp, err = jsonrpc.Call(Address(), "getblock", 0, []interface{}{height})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//	}
//
//	if c.String("blockhash") != "" {
//		resp, err = jsonrpc.Call(Address(), "getblock", 0, []interface{}{blockhash})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//	}
//
//	if bestblockhash {
//		resp, err = jsonrpc.Call(Address(), "getbestblockhash", 0, []interface{}{})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//	}
//
//	if blockcount {
//		resp, err = jsonrpc.Call(Address(), "getblockcount", 0, []interface{}{})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//	}
//
//	if connections {
//		resp, err = jsonrpc.Call(Address(), "getconnectioncount", 0, []interface{}{})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//	}
//
//	if neighbor {
//		resp, err := jsonrpc.Call(Address(), "getneighbor", 0, []interface{}{})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//	}
//
//	if state {
//		resp, err := jsonrpc.Call(Address(), "getnodestate", 0, []interface{}{})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//	}
//
//	if txhash != "" {
//		resp, err = jsonrpc.Call(Address(), "getrawtransaction", 0, []interface{}{txhash})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//	}
//
//	if version {
//		resp, err = jsonrpc.Call(Address(), "getversion", 0, []interface{}{})
//		if err != nil {
//			fmt.Fprintln(os.Stderr, err)
//			return err
//		}
//		output = append(output, resp)
//
//	}
//	for _, v := range output {
//		FormatOutput(v)
//	}
//
//	return nil
//}
//
//func NewCommand() *cli.Command {
//	return &cli.Command{
//		Name:        "info",
//		Usage:       "show blockchain information",
//		Description: "With nodectl info, you could look up blocks, transactions, etc.",
//		ArgsUsage:   "[args]",
//		Flags: []cli.Flag{
//			cli.StringFlag{
//				Name:  "blockhash, b",
//				Usage: "hash for querying a block",
//			},
//			cli.StringFlag{
//				Name:  "txhash, t",
//				Usage: "hash for querying a transaction",
//			},
//			cli.BoolFlag{
//				Name:  "bestblockhash",
//				Usage: "latest block hash",
//			},
//			cli.IntFlag{
//				Name:  "height",
//				Usage: "block height for querying a block",
//				Value: -1,
//			},
//			cli.BoolFlag{
//				Name:  "blockcount, c",
//				Usage: "block number in blockchain",
//			},
//			cli.BoolFlag{
//				Name:  "connections",
//				Usage: "connection count",
//			},
//			cli.BoolFlag{
//				Name:  "neighbor",
//				Usage: "neighbor information of current node",
//			},
//			cli.BoolFlag{
//				Name:  "state, s",
//				Usage: "current node state",
//			},
//			cli.BoolFlag{
//				Name:  "nodeversion, v",
//				Usage: "version of connected remote node",
//			},
//		},
//		Action: infoAction,
//		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
//			PrintError(c, err, "info")
//			return cli.NewExitError("", 1)
//		},
//	}
//}
