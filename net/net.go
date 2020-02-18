package net

import (
	. "github.com/dad-go/common"
	"github.com/dad-go/core/ledger"
	"github.com/dad-go/core/transaction"
	"github.com/dad-go/crypto"
	. "github.com/dad-go/errors"
	"github.com/dad-go/events"
	"github.com/dad-go/net/node"
	"github.com/dad-go/net/protocol"
)

type Neter interface {
	GetTxnPool(byCount bool) map[Uint256]*transaction.Transaction
	Xmit(interface{}) error
	GetEvent(eventName string) *events.Event
	GetBookKeepersAddrs() ([]*crypto.PubKey, uint64)
	CleanSubmittedTransactions(block *ledger.Block) error
	GetNeighborNoder() []protocol.Noder
	Tx(buf []byte)
	AppendTxnPool(*transaction.Transaction) ErrCode
}

func StartProtocol(pubKey *crypto.PubKey) protocol.Noder {
	net := node.InitNode(pubKey)
	net.ConnectSeeds()

	return net
}
