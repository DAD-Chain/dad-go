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

package types

import (
	"bytes"
	"crypto/sha256"
	"errors"

	"github.com/ontio/dad-go/common/config"
	"github.com/ontio/dad-go/common/log"
	"github.com/ontio/dad-go/p2pserver/common"
)

type Message interface {
	Verify([]byte) error
	Serialization() ([]byte, error)
	Deserialization([]byte) error
}

// The message body and header
type MsgCont struct {
	Hdr MsgHdr
	p   interface{}
}

type varStr struct {
	len uint
	buf []byte
}

//split msg type from msg hdr
func MsgType(buf []byte) (string, error) {
	cmd := buf[common.CMD_OFFSET : common.CMD_OFFSET+common.MSG_CMD_LEN]
	n := bytes.IndexByte(cmd, 0)
	if n < 0 || n >= common.MSG_CMD_LEN {
		return "", errors.New("unexpected length of CMD command")
	}
	s := string(cmd[:n])
	return s, nil
}

//check netmagic value
func magicVerify(magic uint32) bool {
	if magic != uint32(config.DefConfig.P2PNode.NetworkId) {
		log.Warnf("unmatched magic number 0x%0x", magic)
		return false
	}
	return true
}

//check wether header is valid
func ValidMsgHdr(buf []byte) bool {
	var h MsgHdr
	h.Deserialization(buf)
	return magicVerify(h.Magic)
}

//caculate payload length
func PayloadLen(buf []byte) (int, error) {
	var h MsgHdr
	err := h.Deserialization(buf)
	if err != nil {
		return 0, err
	}
	if int(h.Length) > common.MAX_MSG_LEN-common.MSG_HDR_LEN {
		return 0, errors.New("calculate PayloadLen error. buf length exceed max payload size")
	}
	return int(h.Length), nil
}

//caculate checksum value
func CheckSum(p []byte) []byte {
	t := sha256.Sum256(p)
	s := sha256.Sum256(t[:])

	// Currently we only need the front 4 bytes as checksum
	return s[:common.CHECKSUM_LEN]
}

// reverse the input
func Reverse(input []byte) []byte {
	if len(input) == 0 {
		return input
	}
	return append(Reverse(input[1:]), input[0])
}
