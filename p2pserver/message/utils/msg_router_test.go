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

package utils

import (
	"testing"

	"github.com/ontio/dad-go-eventbus/actor"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/p2pserver/message/types"
	"github.com/ontio/dad-go/p2pserver/net/netserver"
	"github.com/ontio/dad-go/p2pserver/net/protocol"
	"github.com/stretchr/testify/assert"
)

func testHandler(data *types.MsgPayload, p2p p2p.P2P, pid *actor.PID, args ...interface{}) {
	log.Info("Test handler")
}

// TestMsgRouter tests a basic function of a message router
func TestMsgRouter(t *testing.T) {
	network := netserver.NewNetServer()
	msgRouter := NewMsgRouter(network)
	assert.NotNil(t, msgRouter)

	msgRouter.RegisterMsgHandler("test", testHandler)
	msgRouter.UnRegisterMsgHandler("test")
	msgRouter.Start()
	msgRouter.Stop()
}
