package signtest

import (
	"github.com/dad-go/eventbus/actor"
	"fmt"
	"github.com/dad-go/crypto"
	"bytes"

	"runtime"
	"strconv"
	"time"
	"github.com/dad-go/eventbus/eventhub"
)

//var signActor ,vfActor *actor.PID
const loop = 10
func init(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	signprops := actor.FromProducer(func() actor.Actor { return &SignActor{} })
	vfprops := actor.FromProducer(func() actor.Actor { return &VerifyActor{} })

	//signActor,_ = actor.SpawnNamed(signprops,"sig1")
	signActor := actor.Spawn(signprops)
	signActor2 := actor.Spawn(signprops)
	signActor3 := actor.Spawn(signprops)
	signActor4 := actor.Spawn(signprops)
	signActor5 := actor.Spawn(signprops)
	//vfActor ,_= actor.SpawnNamed(vfprops,"vf1")
	vfActor := actor.Spawn(vfprops)
	vfActor2 := actor.Spawn(vfprops)
	vfActor3 := actor.Spawn(vfprops)
	vfActor4 := actor.Spawn(vfprops)
	vfActor5 := actor.Spawn(vfprops)

	sigTOPIC := "SIGTOPIC"
	verifyTOPIC := "VERIFYTOPIC"
	setTOPIC := "SETTOPIC"

	eventhub.GlobalEventHub.Subscribe(sigTOPIC,signActor)
	eventhub.GlobalEventHub.Subscribe(sigTOPIC,signActor2)
	eventhub.GlobalEventHub.Subscribe(sigTOPIC,signActor3)
	eventhub.GlobalEventHub.Subscribe(sigTOPIC,signActor4)
	eventhub.GlobalEventHub.Subscribe(sigTOPIC,signActor5)

	eventhub.GlobalEventHub.Subscribe(setTOPIC,signActor)
	eventhub.GlobalEventHub.Subscribe(setTOPIC,signActor2)
	eventhub.GlobalEventHub.Subscribe(setTOPIC,signActor3)
	eventhub.GlobalEventHub.Subscribe(setTOPIC,signActor4)
	eventhub.GlobalEventHub.Subscribe(setTOPIC,signActor5)


	eventhub.GlobalEventHub.Subscribe(verifyTOPIC,vfActor)
	eventhub.GlobalEventHub.Subscribe(verifyTOPIC,vfActor2)
	eventhub.GlobalEventHub.Subscribe(verifyTOPIC,vfActor3)
	eventhub.GlobalEventHub.Subscribe(verifyTOPIC,vfActor4)
	eventhub.GlobalEventHub.Subscribe(verifyTOPIC,vfActor5)

}
type RunMsg struct{

}

type BusynessActor struct{
	Datas map[string][]byte
	privatekey []byte
	pubkey crypto.PubKey
	start int64
}


func (s *BusynessActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		fmt.Println("Started, initialize actor here")
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about shut down")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about restart")
	case *RunMsg:
		crypto.SetAlg("SM2")
		bb := bytes.NewBuffer([]byte("s"))
		for i := 0; i < 10000; i++ {
			bb.WriteString("1234567890abcdefghijklmnopqrstuvwxyz")
		}


		privKey, pubkey, _ := crypto.GenKeyPair()

		s.privatekey = privKey
		s.pubkey = pubkey

		setPrivMsg := &SetPrivKey{PrivKey: privKey}

		setEvent := &eventhub.Event{Topic:"SETTOPIC",Publisher:context.Self(),Message:setPrivMsg,Policy:eventhub.PUBLISH_POLICY_ALL}

		eventhub.GlobalEventHub.Publish(setEvent)

		//signActor.Tell(setPrivMsg)

		s.start = time.Now().UnixNano()
		for i := 0; i < loop; i++ {

			idx := strconv.Itoa(i)
			bb.WriteString(idx)
			data := bb.Bytes()
			sigMsg := &SignRequest{Seq: idx, Data: data}
			s.Datas[idx] = data
			sigEvent := &eventhub.Event{Topic:"SIGTOPIC",Publisher:context.Self(),Message:sigMsg,Policy:eventhub.PUBLISH_POLICY_ROUNDROBIN}
			eventhub.GlobalEventHub.Publish(sigEvent)
			//signActor.Request(sigMsg,context.Self())
		}

	case *SignResponse:
		seq := msg.Seq
		sig := msg.Signature
		//fmt.Printf("seq:%s,sig:%v\n",seq,sig)

		vfr:= &VerifyRequest{Signature:sig,Data:s.Datas[seq],PublicKey:s.pubkey,Seq:seq}
		//vfActor.Request(vfr,context.Self())

		vrfEvent := &eventhub.Event{Topic:"VERIFYTOPIC",Publisher:context.Self(),Message:vfr,Policy:eventhub.PUBLISH_POLICY_ROUNDROBIN}

		eventhub.GlobalEventHub.Publish(vrfEvent)
/*		i ,_:= strconv.Atoi(seq)
		if  i == loop-1{
			spend:= (time.Now().UnixNano() - s.start)/1000000
			fmt.Printf("signature spend %d ms\n",int(spend))
		}*/


	case *VerifyResponse:
		seq := msg.Seq
		result:= msg.Result
		errmsg := msg.ErrorMsg

		if !result{
			fmt.Printf("seq:%s faild pass,err:%s\n",seq,errmsg)
		}else{
			fmt.Printf("seq:%s passed verify\n",seq)
		}

/*		i ,_:= strconv.Atoi(seq)
		if  i == loop-1{
			spend:= (time.Now().UnixNano() - s.start)/1000000
			fmt.Printf("verify spend %d ms",int(spend))
		}*/


	default:
		fmt.Printf("unknown msg %v\n", msg)
	}
}

