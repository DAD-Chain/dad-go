package actor

import (
	"github.com/dad-go/eventbus/actor"
	//"github.com/dad-go/net/message"
)

var ConsensusPid *actor.PID

//func PushConsensus(cons *message.ConsensusPayload){
//	ConsensusPid.Tell(cons)
//}

func SetConsensusPid(conPid * actor.PID){
	ConsensusPid = conPid
}