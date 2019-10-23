##命令行测试web功能
暴露web后台
go run ./web.go
启动web 服务
micro --server_address 127.0.0.1:8082 web
命令行发送登陆请求两种方式
方式1 请求自己web server
curl -d "name=lemon&passwd=lemon123" "http://127.0.0.1:12222"
方式2 请求micro web服务
curl -d "name=lemon&passwd=lemon123" "http://127.0.0.1:8082/greeter"

##代码测试
仅测试自己的webserver
1 在config-etcd-srv文件夹中启动etcd服务
./etcd_start.sh
2 在main文件夹启动loginservice
go run ./main.go loginservice
3 在main文件夹启动dbservice
go run ./main.go dbservice
4 在webgate文件夹启动web服务
go run ./web.go
5 进入webgate testcase 文件夹
进入 changepasswdtest测试申请修改密码
服务器会生成token给客户端
进入 tokenchangepasswd文件夹
测试确认修改，将tokenstr发送给服务器
服务器确认无误后修改。

