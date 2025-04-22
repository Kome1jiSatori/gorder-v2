package main

import (
	"github.com/Kome1jiSatori/gorder-v2/common/config"
	"github.com/spf13/viper"
	"log"
)

func init() {
	if err := config.NewViperConfig(); err != nil {
		log.Fatal(err)
	}
}

// 开启服务
func main() {
	log.Printf("%v", viper.Get("order"))
}
