package common

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"serve/model"
)

// InitDB 连接mysql数据库
func InitDB() *gorm.DB {
	host := viper.GetString("mysql.host")
	post := viper.GetString("mysql.post")
	database := viper.GetString("mysql.database")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	charset := viper.GetString("mysql.charset")
	arg := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", username, password, host, post, database, charset)
	DB, err := gorm.Open(mysql.Open(arg), &gorm.Config{})
	if err != nil {
		log.Fatalln("fail database error", err)
	}
	// 默认建表
	DB.AutoMigrate(model.User{})
	return DB
}

// InitPgsql 连接PostgreSQL数据库
func InitPgsql() *gorm.DB {
	host := viper.GetString("postgresql.host")
	user := viper.GetString("postgresql.username")
	password := viper.GetString("postgresql.password")
	database := viper.GetString("postgresql.database")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, user, password, database)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("pgsql initialize fail:" + err.Error())
	}
	db.AutoMigrate(model.Path{})
	return db
}

// InitRedis redis连接
func InitRedis() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: "",
		DB:       0,
	})
	return conn
}
