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

package vbft

import "testing"

func TestAddMsg(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, uint64(1))
	block, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed :%v", err)
		return
	}
	blockproposalmsg := &blockProposalMsg{
		Block: block,
	}
	err = msgpool.AddMsg(blockproposalmsg)
	t.Logf("TestAddMsg %v", err)
}

func TestHasMsg(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, uint64(1))
	block, err := constructBlock()
	if err != nil {
		t.Errorf("constructBlock failed :%v", err)
		return
	}
	blockproposalmsg := &blockProposalMsg{
		Block: block,
	}
	status := msgpool.HasMsg(blockproposalmsg)
	t.Logf("TestHasMsg: %v", status)
}

func TestGetProposalMsgs(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, uint64(1))
	consensusmsgs := msgpool.GetProposalMsgs(uint64(1))
	t.Logf("TestGetProposalMsgs: %v", len(consensusmsgs))
}

func TestGetEndorsementsMsgs(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, uint64(1))
	consensusmsgs := msgpool.GetEndorsementsMsgs(uint64(1))
	t.Logf("TestGetEndorsementsMsgs: %v", len(consensusmsgs))
}

func TestGetCommitMsgs(t *testing.T) {
	server := constructServer()
	msgpool := newMsgPool(server, uint64(1))
	consensusmsgs := msgpool.GetCommitMsgs(uint64(1))
	t.Logf("TestGetCommitMsgs: %v", len(consensusmsgs))
}
