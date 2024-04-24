package services

import (
	"encoding/json"
	"fmt"
	"log"
	"serve/common/global"
	"serve/model"
	"serve/utils"
	"strconv"
	"time"
)

func Coordinate(position model.QtPosition) {
	qpos := position
	border := global.Border
	var pos model.Position
	anc1, err := strconv.ParseFloat(qpos.ANC1, 64)
	anc2, err := strconv.ParseFloat(qpos.ANC2, 64)
	anc3, err := strconv.ParseFloat(qpos.ANC3, 64)
	x, err := strconv.ParseFloat(qpos.X, 64)
	y, err := strconv.ParseFloat(qpos.Y, 64)
	t := time.Now().UnixNano()
	// p := influxdb2.NewPoint("position",
	//	map[string]string{"unit": "pos"},
	//	map[string]interface{}{"x": x, "y": y},
	//	time.Now().In(global.TIMEZONE).Truncate(time.Nanosecond),
	// )
	// influx := new(initialize.Influx)
	// client := influx.GetInflux()
	// err = influx.Write(client, p)
	if err != nil {
		log.Println(err.Error())
		panic(err.Error())
	}
	if err != nil {
		panic(err.Error())
	}
	pos = model.Position{
		X:         x,
		Y:         y,
		TimeStamp: t,
		ANC1:      anc1 / 1000,
		ANC2:      anc2 / 1000,
		ANC3:      anc3 / 1000,
	}
	if !utils.IsOutsideArea(pos, border) {
		if pos.TimeStamp >= global.NextSentTime {
			global.LastSentTime = pos.TimeStamp
			global.NextSentTime = global.LastSentTime + 60*global.Period
			go SendMessage(pos.TimeStamp)
		}
		fmt.Println("未到报警时间")
	}
	posRes := model.PositionResponse{
		DataType: "coordinate",
		Position: pos,
	}
	jsonPos, _ := json.Marshal(posRes)
	global.DATACHAN <- jsonPos
	fmt.Println(border, pos)
}

func Distance(position model.QtPositionRequest) {
	pos := position.QtPosition
	anc1, err := strconv.ParseFloat(pos.ANC1, 64)
	anc2, err := strconv.ParseFloat(pos.ANC2, 64)
	anc3, err := strconv.ParseFloat(pos.ANC3, 64)
	// t, err := strconv.ParseInt(pos.TimeStamp, 10, 64)
	if err != nil {
		panic(err)
	}
	distance := model.Distance{
		// TimeStamp: t,
		ANC1: anc1 / 1000,
		ANC2: anc2 / 1000,
		ANC3: anc3 / 1000,
	}
	distanceRes := model.DistanceResponse{
		DataType: "distance",
		Distance: distance,
	}
	jsonPos, _ := json.Marshal(distanceRes)
	global.DATACHAN <- jsonPos
	fmt.Println(pos)
}
