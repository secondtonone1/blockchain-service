package loginservice

import (
	"blockchain-service/basic/config"
	dbproto "blockchain-service/dbservice/proto"
	lgproto "blockchain-service/loginservice/proto"

	"context"
	"fmt"
	"time"

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
	fmt.Println("receive login msg , name is ", req.Name, "pswd is ", req.Passwd, time.Now())

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
	regrsp, regerr := h.dbservice_cli.RegisterUsr(ctx, &dbproto.RegisterUsrReq{Name: in.Name, Passwd: in.Passwd, Email: in.Email})
	if regerr != nil {
		fmt.Println("register failed ", regerr.Error())
		return regerr
	}
	out.Name = regrsp.Name
	out.Errid = regrsp.Errid
	return nil
}

func (h *UsrLoginImpl) ChangePasswd(ctx context.Context, in *lgproto.ChangewdReq, out *lgproto.ChangewdRsp) error {
	checkrsp, ckerr := h.dbservice_cli.CheckUsrEmail(ctx, &dbproto.CheckUsrEmailReq{Name: in.Name, Email: in.Email})
	if ckerr != nil {
		fmt.Println("check failed", ckerr.Error())
		return ckerr
	}
	out.Name = checkrsp.Name
	out.Email = checkrsp.Email
	out.Errid = checkrsp.Errid

	//此处应该调用第三方发送邮件给认证人邮箱
	return nil
}

func (h *UsrLoginImpl) ChangePasswdConfirm(ctx context.Context, in *lgproto.ChangewdCfmReq, out *lgproto.ChangewdCfmRsp) error {
	confirmrsp, conerr := h.dbservice_cli.ChangePwd(ctx, &dbproto.ResetPwdReq{Name: in.Name, Passwd: in.Passwd})
	if conerr != nil {
		fmt.Println("confirm failed", conerr.Error())
		return conerr
	}
	out.Name = confirmrsp.Name
	out.Email = confirmrsp.Email
	out.Errid = confirmrsp.Errid
	return nil
}

func Start() {

	// 我这里用的etcd 做为服务发现，如果使用consul可以去掉
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{
			//"http://47.105.111.1:2379", "http://47.105.111.1:2379",
			config.GetCommonVipper().GetString("etcdconfig.etcdnode1addr"),
			config.GetCommonVipper().GetString("etcdconfig.etcdnode2addr"),
			config.GetCommonVipper().GetString("etcdconfig.etcdnode3addr"),
		}
	})

	srvName := config.GetCommonVipper().GetString("servicename.loginservicename")
	//	t, io, err := tracer.NewTracer(srvName, "localhost:6831")

	// 初始化服务
	service := micro.NewService(
		micro.Name(srvName),
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
