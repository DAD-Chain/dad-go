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

package remote

import "github.com/dad-go/eventbus/actor"

type EndpointTerminatedEvent struct {
	Address string
}

type EndpointConnectedEvent struct {
	Address string
}

type remoteWatch struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

type remoteUnwatch struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

type remoteDeliver struct {
	header       actor.ReadonlyMessageHeader
	message      interface{}
	target       *actor.PID
	sender       *actor.PID
	serializerID int32
}

type remoteTerminate struct {
	Watcher *actor.PID
	Watchee *actor.PID
}

type JsonMessage struct {
	TypeName string
	Json     string
}

var (
	stopMessage interface{} = &actor.Stop{}
)

var (
	ActorPidRespErr         interface{} = &ActorPidResponse{StatusCode: ResponseStatusCodeERROR.ToInt32()}
	ActorPidRespTimeout     interface{} = &ActorPidResponse{StatusCode: ResponseStatusCodeTIMEOUT.ToInt32()}
	ActorPidRespUnavailable interface{} = &ActorPidResponse{StatusCode: ResponseStatusCodeUNAVAILABLE.ToInt32()}
)
