package global

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"serve/model"
	"time"
)

var (
	PGSQL        *gorm.DB
	DB           *gorm.DB
	DATACHAN     chan []byte
	Border       model.Border
	Period       int64 = 1
	LastSentTime int64
	NextSentTime int64
	TIMEZONE     *time.Location
	RDB          *redis.Client
)
