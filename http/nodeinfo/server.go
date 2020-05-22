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

package nodeinfo

import (
	"fmt"
	"github.com/dad-go/common/config"
	"github.com/dad-go/core/ledger"
	. "github.com/dad-go/net/protocol"
	"html/template"
	"net/http"
	"sort"
	"strconv"
)

type Info struct {
	NodeVersion   string
	BlockHeight   uint32
	NeighborCnt   int
	Neighbors     []NgbNodeInfo
	HttpRestPort  int
	HttpWsPort    int
	HttpJsonPort  int
	HttpLocalPort int
	NodePort      int
	NodeId        string
	NodeType      string
}

const (
	verifyNode = "Verify Node"
	serviceNode = "Service Node"
)

var node Noder

var templates = template.Must(template.New("info").Parse(page))

func newNgbNodeInfo(ngbId string, ngbType string, ngbAddr string, httpInfoAddr string, httpInfoPort int, httpInfoStart bool) *NgbNodeInfo {
	return &NgbNodeInfo{NgbId: ngbId, NgbType: ngbType, NgbAddr: ngbAddr, HttpInfoAddr: httpInfoAddr,
		HttpInfoPort: httpInfoPort, HttpInfoStart: httpInfoStart}
}

func initPageInfo(blockHeight uint32, curNodeType string, ngbrCnt int, ngbrsInfo []NgbNodeInfo) (*Info, error) {
	id := fmt.Sprintf("0x%x", node.GetID())
	return &Info{NodeVersion: config.Version, BlockHeight: blockHeight,
		NeighborCnt: ngbrCnt, Neighbors: ngbrsInfo,
		HttpRestPort:  config.Parameters.HttpRestPort,
		HttpWsPort:    config.Parameters.HttpWsPort,
		HttpJsonPort:  config.Parameters.HttpJsonPort,
		HttpLocalPort: config.Parameters.HttpLocalPort,
		NodePort:      config.Parameters.NodePort,
		NodeId:        id, NodeType: curNodeType}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	var ngbrNodersInfo []NgbNodeInfo
	var ngbId string
	var ngbAddr string
	var ngbType string
	var ngbInfoPort int
	var ngbInfoState bool
	var ngbHttpInfoAddr string

	curNodeType := serviceNode
	bookkeeperState,  _ := ledger.DefLedger.GetBookkeeperState()
	bookkeepers := bookkeeperState.CurrBookkeeper
	bookkeeperLen := len(bookkeepers)
	for i := 0; i < bookkeeperLen; i++ {
		if node.GetPubKey().X.Cmp(bookkeepers[i].X) == 0 {
			curNodeType = verifyNode
			break
		}
	}

	ngbrNoders := node.GetNeighborNoder()
	ngbrsLen := len(ngbrNoders)
	for i := 0; i < ngbrsLen; i++ {
		ngbType = serviceNode
		for j := 0; j < bookkeeperLen; j++ {
			if ngbrNoders[i].GetPubKey().X.Cmp(bookkeepers[j].X) == 0 {
				ngbType = verifyNode
				break
			}
		}

		ngbAddr = ngbrNoders[i].GetAddr()
		ngbInfoPort = ngbrNoders[i].GetHttpInfoPort()
		ngbInfoState = ngbrNoders[i].GetHttpInfoState()
		ngbHttpInfoAddr = ngbAddr + ":" + strconv.Itoa(ngbInfoPort)
		ngbId = fmt.Sprintf("0x%x", ngbrNoders[i].GetID())

		ngbrInfo := newNgbNodeInfo(ngbId, ngbType, ngbAddr, ngbHttpInfoAddr, ngbInfoPort, ngbInfoState)
		ngbrNodersInfo = append(ngbrNodersInfo, *ngbrInfo)
	}
	sort.Sort(NgbNodeInfoSlice(ngbrNodersInfo))

	blockHeight := ledger.DefLedger.GetCurrentBlockHeight()
	pageInfo, err := initPageInfo(blockHeight, curNodeType, ngbrsLen, ngbrNodersInfo)
	if err != nil {
		http.Redirect(w, r, "/info", http.StatusFound)
		return
	}

	err = templates.ExecuteTemplate(w, "info", pageInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func StartServer(n Noder) {
	node = n
	port := int(config.Parameters.HttpInfoPort)
	http.HandleFunc("/info", viewHandler)
	http.ListenAndServe(":" + strconv.Itoa(port), nil)
}
