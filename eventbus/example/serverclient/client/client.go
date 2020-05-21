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

package client

import (
	"fmt"
	"time"

	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/eventbus/example/serverclient/message"
)

type Client struct{}

//Call the server synchronous
func (client *Client) SyncCall(serverPID *actor.PID) (interface{}, error) {
	//props := actor.FromProducer(func() actor.Actor { return &Client{} })
	//actor.Spawn(props)
	future := serverPID.RequestFuture(&message.Request{Who: "dad-go"}, 10*time.Second)
	result, err := future.Result()
	return result, err
}

//Call the server asynchronous
func (client *Client) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize client actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *message.Response:
		fmt.Println("Receive message", msg.Welcome)
	}
}

func (client *Client) AsyncCall(serverPID *actor.PID) *actor.PID {
	props := actor.FromProducer(func() actor.Actor { return &Client{} })
	clientPID := actor.Spawn(props)
	serverPID.Request(&message.Request{Who: "dad-go"}, clientPID)
	return clientPID
}

func (client *Client) Stop(pid *actor.PID) {
	pid.Stop()
}
