package remote

import (
	"errors"
	"fmt"
	"time"

	"github.com/dad-go/eventbus/actor"
	"github.com/dad-go/common/log"
)

var (
	nameLookup   = make(map[string]actor.Props)
	activatorPid *actor.PID
)

func spawnActivatorActor() {
	props := actor.FromProducer(newActivatorActor()).WithGuardian(actor.RestartingSupervisorStrategy())
	activatorPid, _ = actor.SpawnNamed(props, "activator")
}

func stopActivatorActor() {
	activatorPid.GracefulStop()
}

//Register a known actor props by name
func Register(kind string, props *actor.Props) {
	nameLookup[kind] = *props
}

//GetKnownKinds returns a slice of known actor "kinds"
func GetKnownKinds() []string {
	keys := make([]string, 0, len(nameLookup))
	for k := range nameLookup {
		keys = append(keys, k)
	}
	return keys
}

type activator struct {
}

//ErrActivatorUnavailable : this error will not panic the Activator.
//It simply tells Partition this Activator is not available
//Partition will then find next available Activator to spawn
var ErrActivatorUnavailable = &ActivatorError{ResponseStatusCodeUNAVAILABLE.ToInt32(), true}

type ActivatorError struct {
	Code       int32
	DoNotPanic bool
}

func (e *ActivatorError) Error() string {
	return fmt.Sprint(e.Code)
}

//ActivatorForAddress returns a PID for the activator at the given address
func ActivatorForAddress(address string) *actor.PID {
	pid := actor.NewPID(address, "activator")
	return pid
}

//SpawnFuture spawns a remote actor and returns a Future that completes once the actor is started
func SpawnFuture(address, name, kind string, timeout time.Duration) *actor.Future {
	activator := ActivatorForAddress(address)
	f := activator.RequestFuture(&ActorPidRequest{
		Name: name,
		Kind: kind,
	}, timeout)
	return f
}

//Spawn spawns a remote actor of a given type at a given address
func Spawn(address, kind string, timeout time.Duration) (*ActorPidResponse, error) {
	return SpawnNamed(address, "", kind, timeout)
}

//SpawnNamed spawns a named remote actor of a given type at a given address
func SpawnNamed(address, name, kind string, timeout time.Duration) (*ActorPidResponse, error) {
	activator := ActivatorForAddress(address)
	res, err := activator.RequestFuture(&ActorPidRequest{
		Name: name,
		Kind: kind,
	}, timeout).Result()
	if err != nil {
		return nil, err
	}
	switch msg := res.(type) {
	case *ActorPidResponse:
		return msg, nil
	default:
		return nil, errors.New("remote: Unknown response when remote activating")
	}
}

func newActivatorActor() actor.Producer {
	return func() actor.Actor {
		return &activator{}
	}
}

func (*activator) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *actor.Started:
		log.Debug("Started Activator")
	case *ActorPidRequest:
		props := nameLookup[msg.Kind]
		name := msg.Name

		//unnamed actor, assign auto ID
		if name == "" {
			name = actor.ProcessRegistry.NextId()
		}

		pid, err := actor.SpawnNamed(&props, "Remote$"+name)

		if err == nil {
			response := &ActorPidResponse{Pid: pid}
			context.Respond(response)
		} else if err == actor.ErrNameExists {
			response := &ActorPidResponse{
				Pid:        pid,
				StatusCode: ResponseStatusCodePROCESSNAMEALREADYEXIST.ToInt32(),
			}
			context.Respond(response)
		} else if aErr, ok := err.(*ActivatorError); ok {
			response := &ActorPidResponse{
				StatusCode: aErr.Code,
			}
			context.Respond(response)
			if !aErr.DoNotPanic {
				panic(err)
			}
		} else {
			response := &ActorPidResponse{
				StatusCode: ResponseStatusCodeERROR.ToInt32(),
			}
			context.Respond(response)
			panic(err)
		}
	case actor.SystemMessage, actor.AutoReceiveMessage:
		//ignore
	default:
		log.Error("Activator received unknown message!")
	}
}
