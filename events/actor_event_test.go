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
	"fmt"
	"testing"
	"time"

	"github.com/ontio/dad-go-eventbus/actor"
)

const testTopic = "test"

type testMessage struct {
	Message string
}

func testSubReceive(c actor.Context) {
	switch msg := c.Message().(type) {
	case *testMessage:
		fmt.Printf("PID:%s receive message:%s\n", c.Self().Id, msg.Message)
	}
}

func TestActorEvent(t *testing.T) {
	Init()
	subPID1 := actor.Spawn(actor.FromFunc(testSubReceive))
	subPID2 := actor.Spawn(actor.FromFunc(testSubReceive))
	sub1 := NewActorSubscriber(subPID1)
	sub2 := NewActorSubscriber(subPID2)
	sub1.Subscribe(testTopic)
	sub2.Subscribe(testTopic)
	DefActorPublisher.Publish(testTopic, &testMessage{Message: "Hello"})
	time.Sleep(time.Millisecond)
	DefActorPublisher.Publish(testTopic, &testMessage{Message: "Word"})
}
