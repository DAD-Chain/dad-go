package net

import (
	"github.com/dad-go/crypto"
	"github.com/dad-go/events"
	"github.com/dad-go/net/node"
	"github.com/dad-go/net/protocol"
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

func StartProtocol(pubKey *crypto.PubKey) protocol.Noder {
	net := node.InitNode(pubKey)
	net.ConnectSeeds()

	return net
}
