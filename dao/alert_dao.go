package dao

import (
	"fmt"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"serve/common/global"
	"serve/common/tools"
	"serve/model"
	"strconv"
	"strings"
)

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

// CreateGeofence 创建围栏
func CreateGeofence(geo model.GeofenceModel, boundary string, alertArea string, uid string) error {
	point := fmt.Sprintf("POINT(%s %s)", geo.CityCoords.Lng, geo.CityCoords.Lat)
	boundary = fmt.Sprintf("POLYGON((%s))", boundary)
	alertArea = fmt.Sprintf("POLYGON((%s))", alertArea)

	uid = "1,2,4"

	sql := `INSERT INTO geofences (
	name,
	status,
	city_name,
	city_coords,
	boundary,
	manager_ids,
	alert_area,
	description,
	alert_distans
) VALUES (
	?,
	1,
	?,
	ST_GeomFromText(?, 4326),
	ST_GeomFromText(?, 4326),
	?,
	ST_GeomFromText(?, 4326),
	?,
	?
);`

	err := global.PGSQL.Exec(sql,
		geo.Name,
		geo.CityName,
		point,
		boundary,
		pq.Array(geo.ManagerIDs),
		alertArea,
		geo.Description,
		geo.AlertDistans).Error
	if err != nil {
		zap.S().Error("创建围栏失败", err)
		return err
	}
	return nil
}

// GetGeofenceById 获取围栏详情
func GetGeofenceById(geoId int) model.GeofenceModel {
	var geofence model.GeofenceInfo
	sql := fmt.Sprintf("SELECT id, name, status, created_at, updated_at, city_name, city_coords, ST_AsText(boundary::geometry) as boundary, manager_ids, ST_AsText(alert_area) as alert_area , description, alert_distans FROM geofences WHERE id=%d and status != 0;;", geoId)
	global.PGSQL.Raw(sql).Scan(&geofence)

	boundary, _ := parseWKTToCoordinates(geofence.Boundary)
	alertArea, _ := parseWKTToCoordinates(geofence.AlertArea)
	managerIds, _ := tools.ParseIntSliceFromString(geofence.ManagerIds)

	return model.GeofenceModel{
		ID:           geofence.Id,
		Name:         geofence.Name,
		Status:       geofence.Status,
		CreatedAt:    *geofence.CreatedAt,
		UpdatedAt:    *geofence.UpdatedAt,
		CityName:     geofence.CityName,
		Boundary:     boundary,
		ManagerIDs:   managerIds,
		AlertArea:    alertArea,
		Description:  geofence.Description,
		AlertDistans: geofence.AlertDistans,
	}
}

// GetGeofenceList 获取围栏列表
func GetGeofenceList() []model.GeofenceModel {
	var geoList []model.GeofenceInfo
	sql := "SELECT * FROM geofences WHERE status != 0;"
	global.PGSQL.Raw(sql).Scan(&geoList)

	res := make([]model.GeofenceModel, 0, len(geoList))
	for _, geofence := range geoList {

		boundary, _ := parseWKTToCoordinates(geofence.Boundary)
		alertArea, _ := parseWKTToCoordinates(geofence.AlertArea)
		managerIds, _ := tools.ParseIntSliceFromString(geofence.ManagerIds)

		res = append(res, model.GeofenceModel{
			ID:           geofence.Id,
			Name:         geofence.Name,
			Status:       geofence.Status,
			CreatedAt:    *geofence.CreatedAt,
			UpdatedAt:    *geofence.UpdatedAt,
			CityName:     geofence.CityName,
			Boundary:     boundary,
			ManagerIDs:   managerIds,
			AlertArea:    alertArea,
			Description:  geofence.Description,
			AlertDistans: geofence.AlertDistans,
		})
	}
	return res
}

func UpdateGeofence(geoId int, name string, des string, state int, alertDist int, managerIds []int) error {
	update := make(map[string]interface{})
	if name != "" {
		update["name"] = name
	}
	if state >= 0 {
		update["status"] = state
	}
	if des != "" {
		update["description"] = des
	}
	if alertDist != 0 {
		update["alert_distans"] = alertDist
	}
	if len(managerIds) > 0 {
		update["manager_ids"] = pq.Array(managerIds)
	}
	result := global.PGSQL.Table("geofences").Where("id = ?", geoId).Updates(update)
	if result.Error != nil {
		zap.S().Error("Error updating geofence", result.Error)
		return result.Error
	}
	return nil

	//filed := "UPDATE geofences SET"
	//if name != "" {
	//	filed += fmt.Sprintf("name = '%s,'", name)
	//}
	//if state != 0 {
	//	filed += fmt.Sprintf("status = %d,", state)
	//}
	//if des != "" {
	//	filed += fmt.Sprintf("description = '%s',", des)
	//}
	//if alertDist != 0 {
	//	filed += fmt.Sprintf("alert_distans = %d',", alertDist)
	//}
	//if len(managerIds) > 0 {
	//	filed += fmt.Sprintf("manager_ids = %d',", pq.Array(managerIds))
	//}
	//sql := "UPDATE geofences SET  name = '无锡学院', status = 2, description = '吃了吗',  updated_at = CURRENT_TIMESTAMP\nWHERE id = 1;"
}

func parseWKTToCoordinates(wktPolygon string) ([]model.Coordinate, error) {
	// 删除WKT中多余的部分
	trimSentence := strings.TrimPrefix(wktPolygon, "POLYGON((")
	trimSentence = strings.TrimSuffix(trimSentence, "))")

	// 使用逗号分隔每个坐标对
	coordinatePairs := strings.Split(trimSentence, ",")

	var coordinates []model.Coordinate
	for _, pair := range coordinatePairs {
		// 清除空格并分割经度和纬度
		cleanedPair := strings.TrimSpace(pair)
		points := strings.Split(cleanedPair, " ")

		if len(points) != 2 {
			return nil, fmt.Errorf("invalid coordinate pair: %s", pair)
		}

		// 创建并填充坐标结构
		coord := model.Coordinate{
			Lng: points[0],
			Lat: points[1],
		}
		coordinates = append(coordinates, coord)
	}
	return coordinates, nil
}
