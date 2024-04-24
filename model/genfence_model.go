package model

import "time"

type GeofenceInfo struct {
	Id           int
	Name         string
	CityName     string
	Boundary     string
	AlertArea    string
	ManagerIds   string     `json:"manager_ids"` // 这里映射JSON中的 manager_ids 字段
	Description  string     `json:"description"`
	AlertDistans int        `json:"alert_distans"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	Status       int        `json:"status"`
}

// GeofenceModel 表示geofences表中的一行数据
type GeofenceModel struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	Status       int          `json:"status"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	CityName     string       `json:"city_name"`
	CityCoords   *Coordinate  `json:"city_coords"`
	Boundary     []Coordinate `json:"boundary"`
	ManagerIDs   []int        `json:"manager_ids"` // ID数组
	AlertArea    []Coordinate `json:"alert_area"`
	Description  string       `json:"description"`
	AlertDistans int          `json:"alert_distans"`
}

type Coordinate struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
