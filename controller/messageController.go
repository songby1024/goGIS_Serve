package controller

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-playground/validator/v10"
	"serve/response"
	"serve/services"
	"strconv"
)

type MessageController struct {
}

func NewMessage() *MessageController {
	return &MessageController{}
}

// GetMessageList 获取预警消息列表
func (c *MessageController) GetMessageList(ctx *gin.Context) {
	geoIdStr := ctx.DefaultQuery("geoId", "0")
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("pageSize", "10")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	geoId, _ := strconv.Atoi(geoIdStr)

	res := services.GetMessageList(geoId, page, pageSize)
	response.Res200(ctx, gin.H{"msg": "succeed", "data": res})
}

// GetMessageCount 获取预警总数和待处理总数
func (c *MessageController) GetMessageCount(ctx *gin.Context) {
	data := services.GetMessageCount()
	response.Res200(ctx, gin.H{"msg": "succeed", "data": data})
}

// UpdateMessageState 更新预警状态
func (c *MessageController) UpdateMessageState(ctx *gin.Context) {
	type Params struct {
		Id    int `json:"id" validate:"required"`
		State int `json:"state" validate:"required,oneof=0 1"`
	}
	var params Params
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.Res201(ctx, gin.H{"msg": "参数校验失败"})
		return
	}
	if !services.UpdateMessageState(params.Id, params.State) {
		response.Res201(ctx, gin.H{"msg": "更新失败"})
		return
	}
	response.Res200(ctx, gin.H{"msg": "succeed", "data": nil})
}
