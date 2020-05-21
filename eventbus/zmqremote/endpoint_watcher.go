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

package zmqremote

import (
	"fmt"

	"github.com/dad-go/common/log"
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/eventbus/eventhub"
)

func newEndpointWatcher(address string) actor.Producer {
	return func() actor.Actor {
		return &endpointWatcher{
			address: address,
		}
	}
}

type endpointWatcher struct {
	address string
	watched map[string]*actor.PIDSet //key is the watching PID string, value is the watched PID
}

func (state *endpointWatcher) initialize() {
	log.Info("Started EndpointWatcher", fmt.Sprintf("address:%s", state.address))
	state.watched = make(map[string]*actor.PIDSet)
}

func (state *endpointWatcher) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		state.initialize()

	case *remoteTerminate:
		//delete the watch entries
		if pidSet, ok := state.watched[msg.Watcher.Id]; ok {
			pidSet.Remove(msg.Watchee)
			if pidSet.Len() == 0 {
				delete(state.watched, msg.Watcher.Id)
			}
		}

		terminated := &actor.Terminated{
			Who:               msg.Watchee,
			AddressTerminated: false,
		}
		ref, ok := actor.ProcessRegistry.GetLocal(msg.Watcher.Id)
		if ok {
			ref.SendSystemMessage(msg.Watcher, terminated)
		}
	case *EndpointConnectedEvent:
		//Already connected, pass
	case *EndpointTerminatedEvent:
		log.Info("EndpointWatcher handling terminated", fmt.Sprintf("address:%s", state.address))

		for id, pidSet := range state.watched {
			//try to find the watcher ID in the local actor registry
			ref, ok := actor.ProcessRegistry.GetLocal(id)
			if ok {
				pidSet.ForEach(func(i int, pid actor.PID) {
					//create a terminated event for the Watched actor
					eventhub.GlobalEventHub.RemovePID(pid)
					terminated := &actor.Terminated{
						Who:               &pid,
						AddressTerminated: true,
					}

					watcher := actor.NewLocalPID(id)
					//send the address Terminated event to the Watcher
					ref.SendSystemMessage(watcher, terminated)
				})
			}
		}

		//Clear watcher's map
		state.watched = make(map[string]*actor.PIDSet)
		ctx.SetBehavior(state.Terminated)

	case *remoteWatch:
		//add watchee to watcher's map
		if pidSet, ok := state.watched[msg.Watcher.Id]; ok {
			pidSet.Add(msg.Watchee)
		} else {
			state.watched[msg.Watcher.Id] = actor.NewPIDSet(msg.Watchee)
		}

		//recreate the Watch command
		w := &actor.Watch{
			Watcher: msg.Watcher,
		}

		//pass it off to the remote PID
		SendMessage(msg.Watchee, nil, w, nil, -1)

	case *remoteUnwatch:
		//delete the watch entries
		if pidSet, ok := state.watched[msg.Watcher.Id]; ok {
			pidSet.Remove(msg.Watchee)
			if pidSet.Len() == 0 {
				delete(state.watched, msg.Watcher.Id)
			}
		}

		//recreate the Unwatch command
		uw := &actor.Unwatch{
			Watcher: msg.Watcher,
		}

		//pass it off to the remote PID
		SendMessage(msg.Watchee, nil, uw, nil, -1)
	case actor.SystemMessage, actor.AutoReceiveMessage:
		//ignore
	default:
		log.Error("EndpointWatcher received unknown message", fmt.Sprintf("address:%s", state.address), fmt.Sprintf("msg:%v", msg))
	}
}

func (state *endpointWatcher) Terminated(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *remoteWatch:
		//try to find the watcher ID in the local actor registry
		ref, ok := actor.ProcessRegistry.GetLocal(msg.Watcher.Id)
		if ok {

			//create a terminated event for the Watched actor
			terminated := &actor.Terminated{
				Who:               msg.Watchee,
				AddressTerminated: true,
			}
			//send the address Terminated event to the Watcher
			ref.SendSystemMessage(msg.Watcher, terminated)
		}
	case *EndpointConnectedEvent:
		log.Info("EndpointWatcher handling restart", fmt.Sprintf("address:%s", state.address))
		ctx.SetBehavior(state.Receive)
	case *remoteTerminate, *EndpointTerminatedEvent, *remoteUnwatch:
		// pass
		log.Error("EndpointWatcher receive message for already terminated endpoint", fmt.Sprintf("address:%s", state.address), fmt.Sprintf("msg:%v", msg))
	case actor.SystemMessage, actor.AutoReceiveMessage:
		//ignore
	default:
		log.Error("EndpointWatcher received unknown message", fmt.Sprintf("address:%s", state.address), fmt.Sprintf("msg:%v", msg))
	}
}
