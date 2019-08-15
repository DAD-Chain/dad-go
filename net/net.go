package net

import (
	. "dad-go/common"
	"dad-go/config"
	"dad-go/core/ledger"
	"dad-go/core/transaction"
	"dad-go/crypto"
	"dad-go/events"
	"dad-go/net/node"
	"dad-go/net/protocol"
)

type Neter interface {
	GetTxnPool(cleanPool bool) map[Uint256]*transaction.Transaction
	SynchronizeTxnPool()
	Xmit(interface{}) error
	GetEvent(eventName string) *events.Event
	GetBookKeepersAddrs() ([]*crypto.PubKey, uint64)
	CleanSubmittedTransactions(block *ledger.Block) error
	GetNeighborNoder() []protocol.Noder
	Tx(buf []byte)
}

func StartProtocol(pubKey *crypto.PubKey) (Neter, protocol.Noder) {
	seedNodes := config.Parameters.SeedList

	net := node.InitNode(pubKey)
	for _, nodeAddr := range seedNodes {
		go net.Connect(nodeAddr)
	}
	return net, net
}
