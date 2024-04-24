package controller

import (
	"github.com/gin-gonic/gin"
	"serve/common"
	"serve/model"
	"serve/response"
)

// SetInSideInfo 获取报警信息
func SetInSideInfo(ctx *gin.Context) {
	var insideInfo model.InSideInfo
	db := common.InitDB()
	db.AutoMigrate(model.InSideInfo{})
	ctx.Bind(&insideInfo)
	// fmt.Println("insideInfo", insideInfo)
	err := db.Create(&insideInfo).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "写入失败"})
		return
	}
	response.Res200(ctx, gin.H{"msg": "导入信息成功"})
}

// GetAllInSideInfo 获取全部的报警信息
func GetAllInSideInfo(ctx *gin.Context) {
	db := common.InitDB()
	var insideInfo []model.InSideInfo
	err := db.Find(&insideInfo).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "查询用户报警失败"})
		return
	}
	response.Res200(ctx, gin.H{"data": insideInfo})
}
