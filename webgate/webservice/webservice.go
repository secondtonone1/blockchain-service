package webservice

import (
	"sync"

	"github.com/micro/go-micro/web"
)

var webservice web.Service
var webonce sync.Once

func GetWebServiceInst() web.Service {
	webonce.Do(func() {
		webservice = web.NewService(
			web.Name("go.micro.web.greeter"),
			web.Address("127.0.0.1:12222"),
		)
		//填写一系列消息回调，此处为登陆模式
		InitLoginHandler()
	})

	return webservice
}
