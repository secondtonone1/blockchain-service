package webservice

import (
	"context"
	"fmt"
	lgproto "lbaas/loginservice/proto"
	"net/http"
	"sync"

	constdef "lbaas/basic/common"
	"lbaas/basic/config"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

var etcdservice micro.Service = nil
var lgClient lgproto.UsrLoginService = nil
var loginonce sync.Once

func GetLoginCliInst() lgproto.UsrLoginService {
	loginonce.Do(func() {
		reg := etcdv3.NewRegistry(func(op *registry.Options) {
			op.Addrs = []string{
				config.GetCommonVipper().GetString("etcdconfig.etcdnode1addr"),
				config.GetCommonVipper().GetString("etcdconfig.etcdnode2addr"),
				config.GetCommonVipper().GetString("etcdconfig.etcdnode3addr"),
			}
		})

		// 初始化服务
		etcdservice = micro.NewService(
			micro.Registry(reg),
		)
		etcdservice.Init()
		lgClient = lgproto.NewUsrLoginService(config.GetCommonVipper().GetString("servicename.loginservicename"), etcdservice.Client())
	})
	return lgClient
}

func InitLoginHandler() {

	//消息注册
	webservice.HandleFunc("/", LoginCb)

}

//注册函数写在此处
func LoginCb(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()

		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		if name == "" {
			fmt.Println("name empty!")
			w.Write([]byte(`<html><body><h1>` + "name empty! " + `</h1></body></html>`))
			return
		}

		rsp, err := GetLoginCliInst().Login(context.Background(), &lgproto.LoginReq{Name: name, Passwd: passwd})

		if err != nil {
			fmt.Println("login req failed!")
			http.Error(w, err.Error(), 500)
			return
		}

		if rsp.Errid == constdef.RSP_SUCCESS {
			fmt.Println("login req succss: msg is ")
			fmt.Println("name is ", rsp.Name)
			w.Write([]byte(`<html><body><h1>` + "login req succss: name is " + rsp.Name + `</h1></body></html>`))
			return
		}

		if rsp.Errid == constdef.RSP_PASSWD_ERROR {
			fmt.Println("login failed , passwd error!")
			w.Write([]byte(`<html><body><h1>` + "login failed , passwd error! " + `</h1></body></html>`))
			return
		}

		if rsp.Errid == constdef.RSP_LOGINNAME_NOTFOUND {
			fmt.Println("login usr name not found, name is ", rsp.Name)
			rsp, err := GetLoginCliInst().RegisterUsr(context.Background(), &lgproto.RegUsrReq{Name: "lemon", Passwd: "lemon123"})
			if err != nil {
				fmt.Println("register failed")
				http.Error(w, err.Error(), 500)
				return
			}
			if rsp.Errid == constdef.RSP_USRHASREGED {
				fmt.Println("usr has been reged")
				w.Write([]byte(`<html><body><h1>` + "usr has been reged" + `</h1></body></html>`))
				return
			}

			if rsp.Errid == constdef.RSP_SUCCESS {
				fmt.Println("usr has been registered success")
				fmt.Println("usr name is ", rsp.Name)
				w.Write([]byte(`<html><body><h1>` + "usr has been registered success, name is " + rsp.Name + `</h1></body></html>`))
				return
			}
			return
		}

	} else {
		fmt.Println("get request")
		w.Write([]byte(`<html><body><h1>` + "get request" + `</h1></body></html>`))
	}
}

//其他的回调函数写于此处
//todo
//func ChangePasswd(w http.ResponseWriter, r *http.Request){}
