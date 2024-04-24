package global

import (
	"gorm.io/gorm"
	"serve/model"
	"time"
)

var (
	PGSQL        *gorm.DB
	DATACHAN     chan []byte
	Border       model.Border
	Period       int64 = 1
	LastSentTime int64
	NextSentTime int64
	TIMEZONE     *time.Location
)
