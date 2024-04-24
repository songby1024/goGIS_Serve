package routers

import (
	"github.com/gin-gonic/gin"
	"serve/controller"
	"serve/middleware"
)

func Router(r *gin.Engine) {

	// 暴露静态资源
	// 处理跨域
	r.Use(middleware.CorsMiddle())

	// 用户接口
	userApi := r.Group("/user")
	userRouter(userApi)

	// 	业务接口
	r.POST("/geography/savaPath", controller.SavaPath)
	r.GET("/geography/getAllPath", controller.AllPath)
	r.DELETE("/geography/deletePath", controller.DeletePath)

	r.POST("/inSideInfo/setInSideInfo", controller.SetInSideInfo)
	r.GET("/inSideInfo/GetAllInSideInfo", controller.GetAllInSideInfo)

	r.GET("/ws", controller.WsHandle)
	r.POST("/ws/dataHandle", controller.DataHandle)
	// r.GET("/ws", controller.WsHandle)
	// api := r.Group("/api")
	// api.POST("/ws", controller.DataHandle)
	// api.POST("/init", controller.InitConfig)
	// api.GET("/query", controller.Querydb)
	// api.GET("/qrcode", controller.GetQrcodeTicket)

	alert := r.Group("/api")
	{
		//围栏相关接口
		alertService := controller.NewAlertNotice()
		alert.GET("/alterCheck", alertService.AlertCheckCenter)
		alert.POST("alert/add", alertService.AddGeofence)
		alert.GET("/alert/detail", alertService.GetGeofenceById)
		alert.GET("/alert/list", alertService.GetGeofenceList)
		alert.POST("alert/update", alertService.UpdateGeofence)

		//预警统计相关接口
		messageService := controller.NewMessage()
		alert.GET("/message/list", messageService.GetMessageList)
		alert.GET("/message/total", messageService.GetMessageCount)
		alert.POST("/message/update", messageService.UpdateMessageState)
	}

}
