package middleware

import (
	"github.com/gin-gonic/gin"
	"serve/common"
	"serve/model"
	"serve/response"
	"strings"
)

// AuthMiddle 认证中间件
func AuthMiddle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//  1 获取请求头中的("Authorization")数据
		tokenString := ctx.GetHeader("Authorization")

		// 	2 验证数据格式是否正确
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			response.Res417(ctx, gin.H{"msg": "权限不足(标头错误)"})
			ctx.Abort()
			return
		}
		// 	3 解析token
		tokenString = tokenString[6:]
		token, claims, err := common.ParseToken(tokenString)

		// 	4 验证token是否有效
		if err != nil || !token.Valid {
			response.Res417(ctx, gin.H{"msg": "权限不足(token已过期)"})
			ctx.Abort()
			return
		}

		// 	5 获取claims中的userID
		userId := claims.UserId
		user := model.User{}
		DB := common.InitDB()

		// 	6 数据库中查询userID
		err = DB.First(&user, userId).Error
		if err != nil {
			response.Res417(ctx, gin.H{"msg": "权限不足(查询用户失败)"})
			ctx.Abort()
			return
		}

		// 	7 验证userID是否合法
		if user.ID == 0 {
			response.Res417(ctx, gin.H{"msg": "权限不足(用户不存在)"})
			ctx.Abort()
			return
		}
		//  8 userId合法，则将user信息写入上下文context中
		ctx.Set("user", user)
		ctx.Next()
	}
}
