package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"serve/model"
	"serve/response"
	"serve/services"
	"strconv"
)

type AlertNoticeController struct {
}

func NewAlertNotice() *AlertNoticeController {
	return &AlertNoticeController{}
}

func (c *AlertNoticeController) AddGeofence(ctx *gin.Context) {
	var newGeofence model.GeofenceModel
	if err := ctx.ShouldBindJSON(&newGeofence); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := services.CreateGeofence(newGeofence)
	// 插入数据至数据库
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "数据写入失败"})
		return
	}
	response.Res200(ctx, gin.H{"msg": "数据写入成功", "data": nil})
}

// GetGeofenceById 获取围栏详情
func (c *AlertNoticeController) GetGeofenceById(ctx *gin.Context) {
	geoIdStr := ctx.DefaultQuery("geoId", "0")
	geoId, _ := strconv.Atoi(geoIdStr)
	data := services.GetGeofenceById(geoId)
	response.Res200(ctx, gin.H{"msg": "succeed", "data": data})
}

// GetGeofenceList 获取围栏列表
func (c *AlertNoticeController) GetGeofenceList(ctx *gin.Context) {
	list := services.GetGeofenceList()
	data := map[string]interface{}{
		"list":  list,
		"total": len(list),
	}
	response.Res200(ctx, gin.H{"msg": "succeed", "data": data})
}

// UpdateGeofence 更新围栏信息,目前只支持：名称，描述，状态，负责人id等修改
func (c *AlertNoticeController) UpdateGeofence(ctx *gin.Context) {
	var newGeofence model.GeofenceModel
	if err := ctx.ShouldBindJSON(&newGeofence); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := services.UpdateGeofence(newGeofence.ID, newGeofence.Name, newGeofence.Description, newGeofence.Status, newGeofence.AlertDistans, newGeofence.ManagerIDs)
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "更新数据失败"})
		return
	}
	response.Res200(ctx, gin.H{"msg": "succeed", "data": nil})
}

func (c *AlertNoticeController) AlertCheckCenter(ctx *gin.Context) {
	services.AlertCheckCenter(ctx.Writer, ctx.Request)
}
