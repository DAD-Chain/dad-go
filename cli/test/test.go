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

package test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/dad-go/account"
	. "github.com/dad-go/cli/common"
	"github.com/dad-go/common"
	"github.com/dad-go/core/genesis"
	"github.com/dad-go/core/signature"
	"github.com/dad-go/core/types"
	"github.com/dad-go/core/utils"
	"github.com/dad-go/http/base/rpc"
	"github.com/dad-go/smartcontract/service/native/states"
	vmtypes "github.com/dad-go/vm/types"
	"github.com/ontio/dad-go-crypto/keypair"
	"github.com/urfave/cli"
)

func signTransaction(signer *account.Account, tx *types.Transaction) error {
	hash := tx.Hash()
	sign, _ := signature.Sign(signer.PrivateKey, hash[:])
	tx.Sigs = append(tx.Sigs, &types.Sig{
		PubKeys: []keypair.PublicKey{signer.PublicKey},
		M:       1,
		SigData: [][]byte{sign},
	})
	return nil
}

func testAction(c *cli.Context) (err error) {
	txnNum := c.Int("num")
	passwd := c.String("password")

	acct := account.Open(account.WalletFileName, []byte(passwd))
	acc, err := acct.GetDefaultAccount()
	if err != nil {
		fmt.Println("GetDefaultAccount error:", err)
		os.Exit(1)
	}

	transferTest(txnNum, acc)

	return nil
}

func transferTest(n int, acc *account.Account) {
	if n <= 0 {
		n = 1
	}

	for i := 0; i < n; i++ {
		tx := NewOntTransferTransaction(acc.Address, acc.Address, int64(i))
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
			return
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
		}
	}
}

func NewOntTransferTransaction(from, to common.Address, value int64) *types.Transaction {
	var sts []*states.State
	sts = append(sts, &states.State{
		From:  from,
		To:    to,
		Value: big.NewInt(value),
	})
	transfers := new(states.Transfers)
	transfers.States = sts

	bf := new(bytes.Buffer)

	if err := transfers.Serialize(bf); err != nil {
		fmt.Println("Serialize transfers struct error.")
		os.Exit(1)
	}

	cont := &states.Contract{
		Address: genesis.OntContractAddress,
		Method:  "transfer",
		Args:    bf.Bytes(),
	}

	ff := new(bytes.Buffer)
	if err := cont.Serialize(ff); err != nil {
		fmt.Println("Serialize contract struct error.")
		os.Exit(1)
	}

	tx := utils.NewInvokeTransaction(vmtypes.VmCode{
		VmType: vmtypes.Native,
		Code:   ff.Bytes(),
	})

	tx.Nonce = uint32(time.Now().Unix())

	return tx
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "test",
		Usage:       "run test routine",
		Description: "With nodectl test, you could run simple tests.",
		ArgsUsage:   "[args]",
		Flags: []cli.Flag{
			cli.IntFlag{
				Name:  "num, n",
				Usage: "sample transaction numbers",
				Value: 1,
			},
			cli.StringFlag{
				Name:  "password, p",
				Usage: "wallet password",
				Value: "passwordtest",
			},
		},
		Action: testAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			PrintError(c, err, "test")
			return cli.NewExitError("", 1)
		},
	}
}
