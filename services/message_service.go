package services

import (
	"serve/dao"
	"serve/model"
)

func GetMessageList(geoId int, page, pageSize int) map[string]interface{} {
	list, total := dao.GetMessageListByGeoId(geoId, []string{}, page, pageSize)
	res := make([]model.MessageResp, 0, len(list))
	for _, msg := range list {
		res = append(res, model.MessageResp{
			Id:          msg.ID,
			GeoID:       msg.GeoID,
			ClientID:    msg.ClientID,
			AddressName: msg.AddressName,
			AlertTime:   msg.AlertTime,
			AlertClass:  msg.AlertClass,
			AlertDic:    msg.AlertDic,
			MinDistance: msg.MinDistance,
			State:       msg.State,
			PointLat:    msg.PointLat,
			PointLng:    msg.PointLng,
		})
	}
	return map[string]interface{}{
		"list":  res,
		"total": total,
	}
}

// GetMessageCount 获取各个围栏预警总数
func GetMessageCount() map[string]interface{} {
	allCount := dao.GetMessageCount(-1)
	unProcess := dao.GetMessageCount(1)
	return map[string]interface{}{
		"allMessage": allCount,
		"unProcess":  unProcess,
	}
}

// UpdateMessageState 更新状态
func UpdateMessageState(Id, state int) bool {
	return dao.UpdateMessageState(Id, state)
}
