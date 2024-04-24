package dao

import (
	"github.com/spf13/viper"
	"log"
	"serve/initialize"
	"serve/model"
	"testing"
)

//import (
//	"fmt"
//	"serve/model"
//	"testing"
//)

//func TestCheckPointInAlertRange(t *testing.T) {
//	isContant := CheckPointInAlertRange(1, model.Point{
//		Longitude: "-74.0060",
//		Latitude:  "41.7128",
//	}, "")
//	fmt.Println(isContant)
//}

// 在主函数运行之前，读取配置文件
func initConfig() {
	viper.SetConfigName("application-local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("config read fail error = ", err)
	}
}

func TestCreateGeofence(t *testing.T) {
	initialize.InitConfig()
	initialize.InitGlobal()

	p := "-73.9721 40.7934,-73.9601 40.8007,-73.9583 40.7968,-73.9737 40.7845,-73.9721 40.7934"
	err := CreateGeofence(model.GeofenceModel{
		Name:         "无锡苏南硕放机场",
		Status:       1,
		CityName:     "无锡市",
		CityCoords:   &model.Coordinate{Lat: "31.342", Lng: "120.2143"},
		Boundary:     nil,
		ManagerIDs:   []int{1, 2, 7},
		AlertArea:    nil,
		Description:  "无锡军民两用机场",
		AlertDistans: 0,
	}, p, p, "1,5,7")
	if err != nil {
		t.Error("err:", err)
	}
}

func TestGetGeofenceById(t *testing.T) {
	initialize.InitConfig()
	initialize.InitGlobal()
	data := GetGeofenceById(11)
	t.Log("data:", data)
}
