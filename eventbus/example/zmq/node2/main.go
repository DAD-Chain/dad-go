package main

import (
	"runtime"

	"fmt"

	"time"

	"github.com/dad-go/common/log"
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/eventbus/example/zmq/messages"
	"github.com/dad-go/eventbus/mailbox"
	"github.com/dad-go/eventbus/zmqremote"
)

func main() {
	log.Init()
	log.Debug("test")
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	zmqremote.Start("127.0.0.1:8080")

	var sender *actor.PID
	props := actor.
		FromFunc(
			func(context actor.Context) {
				switch msg := context.Message().(type) {
				case *messages.StartRemote:
					//fmt.Println("Done server!")
					fmt.Println("Starting")
					sender = msg.Sender
					context.Respond(&messages.Start{})
				case *messages.Ping:
					//fmt.Println("ping")
					sender.Tell(&messages.Pong{})
				}
			}).
		WithMailbox(mailbox.Bounded(1000000))

	pid, _ := actor.SpawnNamed(props, "remote")
	fmt.Println(pid)

	for {
		time.Sleep(1 * time.Second)
	}
}
