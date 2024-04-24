package dao

import (
	"fmt"
	"serve/common/global"
	"serve/model"
)

func GetUserList(ids []int, field []string) map[int]model.User {
	var userList []model.User
	global.DB.Table("user").Select(field).Where("id in ?", ids).Find(&userList)
	res := make(map[int]model.User)
	for _, user := range userList {
		res[user.Id] = user
	}
	fmt.Println("当前用户表：", userList)
	return res
}
