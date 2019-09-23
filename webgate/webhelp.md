暴露web后台
go run ./web.go
启动web 服务
micro --server_address 127.0.0.1:8082 web

curl -d "name=lemon&passwd=lemon123" "http://127.0.0.1:12222/greeter"
curl -d "name=lemon&passwd=lemon123" "http://127.0.0.1:8082/greeter"