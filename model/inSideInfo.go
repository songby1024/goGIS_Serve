package model

type InSideInfo struct {
	UserName string  `gorm:"column:username" json:"username" `
	Lng      float64 `json:"lng"`
	Lat      float64 `json:"lat"`
	Accuracy float64 `json:"accuracy"`
}

func (InSideInfo) TableName() string {
	return "insideInfo"
}
