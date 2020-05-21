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
	"time"

	"github.com/dad-go/common/log"
	"github.com/dad-go/eventbus/actor"
	zmq "github.com/pebbe/zmq4"
)

type endpointReader struct {
	suspended bool
}

func (s *endpointReader) Receive(stream *zmq.Socket) error {
	targets := make([]*actor.PID, 100)
	for {
		if s.suspended {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		batchstr, err := stream.Recv(0)
		if err != nil {
			log.Error("iEndpointReader failed to recieve.......", err.Error())
			return err
		}

		batchde, _ := Deserialize([]byte(batchstr), "zmqremote.MessageBatch", int32(0))

		batch := batchde.(*MessageBatch)
		//only grow pid lookup if needed
		if len(batch.TargetNames) > len(targets) {
			targets = make([]*actor.PID, len(batch.TargetNames))
		}

		for i := 0; i < len(batch.TargetNames); i++ {
			targets[i] = actor.NewLocalPID(batch.TargetNames[i])
		}

		for _, envelope := range batch.Envelopes {
			pid := targets[envelope.Target]
			message, err := Deserialize(envelope.MessageData, batch.TypeNames[envelope.TypeId], envelope.SerializerId)
			if err != nil {
				log.Error("EndpointReader failed to deserialize........", err)
				return err
			}
			//if message is system message send it as sysmsg instead of usermsg

			sender := envelope.Sender

			switch msg := message.(type) {
			case *actor.Terminated:
				rt := &remoteTerminate{
					Watchee: msg.Who,
					Watcher: pid,
				}
				endpointManager.remoteTerminate(rt)
			case actor.SystemMessage:
				ref, _ := actor.ProcessRegistry.GetLocal(pid.Id)
				ref.SendSystemMessage(pid, msg)
			default:
				var header map[string]string
				if envelope.MessageHeader != nil {
					header = envelope.MessageHeader.HeaderData
				}
				localEnvelope := &actor.MessageEnvelope{
					Header:  header,
					Message: message,
					Sender:  sender,
				}
				pid.Tell(localEnvelope)
			}
		}
	}
}

func (s *endpointReader) suspend(toSuspend bool) {
	s.suspended = toSuspend
}
