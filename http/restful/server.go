package restful

import (
	. "github.com/dad-go/http/restful/restful"
)

func StartServer() {
	func() {
		rt := InitRestServer()
		go rt.Start()
	}()
}

