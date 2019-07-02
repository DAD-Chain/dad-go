package node

import (
	"github.com/DAD-Chain/dad-go/events"
	"fmt"
)

type eventQueue struct {
	Consensus *events.Event
	Block     *events.Event
}

func (eq *eventQueue) init() {
	eq.Consensus = events.NewEvent()
	eq.Block = events.NewEvent()
}

func (eq *eventQueue) GetEvent(eventName string) *events.Event {
	switch eventName {
	case "consensus":
		return eq.Consensus
	case "block":
		return eq.Block
	default:
		fmt.Printf("Unknow event registe")
		return nil
	}
}
