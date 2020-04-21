package eventhub

import (
	"math/rand"

	"github.com/dad-go/common/log"
	"github.com/dad-go/eventbus/actor"
	"github.com/orcaman/concurrent-map"
)

type PublishPolicy int

type RoundRobinState struct {
	state map[string]int
}

const (
	PUBLISH_POLICY_ALL = iota
	PUBLISH_POLICY_ROUNDROBIN
	PUBLISH_POLICY_RANDOM
)

type EventHub struct {
	//sync.RWMutex
	Subscribers cmap.ConcurrentMap
	RoundRobinState
}

type Event struct {
	Publisher *actor.PID
	Topic     string
	Message   interface{}
	Policy    PublishPolicy
}

var GlobalEventHub = &EventHub{Subscribers: cmap.New(), RoundRobinState: RoundRobinState{make(map[string]int)}}

func (this *EventHub) Publish(event *Event) {
	//go func() {
	actors, ok := this.Subscribers.Get(event.Topic)
	if !ok {
		log.Info("no subscribers yet!")
		return
	}
	subscribers := actors.([]*actor.PID)
	this.sendEventByPolicy(subscribers, event, this.RoundRobinState)
	//}()
}

func (this *EventHub) Subscribe(topic string, subscriber *actor.PID) {
	subscribers, _ := this.Subscribers.Get(topic)

	//defer this.RWMutex.Unlock()
	//this.RWMutex.Lock()
	if subscribers == nil {
		this.Subscribers.Set(topic, []*actor.PID{subscriber})
	} else {
		this.Subscribers.Set(topic, append(subscribers.([]*actor.PID), subscriber))
	}

}

func (this *EventHub) Unsubscribe(topic string, subscriber *actor.PID) {

	tmpslice, ok := this.Subscribers.Get(topic)
	if !ok {
		log.Debug("No subscriber on topic:%s yet.\n", topic)
		return
	}
	//defer this.RWMutex.Unlock()
	//this.RWMutex.Lock()
	subscribers := tmpslice.([]*actor.PID)
	for i, s := range subscribers {
		if s == subscriber {
			this.Subscribers.Set(topic, append(subscribers[0:i], subscribers[i+1:]...))
			return
		}
	}

}

func (this *EventHub) sendEventByPolicy(subscribers []*actor.PID, event *Event, state RoundRobinState) {
	switch event.Policy {
	case PUBLISH_POLICY_ALL:
		for _, subscriber := range subscribers {
			subscriber.Request(event.Message, event.Publisher)
		}
	case PUBLISH_POLICY_RANDOM:
		length := len(subscribers)
		if length == 0 {
			log.Info("no subscribers yet!")
			return
		}
		var i int
		i = rand.Intn(length)
		subscribers[i].Request(event.Message, event.Publisher)
	case PUBLISH_POLICY_ROUNDROBIN:
		latestIdx := state.state[event.Topic]
		i := latestIdx + 1
		if i < 0 {
			latestIdx = 0
			i = 0
		}
		state.state[event.Topic] = i
		mod := len(subscribers)
		subscribers[i%mod].Request(event.Message, event.Publisher)
	}
}

func (this *EventHub) RemovePID(pid actor.PID) {
	if this.Subscribers.Count() == 0 {
		return
	}
	keys := this.Subscribers.Keys()
	for index, _ := range keys {
		this.Unsubscribe(keys[index], &pid)
	}
}
