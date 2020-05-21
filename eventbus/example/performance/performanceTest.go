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
	"runtime"
	"time"

	"github.com/dad-go/eventbus/actor"
)

type Ball struct {
	val int
}

var start, end int64

//func Benchmark_Division1(b *testing.B){
func main() {
	fmt.Printf("test performance")
	runtime.GOMAXPROCS(4)
	times := 10000000
	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {

		case Ball:
			val := msg.val
			if val < times {
				context.Sender().Request(Ball{val: val + 1}, context.Self())
			} else {
				end = time.Now().UnixNano()
				fmt.Printf("end at time %d\n", end)
			}
		default:
		}
	})
	playerA, _ := actor.SpawnNamed(props, "playerA")
	playerB, _ := actor.SpawnNamed(props, "playerB")
	start = time.Now().UnixNano()
	fmt.Println("start at time:", start)
	playerA.Request(Ball{val: 1}, playerB)
	time.Sleep(10000 * time.Millisecond)
	fmt.Printf("run time:%d     elapsed time:%d ms", times, (end-start)/1000000)
}
