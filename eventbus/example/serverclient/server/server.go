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

package server

import (
	"fmt"

	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/eventbus/example/serverclient/message"
)

type Server struct{}

func (server *Server) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize server actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *message.Request:
		fmt.Println("Receive message", msg.Who)
		context.Sender().Request(&message.Response{Welcome: "Welcome!"}, context.Self())
	}
}

func (server *Server) Start() *actor.PID {
	props := actor.FromProducer(func() actor.Actor { return &Server{} })
	pid := actor.Spawn(props)
	return pid
}

func (server *Server) Stop(pid *actor.PID) {
	pid.Stop()
}
