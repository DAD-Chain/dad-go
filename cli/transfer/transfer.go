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

package transfer

import (
	"github.com/urfave/cli"
	"fmt"
	"os"
	"github.com/dad-go/http/base/rpc"
	. "github.com/dad-go/cli/common"
	cutils "github.com/dad-go/core/utils"
	vmtypes "github.com/dad-go/vm/types"
	ctypes "github.com/dad-go/core/types"
	"github.com/dad-go/smartcontract/service/native/states"
	"github.com/dad-go/crypto"
	"math/big"
	"github.com/dad-go/common"
	"bytes"
	"github.com/dad-go/account"
	"encoding/json"
	"github.com/dad-go/common/serialization"
	"encoding/hex"
	"time"
)


func transferAction(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	// wallet name is wallet.dat by default
	contract := c.String("contract")
	if contract == "" {
		fmt.Println("Invalid contract address: ", contract)
		os.Exit(1)
	}
	ct, _ := common.HexToBytes(contract)
	ctu, _ := common.AddressParseFromBytes(ct)
	from := c.String("from")
	if from == "" {
		fmt.Println("Invalid sender address: ", from)
		os.Exit(1)
	}
	f, _ := common.HexToBytes(from)
	fu, _ := common.AddressParseFromBytes(f)
	to := c.String("to")
	if to == "" {
		fmt.Println("Invalid revicer address: ", to)
		os.Exit(1)
	}
	t, _ := common.HexToBytes(to)
	tu, _ := common.AddressParseFromBytes(t)
	value := c.Int64("value")
	if value <= 0 {
		fmt.Println("Invalid ont amount: ", value)
		os.Exit(1)
	}

	var sts []*states.State
	sts = append(sts, &states.State{
		From: fu,
		To: tu,
		Value: big.NewInt(value),
	})
	transfers := new(states.Transfers)
	fmt.Println("ctu:", ctu)
	transfers.Params = append(transfers.Params, &states.TokenTransfer{
		Contract: ctu,
		States: sts,
	})

	bf := new(bytes.Buffer)
	if err := serialization.WriteVarBytes(bf, []byte("Token.Common.Transfer")); err != nil {
		fmt.Println("Serialize transfer falg error.")
		os.Exit(1)
	}
	if err := transfers.Serialize(bf); err != nil {
		fmt.Println("Serialize transfers struct error.")
		os.Exit(1)
	}

	tx := cutils.NewInvokeTransaction(vmtypes.VmCode{
		VmType: vmtypes.Native,
		Code: bf.Bytes(),
	})

	tx.Nonce = uint32(time.Now().Unix())

	acct := account.Open(account.WalletFileName, []byte("passwordtest"))
	acc, err := acct.GetDefaultAccount(); if err != nil {
		fmt.Println("GetDefaultAccount error:", err)
		os.Exit(1)
	}
	if err := signTransaction(acc, tx); err != nil {
		fmt.Println("signTransaction error:", err)
		os.Exit(1)
	}

	txbf := new(bytes.Buffer)
	if err := tx.Serialize(txbf); err != nil {
		fmt.Println("Serialize transaction error.")
		os.Exit(1)
	}
	resp, err := rpc.Call(RpcAddress(), "sendrawtransaction", 0,
		[]interface{}{hex.EncodeToString(txbf.Bytes())})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	r := make(map[string]interface{})
	err = json.Unmarshal(resp, &r)
	if err != nil {
		fmt.Println("Unmarshal JSON failed")
		os.Exit(1)
	}
	switch r["result"].(type) {
	case map[string]interface{}:

	case string:
		fmt.Println(r["result"].(string))
		os.Exit(1)
	}
	return nil
}

func signTransaction(signer *account.Account, tx *ctypes.Transaction) error {
	hash := tx.Hash()
	sign, _ := crypto.Sign(signer.PrivateKey, hash[:])
	tx.Sigs = append(tx.Sigs, &ctypes.Sig{
		PubKeys: []*crypto.PubKey{signer.PublicKey},
		M: 1,
		SigData: [][]byte{sign},
	})
	return nil
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "transfer",
		Usage:       "user ont transfer",
		Description: "With nodectl transfer, you could transfer ont.",
		ArgsUsage:   "[args]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "contract, c",
				Usage: "contract address",
			},
			cli.StringFlag{
				Name:  "from, f",
				Usage: "sender address",
			},
			cli.StringFlag{
				Name:  "to, t",
				Usage: "revicer address",
			},
			cli.Int64Flag{
				Name:  "value, v",
				Usage: "ont amount",
			},
		},
		Action: transferAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			return cli.NewExitError("", 1)
		},
	}
}

