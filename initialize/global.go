package initialize

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"serve/common"
	"serve/common/global"
)

func InitLogger() {
	//初始化日志
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("日志初始化失败", err.Error())
	}
	//使用全局logger
	zap.ReplaceGlobals(logger)
}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func InitGlobal() {
	InitLogger()
	global.DB = common.InitDB()
	global.RDB = InitRedis()
	global.PGSQL = common.InitPgsql()
}
