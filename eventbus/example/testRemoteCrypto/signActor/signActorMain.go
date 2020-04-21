package main

import (
	"github.com/dad-go/eventbus/actor"
	"runtime"
	"github.com/dad-go/eventbus/example/testRemoteCrypto/commons"
	"github.com/dad-go/eventbus/remote"
	"github.com/dad-go/common/log"
	"time"
)

func main()  {

	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	log.Init()
	remote.Start("172.26.127.133:9080")
	signprops := actor.FromProducer(func() actor.Actor { return &commons.SignActor{} })
	actor.SpawnNamed(signprops, "sign")


	for{
		time.Sleep(1 * time.Second)
	}
}