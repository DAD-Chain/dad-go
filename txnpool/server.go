package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/dad-go/common"
	"github.com/dad-go/common/log"
	"github.com/dad-go/core/payload"
	"github.com/dad-go/core/types"
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/eventbus/remote"
	tc "github.com/dad-go/txnpool/common"
	tp "github.com/dad-go/txnpool/proc"
	//"github.com/dad-go/validator/db"
	//"github.com/dad-go/validator/statefull"
	"github.com/dad-go/validator/stateless"
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

func startActor(obj interface{}) *actor.PID {
	props := actor.FromProducer(func() actor.Actor {
		return obj.(actor.Actor)
	})

	pid := actor.Spawn(props)
	if pid == nil {
		fmt.Println("Fail to start actor")
		return nil
	}
	return pid
}

func main() {
	remote.Start("192.168.153.130:50010")

	var s *tp.TXPoolServer
	var stopCh chan bool

	stopCh = make(chan bool)

	// Start txnpool server to receive msgs from p2p, consensus and valdiators
	s = tp.NewTxPoolServer(tc.MAXWORKERNUM)

	// Initialize an actor to handle the msgs from valdiators
	rspActor := tp.NewVerifyRspActor(s)
	rspPid := startActor(rspActor)
	if rspPid == nil {
		fmt.Println("Fail to start verify rsp actor")
		return
	}
	s.RegisterActor(tc.VerifyRspActor, rspPid)

	// Initialize an actor to handle the msgs from consensus
	tpa := tp.NewTxPoolActor(s)
	txPoolPid := startActor(tpa)
	if txPoolPid == nil {
		fmt.Println("Fail to start txnpool actor")
		return
	}
	s.RegisterActor(tc.TxPoolActor, txPoolPid)

	// Initialize an actor to handle the msgs from p2p and api
	ta := tp.NewTxActor(s)
	txPid := startActor(ta)
	if txPid == nil {
		fmt.Println("Fail to start txn actor")
		return
	}
	s.RegisterActor(tc.TxActor, txPid)

	// Start stateless validator
	statelessV, err := stateless.NewValidator("stateless")
	if err != nil {
		fmt.Println("failed to new stateless valdiator", err)
		return
	}
	statelessV.Register(rspPid)

	// Todo: depending on ledger db sync, when ledger db ready, enable it
	// Start stateful validator
	/*store, err := db.NewStore("temp.db")
		if err != nil {
			fmt.Println("failed to new store",err)
			return
		}

		statefulV, err := statefull.NewValidator("stateful", store)
		if err != nil {
			fmt.Println("failed to new stateful validator", err)
			return
		}
	    statefulV.Register(rspPid)*/

	<-stopCh
}
