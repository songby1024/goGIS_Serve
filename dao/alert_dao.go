package dao

import (
	"fmt"
	"serve/common/global"
	"serve/model"
	"strconv"
)

//func GetAlertStateById(id int) *model.FenceAlertModel {
//	var alertState model.FenceAlertModel
//	tx := global.DB.Where(model.FenceAlertColumns.GeofenceID, id).
//		Where(model.FenceAlertColumns.IsActive, 1).First(&alertState)
//	if tx.RowsAffected == 0 {
//		return nil
//	}
//	return &alertState
//}
//
//// GetAlertInfoById 获取围栏信息
//func GetAlertInfoById(id int, filed []string) *model.GeofenceModel {
//	var alertState model.GeofenceModel
//	tx := global.DB.Where(model.GeofenceColumns.ID, id).
//		Where(model.GeofenceColumns, 1).First(&alertState)
//	if tx.RowsAffected == 0 {
//		return nil
//	}
//	return &alertState
//}

// CheckPointInAlertRange 判断点是否在围栏中, alert_area预警围栏，boundary实际虚拟围栏
func CheckPointInAlertRange(geoId int, point model.Point, ship string) (model.GeofenceInfo, bool) {
	geoFence := model.GeofenceInfo{}
	// 构建查询
	result := global.PGSQL.Raw("SELECT id, name, city_name, boundary, alert_area, manager_ids FROM geofences WHERE ST_Contains("+ship+", ST_SetSRID(ST_MakePoint(?, ?), 4326)) and id = ?", point.Longitude, point.Latitude, geoId).Scan(&geoFence)
	return geoFence, result.RowsAffected == 1
}

// GetCentroid 获取多边形质心坐标
func GetCentroid(geoId int64) model.Point {
	// 执行查询
	var centroid model.Point
	global.PGSQL.Raw("SELECT ST_X(ST_Centroid(boundary)) AS longitude, ST_Y(ST_Centroid(boundary)) AS latitude FROM geofences WHERE id = ?", geoId).Scan(&centroid)
	return centroid
}

// GetMinDistanceMeters 获取点到多边形的距离
func GetMinDistanceMeters(geoId int, point model.Point) model.Geofence {
	pointSql := fmt.Sprintf("SELECT ST_SetSRID(ST_MakePoint(%s, %s), 4326)", point.Longitude, point.Latitude)
	sql := "WITH TargetPoint AS (" + pointSql + "as geom\n    )\n    SELECT\n        geofences.id,\n        geofences.name,\n        ST_DistanceSphere(geofences.boundary, tp.geom) AS \"MinDistanceMeters\"\n    FROM\n        geofences,\n        TargetPoint tp\n    WHERE\n        geofences.id = " + strconv.Itoa(geoId)
	// 执行查询
	var geofence model.Geofence
	global.PGSQL.Raw(sql).Scan(&geofence)
	return geofence
}
