package initialize

import (
	"github.com/fastwego/offiaccount"
	"github.com/go-redis/redis"
)

var Wechat *offiaccount.OffiAccount
var RDB *redis.Client

func Init() {
	// global.PGSQL = initPgsql()
	Wechat = initWechat()
}

func GetRdb() *redis.Client {
	return RDB
}
func GetWechat() *offiaccount.OffiAccount {
	return Wechat
}
