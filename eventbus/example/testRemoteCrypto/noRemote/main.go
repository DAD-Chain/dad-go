package main

import (
	"runtime"
	"github.com/dad-go/eventbus/example/testRemoteCrypto/commons"
	"github.com/dad-go/eventbus/eventhub"
	"github.com/dad-go/eventbus/actor"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	props := actor.FromProducer(func() actor.Actor { return &commons.BusynessActor{Datas:make(map[string][]byte)} })
	bActor:=actor.Spawn(props)

	signprops := actor.FromProducer(func() actor.Actor { return &commons.SignActor{} })
	signActor := actor.Spawn(signprops)

	eventhub.GlobalEventHub.Subscribe(commons.SetTOPIC, signActor)
	eventhub.GlobalEventHub.Subscribe(commons.SigTOPIC, signActor)

	vfprops := actor.FromProducer(func() actor.Actor { return &commons.VerifyActor{} })
	vfActor := actor.Spawn(vfprops)

	eventhub.GlobalEventHub.Subscribe(commons.VerifyTOPIC,vfActor)

	bActor.Tell(&commons.RunMsg{})


	for{
		time.Sleep(1 * time.Second)
	}
}
