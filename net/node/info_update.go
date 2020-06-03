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

package node

import (
	"github.com/dad-go/common/config"
	"github.com/dad-go/common/log"
	"math/rand"
	"net"
	//"github.com/dad-go/core/ledger"
	"github.com/dad-go/net/actor"
	. "github.com/dad-go/net/message"
	. "github.com/dad-go/net/protocol"
	"strconv"
	"time"
)

func (node *node) GetBlkHdrs() {
	if !node.IsUptoMinNodeCount() {
		return
	}
	noders := node.local.GetNeighborNoder()
	if len(noders) == 0 {
		return
	}
	nodelist := []Noder{}
	for _, v := range noders {
		height, _ := actor.GetCurrentHeaderHeight()
		if uint64(height) < v.GetHeight() {
			nodelist = append(nodelist, v)
		}
	}
	ncout := len(nodelist)
	if ncout == 0 {
		return
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(ncout)
	n := nodelist[index]
	SendMsgSyncHeaders(n)
}

func (node *node) SyncBlk() {
	headerHeight, _ := actor.GetCurrentHeaderHeight()
	currentBlkHeight, _ := actor.GetCurrentBlockHeight()
	if currentBlkHeight >= headerHeight {
		return
	}
	var dValue int32
	var reqCnt uint32
	var i uint32
	noders := node.local.GetNeighborNoder()

	for _, n := range noders {
		if uint32(n.GetHeight()) <= currentBlkHeight {
			continue
		}
		n.RemoveFlightHeightLessThan(currentBlkHeight)
		count := MAX_REQ_BLK_ONCE - uint32(n.GetFlightHeightCnt())
		dValue = int32(headerHeight - currentBlkHeight - reqCnt)
		flights := n.GetFlightHeights()
		if count == 0 {
			for _, f := range flights {
				hash, _ := actor.GetBlockHashByHeight(f)
				isContainBlock, _ := actor.IsContainBlock(hash)
				if isContainBlock == false {
					ReqBlkData(n, hash)
				}
			}

		}
		for i = 1; i <= count && dValue >= 0; i++ {
			hash, _ := actor.GetBlockHashByHeight(currentBlkHeight + reqCnt)
			isContainBlock, _ := actor.IsContainBlock(hash)
			if isContainBlock == false {
				ReqBlkData(n, hash)
				n.StoreFlightHeight(currentBlkHeight + reqCnt)
			}
			reqCnt++
			dValue--
		}
	}
}

func (node *node) SendPingToNbr() {
	noders := node.local.GetNeighborNoder()
	for _, n := range noders {
		if n.GetState() == ESTABLISH {
			buf, err := NewPingMsg()
			if err != nil {
				log.Error("failed build a new ping message")
			} else {
				go n.Tx(buf)
			}
		}
	}
}

func (node *node) HeartBeatMonitor() {
	noders := node.local.GetNeighborNoder()
	var periodUpdateTime uint
	if config.Parameters.GenBlockTime > config.MIN_GEN_BLOCK_TIME {
		periodUpdateTime = config.Parameters.GenBlockTime / TIMES_OF_UPDATE_TIME
	} else {
		periodUpdateTime = config.DEFAULT_GEN_BLOCK_TIME / TIMES_OF_UPDATE_TIME
	}
	for _, n := range noders {
		if n.GetState() == ESTABLISH {
			t := n.GetLastRXTime()
			if t.Before(time.Now().Add(-1 * time.Second * time.Duration(periodUpdateTime) * KEEPALIVE_TIMEOUT)) {
				log.Warn("keepalive timeout!!!")
				n.SetState(INACTIVITY)
				n.CloseConn()
			}
		}
	}
}

func (node *node) ReqNeighborList() {
	buf, _ := NewMsg("getaddr", node.local)
	go node.Tx(buf)
}

func (node *node) ConnectSeeds() {
	if node.IsUptoMinNodeCount() {
		return
	}
	seedNodes := config.Parameters.SeedList
	for _, nodeAddr := range seedNodes {
		found := false
		var n Noder
		var ip net.IP
		node.nbrNodes.Lock()
		for _, tn := range node.nbrNodes.List {
			addr := getNodeAddr(tn)
			ip = addr.IpAddr[:]
			addrstring := ip.To16().String() + ":" + strconv.Itoa(int(addr.Port))
			if nodeAddr == addrstring {
				n = tn
				found = true
				break
			}
		}
		node.nbrNodes.Unlock()
		if found {
			if n.GetState() == ESTABLISH {
				n.ReqNeighborList()
			}
		} else { //not found
			go node.Connect(nodeAddr)
		}
	}
}

func getNodeAddr(n *node) NodeAddr {
	var addr NodeAddr
	addr.IpAddr, _ = n.GetAddr16()
	addr.Time = n.GetTime()
	addr.Services = n.Services()
	addr.Port = n.GetPort()
	addr.ID = n.GetID()
	return addr
}

func (node *node) reconnect() {
	node.RetryConnAddrs.Lock()
	defer node.RetryConnAddrs.Unlock()
	lst := make(map[string]int)
	for addr := range node.RetryAddrs {
		node.RetryAddrs[addr] = node.RetryAddrs[addr] + 1
		rand.Seed(time.Now().UnixNano())
		log.Trace("Try to reconnect peer, peer addr is ", addr)
		<-time.After(time.Duration(rand.Intn(CONN_MAX_BACK)) * time.Millisecond)
		log.Trace("Back off time`s up, start connect node")
		node.Connect(addr)
		if node.RetryAddrs[addr] < MAX_RETRY_COUNT {
			lst[addr] = node.RetryAddrs[addr]
		}
	}
	node.RetryAddrs = lst

}

func (n *node) TryConnect() {
	if n.fetchRetryNodeFromNeiborList() > 0 {
		n.reconnect()
	}
}

func (n *node) fetchRetryNodeFromNeiborList() int {
	n.nbrNodes.Lock()
	defer n.nbrNodes.Unlock()
	var ip net.IP
	neibornodes := make(map[uint64]*node)
	for _, tn := range n.nbrNodes.List {
		addr := getNodeAddr(tn)
		ip = addr.IpAddr[:]
		nodeAddr := ip.To16().String() + ":" + strconv.Itoa(int(addr.Port))
		if tn.GetState() == INACTIVITY {
			//add addr to retry list
			n.AddInRetryList(nodeAddr)
			//close legacy node
			if tn.conn != nil {
				tn.CloseConn()
			}
		} else {
			//add others to tmp node map
			n.RemoveFromRetryList(nodeAddr)
			neibornodes[tn.GetID()] = tn
		}
	}
	n.nbrNodes.List = neibornodes
	return len(n.RetryAddrs)
}

// FIXME part of node info update function could be a node method itself intead of
// a node map method
// Fixme the Nodes should be a parameter
func (node *node) updateNodeInfo() {
	var periodUpdateTime uint
	if config.Parameters.GenBlockTime > config.MIN_GEN_BLOCK_TIME {
		periodUpdateTime = config.Parameters.GenBlockTime / TIMES_OF_UPDATE_TIME
	} else {
		periodUpdateTime = config.DEFAULT_GEN_BLOCK_TIME / TIMES_OF_UPDATE_TIME
	}
	ticker := time.NewTicker(time.Second * (time.Duration(periodUpdateTime)))
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			node.SendPingToNbr()
			node.GetBlkHdrs()
			node.SyncBlk()
			node.HeartBeatMonitor()
		case <-quit:
			ticker.Stop()
			return
		}
	}
	// TODO when to close the timer
}

func (node *node) updateConnection() {
	t := time.NewTimer(time.Second * CONN_MONITOR)
	for {
		select {
		case <-t.C:
			node.ConnectSeeds()
			node.TryConnect()
			t.Stop()
			t.Reset(time.Second * CONN_MONITOR)
		}
	}

}
