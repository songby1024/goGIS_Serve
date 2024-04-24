package services

import (
	"gorm.io/gorm"
	"log"
	"serve/common/global"
	"serve/model"
)

func NewConfig(client model.Client) error {
	var clientConfig model.Client
	err := global.PGSQL.Find(&clientConfig).Where("id=1").Error
	if err != gorm.ErrRecordNotFound {
		clientConfig.PosConfig = client.PosConfig
		clientConfig.Border = client.Border
		clientConfig.WarnPeriod = client.WarnPeriod
		global.PGSQL.Save(&clientConfig)
	}
	err = global.PGSQL.Create(&client).Error
	if err != nil {
		log.Println("create fail" + err.Error())
		return err
	}
	return nil
}

func GteConfig() (client model.Client, err error) {
	err = global.PGSQL.Find(&client).Where("id=1").Error
	return
}
