package controller

import (
	"github.com/gin-gonic/gin"
	"serve/common"
	"serve/model"
	"serve/response"
	"serve/utils"
	"time"
)

func Login(ctx *gin.Context) {

	var DB = common.InitDB()

	//  接受用户登录信息
	var userLoginInfo model.User
	ctx.Bind(&userLoginInfo)

	// 用户信息校验
	// 判断用户名是否存在
	var newUser model.User
	err := DB.Where(" username = ?", userLoginInfo.UserName).First(&newUser).Error
	if err != nil { // 	用户名不存在
		response.Res202(ctx, gin.H{"msg": "账户不存在"})
		return
	}

	// 验证密码
	if userLoginInfo.PassWord != newUser.PassWord {
		// 	密码错误
		response.Res202(ctx, gin.H{"msg": "密码错误"})
		return
	}
	// 返回成功地请求响应
	//if userLoginInfo.UserName == "admin" && userLoginInfo.PassWord == "admin" {
	//	response.Res200(ctx, gin.H{"msg": "success", "user": userLoginInfo.UserName})
	//} else {
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Res202(ctx, gin.H{"msg": "登录失败"})
		return
	}

	data := map[string]interface{}{
		"user":  newUser.ID,
		"token": token,
	}
	response.Res200(ctx, gin.H{"msg": "success", "user": newUser, "token": token, "data": data})
	//}
}

// Register 处理注册请求
func Register(ctx *gin.Context) {
	var DB = common.InitDB()

	//  接受用户登录信息
	var userLoginInfo model.User
	ctx.Bind(&userLoginInfo)

	// 用户信息校验
	// 判断用户名是否存在
	var newUser model.User
	err := DB.Where(" username = ?", userLoginInfo.UserName).First(&newUser).Error
	if err != nil { // 	用户名不存在
		err = DB.Create(&userLoginInfo).Error
		if err != nil {
			response.Res202(ctx, gin.H{"msg": "账户已存在"})
			return
		} else {
			response.Res200(ctx, gin.H{"msg": "注册成功"})
			return
		}
	} else {
		response.Res202(ctx, gin.H{"msg": "账户已存在"})
		return
	}
}

// AllUser 获取所有用户信息
func AllUser(ctx *gin.Context) {
	var users []model.User
	var db = common.InitDB()
	err := db.Find(&users).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "查询用户失败"})
		return
	}
	response.Res200(ctx, gin.H{"msg": "查询成功", "users": users})
}

type SendEmailReq struct {
	UserName string `json:"username"`
	Content  string `json:"content"`
}

// SendEmail 发送告警邮件
func SendEmail(ctx *gin.Context) {
	// 连接数据库
	db := common.InitDB()
	db.AutoMigrate(&model.User{})
	// 获取参数
	emailInfo := SendEmailReq{}
	ctx.Bind(&emailInfo)
	// 查询用户信息
	var userInfo model.User
	err := db.Where("username = ?", emailInfo.UserName).First(&userInfo).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "查询用户失败"})
		return
	}
	if userInfo.Email == "" {
		response.Res201(ctx, gin.H{"msg": "用户邮箱为空"})
		return
	}

	// 	发送邮件
	err = utils.SendEMail(userInfo.Email, emailInfo.Content)
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "邮件发送失败"})
		return
	}

	// 	将发送邮件的数据写入redis中间缓存数据库
	resDb := common.InitRedis()
	defer resDb.Close()

	// 	将用户ID、邮件内容、创建时间构造为一个list列表
	err = resDb.LPush("email", userInfo.UserName, emailInfo.Content, time.Now().
		Format("2006-01-02 15:04:05")).Err()
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "邮件存储失败"})
		return
	}
	response.Res200(ctx, gin.H{"msg": "邮件发送成功"})
}

// GetAdminList 获取所有的管理员列表
func GetAdminList(ctx *gin.Context) {
	// 连接数据库
	db := common.InitDB()
	var adminList []model.User
	err := db.Where("ruler = ?", 1).Find(&adminList).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "查询管理员列表失败", "list": ""})
		return
	}
	response.Res200(ctx, gin.H{"msg": "查询成功", "list": adminList})
}
