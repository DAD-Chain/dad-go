package net

import (
	"github.com/DAD-Chain/dad-go/common"
	"github.com/DAD-Chain/dad-go/config"
	"github.com/DAD-Chain/dad-go/core/transaction"
	"github.com/DAD-Chain/dad-go/crypto"
	"github.com/DAD-Chain/dad-go/events"
	"github.com/DAD-Chain/dad-go/net/node"
	"github.com/DAD-Chain/dad-go/net/protocol"
)

type Neter interface {
	GetTxnPool(cleanPool bool) map[common.Uint256]*transaction.Transaction
	SynchronizeTxnPool()
	Xmit(common.Inventory) error // The transmit interface
	GetEvent(eventName string) *events.Event
	GetMinersAddrs() ([]*crypto.PubKey, uint64)
}

func StartProtocol(pubKey *crypto.PubKey) (Neter, protocol.Noder) {
	seedNodes := config.Parameters.SeedList

	net := node.InitNode(pubKey)
	for _, nodeAddr := range seedNodes {
		go net.Connect(nodeAddr)
	}
	return net, net
}
