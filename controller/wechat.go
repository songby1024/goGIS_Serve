package controller

import (
	"github.com/gin-gonic/gin"
	"serve/services"
)

func GetQrcodeTicket(ctx *gin.Context) {
	ticket := services.GetWechatQrcodeTicket()
	if ticket == "" {
		ctx.JSON(500, gin.H{"code": 00, "errMsg": "获取失败"})
		return
	}
	ctx.JSON(200, gin.H{"code": 01, "data": ticket})
	return
}
