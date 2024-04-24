package controller

import (
	"github.com/gin-gonic/gin"
	"serve/services"
)

type AlertNoticeController struct {
}

func NewAlertNotice() *AlertNoticeController {
	return &AlertNoticeController{}
}

//	data, err := dao.GetDonateProjectsDetail(id)
//	if err != nil {
//		common.RespFail(ctx.Writer, "Failed to obtain data, please try again")
//		return
//	}
//	common.RespOKList(ctx.Writer, data)
//}

func (c *AlertNoticeController) AlertCheckCenter(ctx *gin.Context) {
	services.AlertCheckCenter(ctx.Writer, ctx.Request)
}
