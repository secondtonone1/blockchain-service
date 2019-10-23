package dbservice

import (
	constdef "blockchain-service/basic/common"
	"blockchain-service/basic/config"
	dbproto "blockchain-service/dbservice/proto"
	"context"
	"fmt"
	"time"

	"blockchain-service/dbservice/dbtable"

	"blockchain-service/dbservice/dbmanager"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

type UsrDBImpl struct {
	usrdata map[string]*dbtable.UsrBase
}

func (cl *UsrDBImpl) CheckLogin(ctx context.Context, in *dbproto.CheckLoginReq, out *dbproto.CheckLoginRsp) error {
	fmt.Println("servicedb receive check login msg", time.Now())
	value, ok := cl.usrdata[in.Name]
	if ok {
		out.Name = in.Name
		out.Errid = constdef.RSP_SUCCESS

		if value.Passwd != in.Passwd {
			fmt.Println("usr passwd error!")
			out.Errid = constdef.RSP_PASSWD_ERROR
			return nil
		}

		fmt.Println("user has been found, name is ", value.Name)
		return nil
	}

	temp := &dbtable.UsrBase{Name: in.Name}
	findres := dbtable.UsrBaseTblFind(temp)
	if findres != nil {
		out.Name = in.Name
		out.Errid = constdef.RSP_LOGINNAME_NOTFOUND
		fmt.Println("usr not found in db, name is ", in.Name)
		return nil
	}

	cl.usrdata[temp.Name] = temp
	out.Name = in.Name
	out.Errid = constdef.RSP_SUCCESS
	fmt.Println("usr has been found in DB, name is ", temp.Name)
	fmt.Println("usr has been found in DB, passwd is ", temp.Passwd)

	return nil
}

func (ru *UsrDBImpl) RegisterUsr(ctx context.Context, in *dbproto.RegisterUsrReq, out *dbproto.RegisterUsrRsp) error {
	//constdef.RSP_USRHASREGED
	value, ok := ru.usrdata[in.Name]
	if ok {
		fmt.Printf("usr %s has been reged \n", in.Name)
		out.Name = value.Name
		out.Errid = constdef.RSP_USRHASREGED
		return nil
	}

	usrbase := &dbtable.UsrBase{Name: in.Name}
	finderr := dbtable.UsrBaseTblFind(usrbase)
	if finderr == nil {
		ru.usrdata[usrbase.Name] = usrbase
		fmt.Printf("usr %s has found \n", in.Name)
		out.Name = in.Name
		out.Errid = constdef.RSP_USRHASREGED
		return nil
	}

	usrbase.Name = in.Name
	usrbase.Passwd = in.Passwd
	usrbase.Email = in.Email
	inserterr := dbtable.UsrBaseTblInsert(usrbase)
	if inserterr != nil {
		fmt.Printf("usr %s insert failed \n", in.Name)
		out.Name = in.Name
		out.Errid = constdef.RSP_USRREG_FAILED
		return inserterr
	}
	fmt.Println("usrbase usr id is ", usrbase.ID)
	ru.usrdata[usrbase.Name] = usrbase
	out.Name = usrbase.Name
	out.Errid = constdef.RSP_SUCCESS
	out.Email = usrbase.Email
	return nil
}

func (h *UsrDBImpl) ChangePwd(ctx context.Context, in *dbproto.ResetPwdReq, out *dbproto.ResetPwdRsp) error {
	value, ok := h.usrdata[in.Name]
	if !ok {
		temp := &dbtable.UsrBase{Name: in.Name}
		findres := dbtable.UsrBaseTblFind(temp)
		if findres != nil {
			out.Name = in.Name
			out.Errid = constdef.RSP_USERNAME_ERROR
			fmt.Println("usr not found in db, name is ", in.Name)
			return nil
		}
		h.usrdata[temp.Name] = temp
		value = temp
	}

	value.Passwd = in.Passwd
	upres := dbtable.UsrBaseTblUpdate(&dbtable.UsrBase{ID: value.ID, Name: value.Name}, value)
	if upres != nil {
		out.Name = in.Name
		out.Email = value.Email
		out.Errid = constdef.RSP_UPDATE_ERROR
		return nil
	}
	out.Name = in.Name
	out.Errid = constdef.RSP_SUCCESS
	out.Email = value.Email
	return nil
}

func (h *UsrDBImpl) CheckUsrEmail(ctx context.Context, in *dbproto.CheckUsrEmailReq, out *dbproto.CheckUsrEmailRsp) error {
	value, ok := h.usrdata[in.Name]

	if !ok {
		temp := &dbtable.UsrBase{Name: in.Name}
		findres := dbtable.UsrBaseTblFind(temp)
		if findres != nil {
			out.Name = in.Name
			out.Email = in.Email
			out.Errid = constdef.RSP_USERNAME_ERROR
			fmt.Println("usr not found in db, name is ", in.Name)
			return nil
		}

		h.usrdata[temp.Name] = temp
		value = temp
	}

	if value.Email != in.Email {
		fmt.Println("usr email : ", in.Email, " not found!")
		out.Name = value.Name
		out.Email = value.Email
		out.Errid = constdef.RSP_EMAIL_ERROR
		return nil
	}

	out.Name = value.Name
	out.Email = value.Email
	out.Errid = constdef.RSP_SUCCESS
	fmt.Println("db check usr and email success!")
	return nil

}

func newUsrDBImpl() *UsrDBImpl {
	udb := new(UsrDBImpl)
	udb.usrdata = make(map[string]*dbtable.UsrBase)
	initTbls()
	return udb
}

func initTbls() {
	dbtable.UsrBaseTblCreate()
	dbtable.ProjectTblCreate()
}

func Start() {

	dbmgr := dbmanager.GetDBMgrInst()
	if dbmgr == nil {
		fmt.Println("dbmgr init failed!")
	}
	defer dbmgr.ReleaseDB()

	// 我这里用的etcd 做为服务发现，如果使用consul可以去掉
	reg := etcdv3.NewRegistry(func(op *registry.Options) {
		op.Addrs = []string{

			config.GetCommonVipper().GetString("etcdconfig.etcdnode1addr"),
			config.GetCommonVipper().GetString("etcdconfig.etcdnode2addr"),
			config.GetCommonVipper().GetString("etcdconfig.etcdnode3addr"),
			//"http://47.105.111.1:2379", "http://47.105.111.1:2379",
		}
	})

	// 初始化服务
	service := micro.NewService(
		micro.Name(config.GetCommonVipper().GetString("servicename.dbservicename")),
		micro.Registry(reg),
	)
	service.Init()
	dbhandler := newUsrDBImpl()
	// 注册 Handler
	dbproto.RegisterDBServiceHandler(service.Server(), dbhandler)

	// run server
	if err := service.Run(); err != nil {
		panic(err)
	}

}
