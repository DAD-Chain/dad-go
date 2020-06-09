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

package dbft

import (
	"io"

	ser "github.com/ontio/dad-go/common/serialization"
)

type ChangeView struct {
	msgData       ConsensusMessageData
	NewViewNumber byte
}

func (cv *ChangeView) Serialize(w io.Writer) error {
	cv.msgData.Serialize(w)
	w.Write([]byte{cv.NewViewNumber})
	return nil
}

//read data to reader
func (cv *ChangeView) Deserialize(r io.Reader) error {
	cv.msgData.Deserialize(r)
	viewNum, err := ser.ReadBytes(r, 1)
	if err != nil {
		return err
	}
	cv.NewViewNumber = viewNum[0]
	return nil
}

func (cv *ChangeView) Type() ConsensusMessageType {
	return cv.ConsensusMessageData().Type
}

func (cv *ChangeView) ViewNumber() byte {
	return cv.msgData.ViewNumber
}

func (cv *ChangeView) ConsensusMessageData() *ConsensusMessageData {
	return &(cv.msgData)
}
