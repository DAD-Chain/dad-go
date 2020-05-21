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

package actor

import (
	"github.com/ontio/dad-go-eventbus/actor"
	"github.com/dad-go/events/message"
	"github.com/dad-go/events"
)

type EventActor struct{
	blockPersistCompleted func(v interface{})
	smartCodeEvt func(v interface{})
}

func (t *EventActor) Receive(c actor.Context) {
	switch msg := c.Message().(type) {
	case *message.SaveBlockCompleteMsg:
		t.blockPersistCompleted(*msg.Block)
	case *message.SmartCodeEventMsg:
		t.smartCodeEvt(*msg.Event)
	default:
		//fmt.Println(msg)
	}
}

func SubscribeEvent(topic string,handler func(v interface{})) {
	var props = actor.FromProducer(func() actor.Actor {
		if topic == message.TopicSaveBlockComplete{
			return &EventActor{blockPersistCompleted:handler}
		}else if topic == message.TopicSmartCodeEvent{
			return &EventActor{smartCodeEvt:handler}
		}else{
			return &EventActor{}
		}
	})
	var pid = actor.Spawn(props)
	var sub = events.NewActorSubscriber(pid)
	sub.Subscribe(topic)
}
