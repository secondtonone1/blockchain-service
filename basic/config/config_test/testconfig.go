package main

import (
	"blockchain-service/basic/config"
	"fmt"
)

func main() {
	fmt.Println(config.GetCommonVipper().GetString("TimeStamp"))
	fmt.Println(config.GetCommonVipper().GetString("etcdconfig.etcdnode2addr"))
	fmt.Println(config.GetCommonVipper().GetString("etcdconfig.etcdnode3addr"))
	fmt.Println(config.GetCommonVipper().GetString("servicename.loginservicename"))

}
