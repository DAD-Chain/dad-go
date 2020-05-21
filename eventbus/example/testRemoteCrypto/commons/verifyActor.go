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

package commons

import (
	"fmt"
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/crypto"
	"bytes"
)



type VerifyActor struct{
	Count int
}

func (s *VerifyActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")

	case *VerifyRequest:
		s.Count++
		//fmt.Println(context.Self().Id, "is verifying...")
		crypto.SetAlg("")
		buf := bytes.NewBuffer(msg.PublicKey)
		pubKey := new(crypto.PubKey)
		err := pubKey.DeSerialize(buf)
		if err != nil {
			fmt.Println("DeSerialize failed.", err)
		}
		err = crypto.Verify(*pubKey,msg.Data,msg.Signature)
		//fmt.Println(context.Self().Id, "done verifying...")
		if err != nil{
			fmt.Println("verify error :", err)
			response:=&VerifyResponse{Seq:msg.Seq,Result:false,ErrorMsg:err.Error()}
			context.Sender().Tell(response)
		}else{
			response:=&VerifyResponse{Seq:msg.Seq,Result:true,ErrorMsg:""}
			context.Sender().Tell(response)
		}
		if s.Count%1 == 0 {
			fmt.Println(s.Count)
		}

	default:
		fmt.Printf("---unknown message%v\n",msg)
	}
}