package actor

import (
	"reflect"
	"testing"
	"github.com/dad-go/common/log"

	"github.com/stretchr/testify/assert"
)

type ShortLivingActor struct {
}

func (self *ShortLivingActor) Receive(ctx Context) {

}

func TestStopFuture(t *testing.T) {
	log.Debug("hello world")

	ID := "UniqueID"
	{
		props := FromProducer(func() Actor { return &ShortLivingActor{} })
		a, _ := SpawnNamed(props, ID)

		fut := a.StopFuture()

		res, errR := fut.Result()
		if errR != nil {
			assert.Fail(t, "Failed to wait stop actor %s", errR)
			return
		}

		_, ok := res.(*Terminated)
		if !ok {
			assert.Fail(t, "Cannot cast %s", reflect.TypeOf(res))
			return
		}

		_, found := ProcessRegistry.Get(a)
		assert.False(t, found)
	}
}
