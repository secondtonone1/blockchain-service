package main

import (
	"context"
	"fmt"
	constdef "lbaas/basic/common"
	lgproto "lbaas/loginservice/proto"

	"lbaas/basic/config"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

func main() {
	// 我这里用的etcd 做为服务发现，如果使用consul可以去掉

	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			config.GetCommonVipper().GetString("etcdconfig.etcdnode1addr"),
			config.GetCommonVipper().GetString("etcdconfig.etcdnode2addr"),
			config.GetCommonVipper().GetString("etcdconfig.etcdnode3addr"),
		}
	})

	// 初始化服务
	service := micro.NewService(
		micro.Registry(reg),
	)
	service.Init()

	lgClient := lgproto.NewUsrLoginService(config.GetCommonVipper().GetString("servicename.loginservicename"), service.Client())
	rsp, err := lgClient.Login(context.Background(), &lgproto.LoginReq{Name: config.GetCommonVipper().GetString("dbconfig.dbuser"),
		Passwd: config.GetCommonVipper().GetString("dbconfig.dbpswd")})

	if err != nil {
		fmt.Println("login req failed!", err)
		return
	}

	if rsp.Errid == constdef.RSP_SUCCESS {
		fmt.Println("login req succss: msg is ")
		fmt.Println("name is ", rsp.Name)
		return
	}

	if rsp.Errid == constdef.RSP_PASSWD_ERROR {
		fmt.Println("login failed , passwd error!")
		return
	}

	if rsp.Errid == constdef.RSP_LOGINNAME_NOTFOUND {
		fmt.Println("login usr name not found, name is ", rsp.Name)
		rsp, err := lgClient.RegisterUsr(context.Background(), &lgproto.RegUsrReq{Name: "lemon", Passwd: "lemon123", Email: "lemon@163.com"})
		if err != nil {
			fmt.Println("register failed")
			return
		}
		if rsp.Errid == constdef.RSP_USRHASREGED {
			fmt.Println("usr has been reged")
			return
		}

		if rsp.Errid == constdef.RSP_SUCCESS {
			fmt.Println("usr has been registered success")
			fmt.Println("usr name is ", rsp.Name)
			return
		}
		return
	}

}
