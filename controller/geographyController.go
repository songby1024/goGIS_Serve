package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"serve/common"
	"serve/model"
	"serve/response"
	"strings"
)

// SavaPath 拿到前端返回的Path保存到数据库
func SavaPath(ctx *gin.Context) {
	var psql = common.InitPgsql()
	psql.AutoMigrate(&model.Path{})

	paramData := model.PathStruct{}
	ctx.Bind(&paramData)
	// fmt.Println("拿到的地理信息", paramData)
	var paths string
	for _, v := range paramData.Path {
		path, _ := json.Marshal(&v)
		paths = paths + string(path) + ","
	}
	path := model.Path{
		Path: paths,
	}
	// 将数据写入postgreSQL数据库
	err := psql.Create(&path).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "数据写入失败"})
		return
	}
	response.Res200(ctx, gin.H{"msg": "数据写入成功"})
}

// AllPath 拿到数据库所有的path返回给前端
func AllPath(ctx *gin.Context) {
	pqsql := common.InitPgsql()
	pqsql.AutoMigrate(&model.Path{})
	var paths []model.Path
	err := pqsql.Find(&paths).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "数据查询失败"})
		return
	}
	var resPaths []model.PathStruct
	for _, v := range paths {
		// 设置一个空的PathStruct
		var path model.PathStruct

		// 反序列化path到path.path
		strS := strings.Split(v.Path, "},")
		var pathS []model.PathType
		for _, v3 := range strS {
			var p model.PathType
			json.Unmarshal([]byte(v3+"}"), &p)
			if p.Q == 0 && p.R == 0 && p.Lat == 0 && p.Lng == 0 {
				continue
			}
			pathS = append(pathS, p)
		}
		path.Path = pathS

		path.ID = v.ID
		resPaths = append(resPaths, path)
	}
	response.Res200(ctx, gin.H{"msg": "查询成功", "data": resPaths})
}

// DeletePath 删除数据库中的指定idPath
func DeletePath(ctx *gin.Context) {
	db := common.InitPgsql()
	db.AutoMigrate(model.Path{})
	id, _ := ctx.GetQuery("id")
	var newInfo model.Path
	err := db.Where("id = ?", id).First(&newInfo).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "删除失败"})
		return
	}
	err = db.Delete(&newInfo).Error
	if err != nil {
		response.Res201(ctx, gin.H{"msg": "删除失败"})
		return
	}
	response.Res200(ctx, gin.H{"msg": "删除成功"})
}
