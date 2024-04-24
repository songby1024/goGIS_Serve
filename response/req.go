package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 将传入的数据进行解析合并
func resData(code int, data gin.H) gin.H {
	resVal := gin.H{"code": code}
	for x, y := range data {
		resVal[x] = y
	}
	return resVal
}

// 封装响应返回函数

// Res200 成功得回调
func Res200(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusOK, resData(200, data))
}

// Res201 请求成功并且创建资源
func Res201(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusCreated, resData(201, data))
}

// Res202 已被接收得回调
func Res202(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusAccepted, resData(202, data))
}

// Res400 服务器不理解请求的语法。
func Res400(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusBadRequest, resData(400, data))
}

// Res401 未授权
func Res401(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusBadRequest, resData(401, data))
}

// Res405 服务端禁止的访问方法
func Res405(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusBadRequest, resData(405, data))
}

// Res417 未满足期望请求标头字段
func Res417(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusBadRequest, resData(417, data))
}

// Res500 服务器内部错误
func Res500(ctx *gin.Context, data gin.H) {
	ctx.JSON(http.StatusBadRequest, resData(500, data))
}
