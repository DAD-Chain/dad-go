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

package router

import (
	"sync"

	"github.com/dad-go/eventbus/actor"
)

func spawn(id string, config RouterConfig, props *actor.Props, parent *actor.PID) (*actor.PID, error) {
	ref := &process{}
	proxy, absent := actor.ProcessRegistry.Add(ref, id)
	if !absent {
		return proxy, actor.ErrNameExists
	}

	var pc = *props
	pc.WithSpawnFunc(nil)
	ref.state = config.CreateRouterState()

	if config.RouterType() == GroupRouterType {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		ref.router, _ = actor.DefaultSpawner(id+"/router", actor.FromProducer(func() actor.Actor {
			return &groupRouterActor{
				props:  &pc,
				config: config,
				state:  ref.state,
				wg:     wg,
			}
		}), parent)
		wg.Wait() // wait for routerActor to start
	} else {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		ref.router, _ = actor.DefaultSpawner(id+"/router", actor.FromProducer(func() actor.Actor {
			return &poolRouterActor{
				props:  &pc,
				config: config,
				state:  ref.state,
				wg:     wg,
			}
		}), parent)
		wg.Wait() // wait for routerActor to start
	}

	return proxy, nil
}
