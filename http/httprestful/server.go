package httprestful

import (
	"github.com/dad-go/net/httprestful/common"
	. "github.com/dad-go/net/httprestful/restful"
	. "github.com/dad-go/net/protocol"
)

func StartServer(n Noder) {
	common.SetNode(n)
	func() {
		rest := InitRestServer()
		go rest.Start()
	}()
}

