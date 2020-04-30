package net

import (
	"github.com/dad-go/crypto"
	"github.com/dad-go/events"
	"github.com/dad-go/net/node"
	"github.com/dad-go/net/protocol"
	ns "github.com/dad-go/net/actor"
	"github.com/dad-go/eventbus/actor"
)

type Neter interface {
	//GetTxnPool(byCount bool) (map[Uint256]*types.Transaction, Fixed64)
	Xmit(interface{}) error
	GetEvent(eventName string) *events.Event
	GetBookKeepersAddrs() ([]*crypto.PubKey, uint64)
	//CleanTransactions(txns []*types.Transaction) error
	GetNeighborNoder() []protocol.Noder
	Tx(buf []byte)
	//AppendTxnPool(*types.Transaction) ErrCode
}

func SetTxnPoolPid(txnPid *actor.PID) {
	ns.SetTxnPoolPid(txnPid)
}

func SetConsensusPid(conPid *actor.PID) {
	ns.SetConsensusPid(conPid)
}

func SetLedgePid(conPid *actor.PID) {
	ns.SetLedgePid(conPid)
}

func InitNetServerActor(noder protocol.Noder) (*actor.PID, error){
	netServerPid, err := ns.InitNetServer(noder)
	return netServerPid, err
}

func StartProtocol(pubKey *crypto.PubKey) protocol.Noder{
	net := node.InitNode(pubKey)
	net.ConnectSeeds()
	return net
}
