package main

import (
	constdef "blockchain-service/basic/common"
	"blockchain-service/basic/config"
	lgproto "blockchain-service/loginservice/proto"
	"context"
	"fmt"

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
	rsp, err := lgClient.Login(context.Background(), &lgproto.LoginReq{Name: "Zack", Passwd: "123"})

	if err != nil {
		fmt.Println("login req failed!")
		return
	}
	fmt.Println("login req succss: msg is ")
	fmt.Println("erroid is ", rsp.Errid)
	if rsp.Errid == constdef.RSP_SUCCESS {
		fmt.Println("name is ", rsp.Name)
	}

}
