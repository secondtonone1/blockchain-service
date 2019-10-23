package webservice

import (
	"context"
	"encoding/json"
	"fmt"
	lgproto "lbaas/loginservice/proto"
	"net/http"
	"sync"

	constdef "lbaas/basic/common"
	"lbaas/basic/config"

	"time"

	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/registry/etcdv3"
)

var etcdservice micro.Service = nil
var lgClient lgproto.UsrLoginService = nil
var loginonce sync.Once

type Token struct {
	Token string `json:"token"`
}

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
	webservice.HandleFunc("/register", RegisterUsr)
	webservice.HandleFunc("/changepasswd", ChangePasswd)
	webservice.HandleFunc("/cfchangepasswd", CFChangePasswd)
	webservice.Handle("/tokenchangepasswd", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(TokenChangePasswd)),
	))

}

//注册函数写在此处
func LoginCb(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("receive post request")
		r.ParseForm()
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		if name == "" {
			fmt.Println("name empty!")
			w.Write([]byte(`<html><body><h1>` + "name empty! " + `</h1></body></html>`))
			return
		}

		if passwd == "" {
			fmt.Println("passwd empty")
			w.Write([]byte(`<html><body><h1>` + "passwd empty! " + `</h1></body></html>`))
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
//修改密码
func ChangePasswd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("receive post request")
		r.ParseForm()
		name := r.Form.Get("name")
		emailaddr := r.Form.Get("email")
		if name == "" {
			fmt.Println("name empty!")
			w.Write([]byte(`<html><body><h1>` + "name empty! " + `</h1></body></html>`))
			return
		}

		if emailaddr == "" {
			fmt.Println("email empty")
			w.Write([]byte(`<html><body><h1>` + "passwd empty! " + `</h1></body></html>`))
			return
		}

		rsp, err := GetLoginCliInst().ChangePasswd(context.Background(), &lgproto.ChangewdReq{Name: name, Email: emailaddr})

		if err != nil {
			fmt.Println("changepasswd req failed!")
			http.Error(w, err.Error(), 500)
			return
		}

		if rsp.Errid == constdef.RSP_SUCCESS {
			fmt.Println("changepasswd req succss: msg is ")
			fmt.Println("name is ", rsp.Name)
			tokenstr, generr := GenerateJwt()
			if generr != nil {
				fmt.Println("generate token error!")
				return
			}
			fmt.Println("generate token success, tokenstr is : ", tokenstr)
			tk := Token{tokenstr}
			json, errmas := json.Marshal(tk)
			if errmas != nil {
				fmt.Println("json marsh error!")
				return
			}

			w.Write([]byte(`<html><body><h1>` + "changepasswd req succss: name is " + rsp.Name + `</h1></body></html>`))
			w.Write(json)
			return
		}

		if rsp.Errid == constdef.RSP_USERNAME_ERROR {
			fmt.Println("changepasswd failed , RSP_USERNAME_ERROR")
			w.Write([]byte(`<html><body><h1>` + "changepasswd failed , user has been registered! " + `</h1></body></html>`))
			return
		}

		if rsp.Errid == constdef.RSP_EMAIL_ERROR {
			fmt.Println("registered failed , RSP_EMAIL_ERROR!")
			w.Write([]byte(`<html><body><h1>` + "changepasswd failed , RSP_EMAIL_ERROR " + `</h1></body></html>`))
			return
		}

	} else {
		fmt.Println("get request")
		w.Write([]byte(`<html><body><h1>` + "get request" + `</h1></body></html>`))
	}
}

//注册账号
func RegisterUsr(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println("receive post request")
		r.ParseForm()
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		emailaddr := r.Form.Get("email")
		if name == "" {
			fmt.Println("name empty!")
			w.Write([]byte(`<html><body><h1>` + "name empty! " + `</h1></body></html>`))
			return
		}

		if passwd == "" {
			fmt.Println("passwd empty")
			w.Write([]byte(`<html><body><h1>` + "passwd empty! " + `</h1></body></html>`))
			return
		}

		if emailaddr == "" {
			fmt.Println("email empty")
			w.Write([]byte(`<html><body><h1>` + "passwd empty! " + `</h1></body></html>`))
			return
		}

		rsp, err := GetLoginCliInst().RegisterUsr(context.Background(), &lgproto.RegUsrReq{Name: name, Passwd: passwd, Email: emailaddr})

		if err != nil {
			fmt.Println("register req failed!")
			http.Error(w, err.Error(), 500)
			return
		}

		if rsp.Errid == constdef.RSP_SUCCESS {
			fmt.Println("register req succss: msg is ")
			fmt.Println("name is ", rsp.Name)
			w.Write([]byte(`<html><body><h1>` + "register req succss: name is " + rsp.Name + `</h1></body></html>`))
			return
		}

		if rsp.Errid == constdef.RSP_USRHASREGED {
			fmt.Println("register failed , user has been registered")
			w.Write([]byte(`<html><body><h1>` + "register failed , user has been registered! " + `</h1></body></html>`))
			return
		}

		if rsp.Errid == constdef.RSP_USRREG_FAILED {
			fmt.Println("registered failed , unkown error!")
			w.Write([]byte(`<html><body><h1>` + "registered failed , unkown error! " + `</h1></body></html>`))
			return
		}

	} else {
		fmt.Println("get request")
		w.Write([]byte(`<html><body><h1>` + "get request" + `</h1></body></html>`))
	}
}

//确认更改密码
func CFChangePasswd(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		fmt.Println("receive post request")
		r.ParseForm()
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		emailaddr := r.Form.Get("email")
		if name == "" {
			fmt.Println("name empty!")
			w.Write([]byte(`<html><body><h1>` + "name empty! " + `</h1></body></html>`))
			return
		}

		if passwd == "" {
			fmt.Println("passwd empty")
			w.Write([]byte(`<html><body><h1>` + "passwd empty! " + `</h1></body></html>`))
			return
		}

		if emailaddr == "" {
			fmt.Println("email empty")
			w.Write([]byte(`<html><body><h1>` + "passwd empty! " + `</h1></body></html>`))
			return
		}

		rsp, err := GetLoginCliInst().ChangePasswdConfirm(context.Background(), &lgproto.ChangewdCfmReq{Name: name, Passwd: passwd, Email: emailaddr})

		if err != nil {
			fmt.Println("confirm req failed!")
			http.Error(w, err.Error(), 500)
			return
		}

		if rsp.Errid == constdef.RSP_SUCCESS {
			fmt.Println("confirm req succss: msg is ")
			fmt.Println("name is ", rsp.Name)
			w.Write([]byte(`<html><body><h1>` + "confirm req succss: name is " + rsp.Name + `</h1></body></html>`))
			return
		}

		if rsp.Errid == constdef.RSP_USERNAME_ERROR {
			fmt.Println("confirm failed , RSP_USERNAME_ERROR")
			w.Write([]byte(`<html><body><h1>` + "confirm failed , RSP_USERNAME_ERROR! " + `</h1></body></html>`))
			return
		}

	} else {
		fmt.Println("get request")
		w.Write([]byte(`<html><body><h1>` + "get request" + `</h1></body></html>`))
	}
}

func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(constdef.Tokenkey), nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

func GenerateJwt() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString([]byte(constdef.Tokenkey))
	if err != nil {
		fmt.Println("Error while signing the token")
		return "", err
	}
	return tokenString, nil
}

func TokenChangePasswd(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		fmt.Println("receive token change password request")
		r.ParseForm()
		name := r.Form.Get("name")
		passwd := r.Form.Get("passwd")
		emailaddr := r.Form.Get("email")
		if name == "" {
			fmt.Println("name empty!")
			w.Write([]byte(`<html><body><h1>` + "name empty! " + `</h1></body></html>`))
			return
		}

		if passwd == "" {
			fmt.Println("passwd empty")
			w.Write([]byte(`<html><body><h1>` + "passwd empty! " + `</h1></body></html>`))
			return
		}

		if emailaddr == "" {
			fmt.Println("email empty")
			w.Write([]byte(`<html><body><h1>` + "passwd empty! " + `</h1></body></html>`))
			return
		}

		rsp, err := GetLoginCliInst().ChangePasswdConfirm(context.Background(), &lgproto.ChangewdCfmReq{Name: name, Passwd: passwd, Email: emailaddr})

		if err != nil {
			fmt.Println("confirm req failed!")
			http.Error(w, err.Error(), 500)
			return
		}

		if rsp.Errid == constdef.RSP_SUCCESS {
			fmt.Println("confirm req succss: msg is ")
			fmt.Println("name is ", rsp.Name)
			w.Write([]byte(`<html><body><h1>` + "confirm req succss: name is " + rsp.Name + `</h1></body></html>`))
			return
		}

		if rsp.Errid == constdef.RSP_USERNAME_ERROR {
			fmt.Println("confirm failed , RSP_USERNAME_ERROR")
			w.Write([]byte(`<html><body><h1>` + "confirm failed , RSP_USERNAME_ERROR! " + `</h1></body></html>`))
			return
		}

	} else {
		fmt.Println("get request")
		w.Write([]byte(`<html><body><h1>` + "get request" + `</h1></body></html>`))
	}
}
