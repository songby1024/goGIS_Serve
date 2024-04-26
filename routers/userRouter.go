package routers

import (
	"github.com/gin-gonic/gin"
	"serve/controller"
)

func userRouter(r *gin.RouterGroup) {
	// 登录注册的业务
	// 登录
	r.POST("/login", controller.Login)
	// 注册
	r.POST("/register", controller.Register)
	// 	返回所有用户信息
	r.GET("/allUser", controller.AllUser)

	// 	发送告警信息给用户
	r.POST("/sendEmail", controller.SendEmail)
	//r.POST("/getAllAdminList", controller.GetAdminList)
}
