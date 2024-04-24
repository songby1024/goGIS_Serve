package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"serve/model"
	"serve/services"
)

func InitConfig(ctx *gin.Context) {
	// var clientConfig interface{}
	clientConfig := model.Client{}
	_ = ctx.ShouldBindJSON(&clientConfig)
	err := services.NewConfig(clientConfig)
	if err != nil {
		log.Printf("add config err :%v", err)
		ctx.JSON(200, gin.H{"code": 201, "errMsg": "fail"})
		return
	}
	ctx.JSON(200, gin.H{"code": 200, "errMsg": "OK"})
	return
}
