package restful

import (
	. "github.com/dad-go/http/base/rest"
	. "github.com/dad-go/http/restful/restful"
	. "github.com/dad-go/net/protocol"
)

func StartServer(n Noder) {
	SetNode(n)
	func() {
		rt := InitRestServer()
		go rt.Start()
	}()
}

