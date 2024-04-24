package model

import (
	"github.com/dgrijalva/jwt-go"
)

type Point struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type PointFloat struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

// Content 预警数据
type Content struct {
	Id          int     `json:"id"`
	AddressName string  `json:"addressName"`
	AlertTime   string  `json:"alterTime"`
	AlertClass  string  `json:"alertClass"`
	Point       Point   `json:"point"`
	AlertDic    string  `json:"alertDic"`
	MinDistance float64 `json:"minDistance"`
	State       int     `json:"state"`
}

type CustomClaims struct {
	Address string
	jwt.StandardClaims
}

type Geofence struct {
	ID                int
	Name              string
	MinDistanceMeters float64 `gorm:"column:MinDistanceMeters"`
}
