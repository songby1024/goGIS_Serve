package model

type GeofenceInfo struct {
	Id         int
	Name       string
	CityName   string
	Boundary   string
	AlertArea  string
	ManagerIds string `json:"manager_ids"` // 这里映射JSON中的 manager_ids 字段
}
