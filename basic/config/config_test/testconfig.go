package main

import (
	"fmt"
	"lbaas/basic/config"
)

func main() {
	fmt.Println(config.GetCommonVipper().GetString("TimeStamp"))
	fmt.Println(config.GetCommonVipper().GetString("etcdconfig.etcdnode2addr"))
	fmt.Println(config.GetCommonVipper().GetString("etcdconfig.etcdnode3addr"))
	fmt.Println(config.GetCommonVipper().GetString("servicename.loginservicename"))

}
