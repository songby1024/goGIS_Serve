package dao

import (
	"fmt"
	"go.uber.org/zap"
	"serve/common/global"
	"serve/model"
	"time"
)

// GetMessageListByGeoId 获取预警历史
func GetMessageListByGeoId(geoId int, filed []string, page, pageSize int) ([]model.Messages, int64) {
	var list []model.Messages
	offset := (page - 1) * pageSize
	var count int64
	global.DB.Table("messages").Where("geo_id = ?", geoId).Count(&count)
	global.DB.Table("messages").
		Select(filed).
		Where("geo_id = ?", geoId).
		Order("alert_time DESC").
		Limit(pageSize).
		Offset(offset).Find(&list)
	return list, count
}

// CreateMessage 创建预警记录
func CreateMessage(message model.Messages) bool {
	tx := global.DB.Table("messages").Create(&message)
	return tx.RowsAffected == 1
}

// UpdateMessageState 修改预警状态
func UpdateMessageState(id, state int) bool {
	var msg model.Messages
	tx := global.DB.Table("messages").Where("id = ?", id).First(&msg)
	if tx.RowsAffected == 0 {
		return false
	}
	msg.State = state
	msg.UpdatedAt = time.Now()
	tx = global.DB.Save(&msg)
	if tx.RowsAffected == 0 {
		return false
	}
	return true
}

// GetMessageCount 获取各个围栏的预警总数 state: -1 为所有，0已处理，1未处理
func GetMessageCount(state int) *[]model.GeoIDCount {
	var results []model.GeoIDCount
	query := global.DB.Table("messages").Select("geo_id, COUNT(*) as total_count").
		Where("deleted_at IS NULL")
	if state > 0 {
		query = query.Where("state", state)
	}
	err := query.Group("geo_id").Scan(&results).Error
	if err != nil {
		zap.S().Error("Failed to execute query:", err)
		return nil
	}
	ids := []int64{}
	for _, v := range results {
		ids = append(ids, v.GeoID)
	}
	type Geo struct {
		GeoId int    `json:"geoId" gorm:"geo_id"`
		Name  string `json:"name" gorm:"name"`
	}
	var geoList []Geo
	global.DB.Table("messages").Select([]string{"geo_id", "name"}).Where("geo_id in ?", ids).Find(&geoList)
	fmt.Println("data:", geoList)
	tag := make(map[int64]Geo)
	for _, v := range geoList {
		tag[int64(v.GeoId)] = v
	}
	for i, v := range results {
		if item, ok := tag[v.GeoID]; ok {
			results[i].Name = item.Name
		}
	}
	return &results
}
