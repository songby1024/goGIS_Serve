package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"serve/initialize"
	"serve/routers"
)

// 在主函数运行之前，读取配置文件
func initConfig() {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("config read fail error = ", err)
	}
}

func main() {
	initConfig()
	initialize.InitGlobal()
	r := gin.Default()
	routers.Router(r)
	r.Run(viper.GetString("server.host") + ":" + viper.GetString("server.post"))
}
