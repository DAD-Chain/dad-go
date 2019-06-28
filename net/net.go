package net

import (
	"dad-go/common"
	"dad-go/config"
	"dad-go/events"
	"dad-go/core/transaction"
	"dad-go/net/node"
)

type Neter interface {
	GetMemoryPool() map[common.Uint256]*transaction.Transaction
	SynchronizeMemoryPool()
	Xmit(common.Inventory) error // The transmit interface
	GetEvent(eventName string) *events.Event
}

func StartProtocol() Neter {
	seedNodes := config.Parameters.SeedList

	net := node.InitNode()
	for _, nodeAddr := range seedNodes {
		net.Connect(nodeAddr)
	}
	return net
}
