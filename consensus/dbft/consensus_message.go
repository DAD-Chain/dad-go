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
	"bytes"
	"errors"
	"github.com/dad-go/common/log"
	ser "github.com/dad-go/common/serialization"
	"io"
)

type ConsensusMessage interface {
	ser.SerializableData
	Type() ConsensusMessageType
	ViewNumber() byte
	ConsensusMessageData() *ConsensusMessageData
}

type ConsensusMessageData struct {
	Type       ConsensusMessageType
	ViewNumber byte
}

func DeserializeMessage(data []byte) (ConsensusMessage, error) {
	log.Debug()
	msgType := ConsensusMessageType(data[0])

	r := bytes.NewReader(data)
	switch msgType {
	case PrepareRequestMsg:
		prMsg := &PrepareRequest{}
		err := prMsg.Deserialize(r)
		if err != nil {
			log.Error("[DeserializeMessage] PrepareRequestMsg Deserialize Error: ", err.Error())
			return nil, err
		}
		return prMsg, nil

	case PrepareResponseMsg:
		presMsg := &PrepareResponse{}
		err := presMsg.Deserialize(r)
		if err != nil {
			log.Error("[DeserializeMessage] PrepareResponseMsg Deserialize Error: ", err.Error())
			return nil, err
		}
		return presMsg, nil
	case ChangeViewMsg:
		cv := &ChangeView{}
		err := cv.Deserialize(r)
		if err != nil {
			log.Error("[DeserializeMessage] ChangeViewMsg Deserialize Error: ", err.Error())
			return nil, err
		}
		return cv, nil

	case BlockSignaturesMsg:
		blockSigs := &BlockSignatures{}
		err := blockSigs.Deserialize(r)
		if err != nil {
			log.Error("[DeserializeMessage] BlockSignaturesMsg Deserialize Error: ", err.Error())
			return nil, err
		}

		return blockSigs, nil
	}

	return nil, errors.New("The message is invalid.")
}

func (cd *ConsensusMessageData) Serialize(w io.Writer) {
	log.Debug()
	//ConsensusMessageType
	w.Write([]byte{byte(cd.Type)})

	//ViewNumber
	w.Write([]byte{byte(cd.ViewNumber)})

}

//read data to reader
func (cd *ConsensusMessageData) Deserialize(r io.Reader) error {
	log.Debug()
	//ConsensusMessageType
	var msgType [1]byte
	_, err := io.ReadFull(r, msgType[:])
	if err != nil {
		return err
	}
	cd.Type = ConsensusMessageType(msgType[0])

	//ViewNumber
	var vNumber [1]byte
	_, err = io.ReadFull(r, vNumber[:])
	if err != nil {
		return err
	}
	cd.ViewNumber = vNumber[0]

	return nil
}
