package main

import (
	"github.com/dad-go/account"
	"github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	"github.com/dad-go/consensus"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/crypto"
	"github.com/dad-go/net"
	"github.com/dad-go/http/jsonrpc"
	"github.com/dad-go/http/nodeinfo"
	"github.com/dad-go/http/restful"
	"github.com/dad-go/http/websocket"
	"github.com/dad-go/http/localrpc"
	"github.com/dad-go/net/protocol"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	DefaultMultiCoreNum = 4
)

func init() {
	log.Init(log.Path, log.Stdout)
	var coreNum int
	if config.Parameters.MultiCoreNum > DefaultMultiCoreNum {
		coreNum = int(config.Parameters.MultiCoreNum)
	} else {
		coreNum = DefaultMultiCoreNum
	}
	log.Debug("The Core number is ", coreNum)
	runtime.GOMAXPROCS(coreNum)
}

func main() {
	var acct *account.Account
	var err error
	var noder protocol.Noder
	log.Trace("Node version: ", config.Version)

	if len(config.Parameters.BookKeepers) < account.DefaultBookKeeperCount {
		log.Fatal("At least ", account.DefaultBookKeeperCount, " BookKeepers should be set at config.json")
		os.Exit(1)
	}
	crypto.SetAlg(config.Parameters.EncryptAlg)

	log.Info("0. Open the account")
	client := account.GetClient()
	if client == nil {
		log.Fatal("Can't get local account.")
		os.Exit(1)
	}
	acct, err = client.GetDefaultAccount()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Debug("The Node's PublicKey ", acct.PublicKey)
	defBookKeepers, err := client.GetBookKeepers()
	if err != nil {
		log.Fatalf("GetBookKeepers error:%s", err)
		os.Exit(1)
	}
	log.Info("1. Loading the Ledger")
	ledger.DefLedger, err = ledger.NewLedger()
	if err != nil {
		log.Fatalf("NewLedger error %s", err)
		os.Exit(1)
	}
	err = ledger.DefLedger.Init(defBookKeepers)
	if err != nil {
		log.Fatalf("DefLedger.Init error %s", err)
		os.Exit(1)
	}

	log.Info("3. Start the P2P networks")
	// Don't need two return value.
	noder = net.StartProtocol(acct.PublicKey)
	go restful.StartServer()
	//jsonrpc.RegistRpcNode(noder)

	noder.SyncNodeHeight()
	noder.WaitForPeersStart()
	noder.WaitForSyncBlkFinish()
	if protocol.SERVICENODENAME != config.Parameters.NodeType {
		log.Info("4. Start Consensus Services")
		consensusSrv := consensus.NewConsensusService(client, noder)
		//jsonrpc.RegistConsensusService(consensusSrv)
		go consensusSrv.Start()
		time.Sleep(5 * time.Second)
	}

	log.Info("--Start the RPC interface")
	go jsonrpc.StartRPCServer()
	go localrpc.StartLocalServer()
	go websocket.StartServer()
	if config.Parameters.HttpInfoStart {
		go nodeinfo.StartServer(noder)
	}

	go logCurrBlockHeight()

	//等待退出信号
	waitToExit()
}

func logCurrBlockHeight(){
	ticker := time.NewTicker(config.DEFAULTGENBLOCKTIME * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Trace("BlockHeight = ", ledger.DefLedger.GetCurrentBlockHeight())
			isNeedNewFile := log.CheckIfNeedNewFile()
			if isNeedNewFile {
				log.ClosePrintLog()
				log.Init(log.Path, os.Stdout)
			}
		}
	}
}

func waitToExit(){
	exit := make(chan bool, 0)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		for sig := range sc {
			log.Infof("dad-go received exit signal:%v.", sig.String())
			close(exit)
			break
		}
	}()
	<-exit
}
