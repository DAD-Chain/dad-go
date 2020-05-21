/*
 * Copyright (C) 2018 The dad-go Authors
 * This file is part of The dad-go library.
 *
 * The dad-go is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The dad-go is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The dad-go.  If not, see <http://www.gnu.org/licenses/>.
 */

package events

import (
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/eventbus/eventhub"
)

var DefEvtHub *eventhub.EventHub
var DefPublisherPID *actor.PID
var DefActorPublisher *ActorPublisher
var defPublisherProps *actor.Props

func Init() {
	DefEvtHub = eventhub.GlobalEventHub
	defPublisherProps = actor.FromFunc(func(context actor.Context) {})
	DefPublisherPID = actor.Spawn(defPublisherProps)
	DefActorPublisher = NewActorPublisher(DefPublisherPID)
}

func NewActorPublisher(publisher *actor.PID, evtHub ...*eventhub.EventHub) *ActorPublisher {
	var hub *eventhub.EventHub
	if len(evtHub) == 0 {
		hub = DefEvtHub
	} else {
		hub = evtHub[0]
	}
	if publisher == nil {
		publisher = DefPublisherPID
	}
	return &ActorPublisher{
		EvtHub:    hub,
		Publisher: publisher,
	}
}

type ActorPublisher struct {
	EvtHub    *eventhub.EventHub
	Publisher *actor.PID
}

func (this *ActorPublisher) Publish(topic string, msg interface{}) {
	event := &eventhub.Event{
		Publisher: this.Publisher,
		Message:   msg,
		Topic:     topic,
		Policy:    eventhub.PUBLISH_POLICY_ALL,
	}
	this.EvtHub.Publish(event)
}

func (this *ActorPublisher) PublishEvent(evt *eventhub.Event) {
	this.EvtHub.Publish(evt)
}

type ActorSubscriber struct {
	EvtHub     *eventhub.EventHub
	Subscriber *actor.PID
}

func NewActorSubscriber(subscriber *actor.PID, evtHub ...*eventhub.EventHub) *ActorSubscriber {
	var hub *eventhub.EventHub
	if len(evtHub) == 0 {
		hub = DefEvtHub
	} else {
		hub = evtHub[0]
	}
	return &ActorSubscriber{
		EvtHub:     hub,
		Subscriber: subscriber,
	}
}

func (this *ActorSubscriber) Subscribe(topic string) {
	this.EvtHub.Subscribe(topic, this.Subscriber)
}

func (this *ActorSubscriber) Unsubscribe(topic string) {
	this.EvtHub.Unsubscribe(topic, this.Subscriber)
}
