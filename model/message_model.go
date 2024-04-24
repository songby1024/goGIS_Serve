package model

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Messages struct {
	Model
	GeoID       uint64  `gorm:"not null"` // 地理位置ID
	Name        string  `gorm:"not null"` // 围栏名称
	ClientID    uint64  `gorm:"not null"` // 客户端ID
	AddressName string  `gorm:"size:255"` // 地址名称
	AlertTime   string  `gorm:"size:255"` // 警告时间
	AlertClass  string  `gorm:"size:255"` // 警告类别
	AlertDic    string  `gorm:"size:255"` // 警告词典
	MinDistance float64 `gorm:""`         // 最小距离
	State       int     `gorm:""`         // 状态
	PointLat    float64 `gorm:""`         // 维度
	PointLng    float64 `gorm:""`         // 经度
}

func (m Messages) GetTableName() string {
	return "messages"
}

type MessageResp struct {
	Id          uint    `json:"id"`
	GeoID       uint64  `json:"geoId"`        // 地理位置ID
	ClientID    uint64  `json:"clientId"`     // 客户端ID
	AddressName string  `json:"addressName""` // 地址名称
	AlertTime   string  `json:"alertTime"`    // 警告时间
	AlertClass  string  `json:"alertClass"`   // 警告类别
	AlertDic    string  `json:"alertDic"`     // 警告词典
	MinDistance float64 `json:"minDistance"`  // 最小距离
	State       int     `json:"state"`        // 状态
	PointLat    float64 `json:"pointLat"`     // 维度
	PointLng    float64 `json:"pointLng"`     // 经度
}

type GeoIDCount struct {
	GeoID      int64  `json:"geoID"`
	TotalCount int64  `json:"totalCount"`
	Name       string `json:"name"`
}
