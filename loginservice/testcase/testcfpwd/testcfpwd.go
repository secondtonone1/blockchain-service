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
	rsp, err := lgClient.ChangePasswdConfirm(context.Background(), &lgproto.ChangewdCfmReq{Name: "lemon", Passwd: "lemon1234"})

	if err != nil {
		fmt.Println("change passwd failed!")
		fmt.Println(err.Error())
		return
	}

	if rsp.Errid == constdef.RSP_USERNAME_ERROR {
		fmt.Println("user name error")
		return
	}

	if rsp.Errid == constdef.RSP_EMAIL_ERROR {
		fmt.Println("email  error")
		return
	}

	fmt.Println("handle change passwd req success!")
	return

}
