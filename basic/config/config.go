package config

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var commonvipper *viper.Viper = nil
var commonce sync.Once

func GetCommonVipper() *viper.Viper {
	commonce.Do(func() {
		// dir, err := filepath.Abs(filepath.Dir("../basic/config/common.yaml"))
		// if err != nil {
		// 	return
		// }
		//fmt.Println(dir)
		commonvipper = viper.New()
		commonvipper.SetConfigType("yaml")
		commonvipper.SetConfigName("common")
		commonvipper.AddConfigPath("$GOPATH/src/lbaas/basic/config")

		if err := commonvipper.ReadInConfig(); err != nil {
			fmt.Printf("err:%s\n", err)
		}

		commonvipper.WatchConfig()
		commonvipper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("Config file changed:", e.Name)
		})
	})

	return commonvipper
}
