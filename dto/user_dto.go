package dto

import (
	"serve/model"
)

type UserInfo struct {
	Id         uint   `json:"id"`
	Username   string `json:"username"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

type UserResData struct {
	Users       []UserInfo `json:"users"`
	CurrentPage int        `json:"currentPage"`
	PageSize    int        `json:"pageSize"`
	Total       int        `json:"total"`
}

func UserInfoDto(user model.User) UserInfo {
	return UserInfo{
		Id:         user.ID,
		Username:   user.UserName,
		CreateTime: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateTime: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToUsers(users []model.User) []UserInfo {
	var res []UserInfo
	for _, x := range users {
		res = append(res, UserInfoDto(x))
	}
	return res
}
