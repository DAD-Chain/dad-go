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
	"github.com/ontio/dad-go-crypto/keypair"
	types "github.com/ontio/dad-go/p2pserver/common"
)

//stop net server
type StopServerReq struct {
}

//response of stop request
type StopServerRsp struct {
	Error error
}

//version request
type GetVersionReq struct {
}

//response of version request
type GetVersionRsp struct {
	Version uint32
}

//connection count requet
type GetConnectionCntReq struct {
}

//response of connection count requet
type GetConnectionCntRsp struct {
	Cnt uint32
}

//get net module id
type GetIdReq struct {
}

//response of net module id
type GetIdRsp struct {
	Id uint64
}

//get net sync port
type GetSyncPortReq struct {
}

//response of net sync port
type GetSyncPortRsp struct {
	SyncPort uint16
}

//get net cons port
type GetConsPortReq struct {
}

//response of net cons port
type GetConsPortRsp struct {
	ConsPort uint16
}

//get connection port requet
type GetPortReq struct {
}

//response of connection port requet
type GetPortRsp struct {
	SyncPort uint16
	ConsPort uint16
}

//get connection state requet
type GetConnectionStateReq struct {
}

//response of connection state requet
type GetConnectionStateRsp struct {
	State uint32
}

//get net timestamp request
type GetTimeReq struct {
}

//response of net timestamp
type GetTimeRsp struct {
	Time int64
}

type GetNodeTypeReq struct {
}
type GetNodeTypeRsp struct {
	NodeType uint64
}

//whether net can relay
type GetRelayStateReq struct {
}

//response of whether net can relay
type GetRelayStateRsp struct {
	Relay bool
}

//get all nbr`s address request
type GetNeighborAddrsReq struct {
}

//response of all nbr`s address
type GetNeighborAddrsRsp struct {
	Addrs []types.PeerAddr
}

type TransmitConsensusMsgReq struct {
	Target keypair.PublicKey
	Msg    []byte
}
