package main

import (
	//"fmt"
	"bytes"
	"encoding/hex"
	"github.com/dad-go/common"
	"github.com/dad-go/common/log"
	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/types"
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/eventbus/remote"
	"github.com/dad-go/eventbus/zmqremote"
	"time"
)

var (
	tx *types.Transaction
)

func init() {
	log.Init(log.Path, log.Stdout)

	bookKeepingPayload := &payload.BookKeeping{
		Nonce: uint64(time.Now().UnixNano()),
	}

	tx = &types.Transaction{
		Version:    0,
		Attributes: []*types.TxAttribute{},
		TxType:     types.BookKeeping,
		Payload:    bookKeepingPayload,
	}

	tempStr := "3369930accc1ddd067245e8edadcd9bea207ba5e1753ac18a51df77a343bfe92"
	hex, _ := hex.DecodeString(tempStr)
	var hash common.Uint256
	hash.Deserialize(bytes.NewReader(hex))
	tx.SetHash(hash)
}

func main() {
	var stopCh chan bool

	stopCh = make(chan bool)

	remote.Start("192.168.153.130:50011")
	server := actor.NewPID("192.168.153.130:50010", "Txn")
	/*props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *tp.CheckTxnRsp:
			fmt.Println(msg)
		}
	})*/

	//client := actor.Spawn(props)
	tmpBuffer := bytes.NewBuffer([]byte{})
	tx.Serialize(tmpBuffer)
	server.Tell(&zmqremote.MsgData{MsgType: 1, Data: tmpBuffer.Bytes()})
	<-stopCh
}
