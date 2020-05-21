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

package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dad-go/eventbus/actor"
	. "github.com/dad-go/eventbus/example/services/messages"
	"github.com/dad-go/eventbus/example/services/serviceA"
	"github.com/dad-go/eventbus/example/services/serviceB"
)

func main() {
	sva := actor.FromProducer(func() actor.Actor { return &serviceA.ServiceA{} })
	svb := actor.FromProducer(func() actor.Actor { return &serviceB.ServiceB{} })

	pipA, _ := actor.SpawnNamed(sva, "serviceA")
	pipB, _ := actor.SpawnNamed(svb, "serviceB")

	pipA.Request(&ServiceARequest{"TESTA"}, pipB)

	pipB.Request(&ServiceBRequest{"TESTB"}, pipA)
	time.Sleep(2 * time.Second)

	f := pipA.RequestFuture(1, 50*time.Microsecond)
	result, err := f.Result()
	if err != nil {
		fmt.Println("errors:", err.Error())
	}
	fmt.Println("get sync call result :" + strconv.Itoa(result.(int)))

}
