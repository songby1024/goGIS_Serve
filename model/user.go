package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	*gorm.Model
	Id       int    `json:"id"`
	UserName string `gorm:"column:username" json:"username"`
	PassWord string `gorm:"column:password" json:"password"`
	Email    string `json:"email" gorm:"column:email"`
	Ruler    int    `json:"ruler" gorm:"column:ruler"`
}

func (User) TableName() string {
	return "user"
}

type Email struct {
	UserID    int       `json:"userId"`
	EmailText string    `json:"emailText"`
	SendTime  time.Time `json:"sendTime"`
}
