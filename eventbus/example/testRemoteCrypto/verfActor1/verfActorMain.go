package main

import (
	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/eventbus/example/testRemoteCrypto/commons"
	"runtime"
	"github.com/dad-go/eventbus/remote"
	"github.com/dad-go/common/log"
	"time"
)



func main()  {

	runtime.GOMAXPROCS(runtime.NumCPU() * 1)
	runtime.GC()

	log.Init()
	remote.Start("172.26.127.133:9081")
	vfprops := actor.FromProducer(func() actor.Actor { return &commons.VerifyActor{} })
	actor.SpawnNamed(vfprops, "verify1")


	for{
		time.Sleep(1 * time.Second)
	}
}