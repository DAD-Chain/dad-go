package router

import (
	"math/rand"

	"github.com/dad-go/eventbus/actor"
)

type randomGroupRouter struct {
	GroupRouter
}

type randomPoolRouter struct {
	PoolRouter
}

type randomRouterState struct {
	routees *actor.PIDSet
	values  []actor.PID
}

func (state *randomRouterState) SetRoutees(routees *actor.PIDSet) {
	state.routees = routees
	state.values = routees.Values()
}

func (state *randomRouterState) GetRoutees() *actor.PIDSet {
	return state.routees
}

func (state *randomRouterState) RouteMessage(message interface{}) {
	l := len(state.values)
	r := rand.Intn(l)
	pid := state.values[r]
	pid.Tell(message)
}

func NewRandomPool(size int) *actor.Props {
	return actor.FromSpawnFunc(spawner(&randomPoolRouter{PoolRouter{PoolSize: size}}))
}

func NewRandomGroup(routees ...*actor.PID) *actor.Props {
	return actor.FromSpawnFunc(spawner(&randomGroupRouter{GroupRouter{Routees: actor.NewPIDSet(routees...)}}))
}

func (config *randomPoolRouter) CreateRouterState() Interface {
	return &randomRouterState{}
}

func (config *randomGroupRouter) CreateRouterState() Interface {
	return &randomRouterState{}
}
