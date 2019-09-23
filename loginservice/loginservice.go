package loginservice

import (
	"context"
	"fmt"
	"lbaas/basic/config"
	dbproto "lbaas/dbservice/proto"
	lgproto "lbaas/loginservice/proto"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

type UsrLoginImpl struct {
	dbservice_cli dbproto.DBService
}

func NewUsrLoginImpl() *UsrLoginImpl {
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

	lgClient := dbproto.NewDBService(config.GetCommonVipper().GetString("servicename.dbservicename"), service.Client())
	return &UsrLoginImpl{dbservice_cli: lgClient}
}

func (ul *UsrLoginImpl) Login(ctx context.Context, req *lgproto.LoginReq, rsp *lgproto.LoginRsp) error {
	fmt.Println("receive login msg , name is ", req.Name, "pswd is ", req.Passwd)

	loginrsp, loginer := ul.dbservice_cli.CheckLogin(ctx, &dbproto.CheckLoginReq{Name: req.Name, Passwd: req.Passwd})
	if loginer != nil {
		fmt.Println("login failed ", loginer.Error())
		return loginer
	}
	rsp.Name = loginrsp.Name
	rsp.Errid = loginrsp.Errid
	return nil
}

func (h *UsrLoginImpl) RegisterUsr(ctx context.Context, in *lgproto.RegUsrReq, out *lgproto.RegUsrRsp) error {
	regrsp, regerr := h.dbservice_cli.RegisterUsr(ctx, &dbproto.RegisterUsrReq{Name: in.Name, Passwd: in.Passwd})
	if regerr != nil {
		fmt.Println("register failed ", regerr.Error())
		return regerr
	}
	out.Name = regrsp.Name
	out.Errid = regrsp.Errid
	return nil
}

func Start() {

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
		micro.Name(config.GetCommonVipper().GetString("servicename.loginservicename")),
		micro.Registry(reg),
		micro.Address("0.0.0.0:8000"),
	)
	service.Init()
	loginhd := NewUsrLoginImpl()
	// 注册 Handler
	lgproto.RegisterUsrLoginHandler(service.Server(), loginhd)
	// run server
	if err := service.Run(); err != nil {
		panic(err)
	}

}
