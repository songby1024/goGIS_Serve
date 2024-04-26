package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/http"
	"serve/common/global"
	"serve/common/tools"
	"serve/dao"
	"serve/model"
	"serve/utils"
	"strconv"
	"sync"
	"time"
)

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type pointInfo struct {
	point  model.PointFloat
	geoId  int
	uid    int
	isReal bool
}

// Message 前端只需要传入：geoId，alertGeoId，clientId和当前定位point,取余字段由后端返回
type Message struct {
	Model
	GeofenceId int64         `json:"geoId"`      //围栏
	AlertGeoId int64         `json:"alertGeoId"` //触发预警围栏
	clientId   int64         `json:"clientId"`   //信息接收者
	AlertType  int           //预警类型  0正常， 1：低级， 2中级，3高级
	Content    model.Content `json:"content"` //预警内容
}

// Node 构造连接
type Node struct {
	Conn      *websocket.Conn //socket连接
	Addr      string          //客户端地址
	DataQueue chan []byte     //消息内容
	ClientId  int             //客户端唯一id
	JoinTime  int64           //节点创建时间
	GeoId     int             //归属围栏
	isReal    bool            //是否在围栏内
}

// 映射关系：围栏id绑定
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 绑定用户模拟用户信息
var pointLocker sync.RWMutex
var manyMeantimePointMap map[int64]map[int]pointInfo = make(map[int64]map[int]pointInfo)
var manyPointInfoMap map[int]pointInfo = make(map[int]pointInfo)
var pointQueue chan pointInfo = make(chan pointInfo, 3)

// 读写锁，绑定node时需要线程安全
var rwLocker sync.RWMutex

// AlertCheckCenter 需要 ：ID ，接受者ID ，消息类型，发送的内容
func AlertCheckCenter(w http.ResponseWriter, r *http.Request) {
	//1.  获取参数信息发送者userId
	query := r.URL.Query()
	uuid := query.Get("clientId")

	clientId, err := strconv.ParseInt(uuid, 10, 64)
	if err != nil {
		zap.S().Error("类型转换失败", err)
		return
	}

	//升级为socket
	var isvalida = true
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(w, r, nil)
	if err != nil {
		zap.S().Error(err)
		return
	}

	//获取socket连接,构造消息节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		ClientId:  int(clientId),
		JoinTime:  time.Now().Unix(),
	}

	//将userId和Node绑定
	rwLocker.Lock()
	fmt.Println("client:", clientId)
	clientMap[clientId] = node
	rwLocker.Unlock()

	//服务发送消息
	go sendProc(node)

	//服务接收消息
	go recProc(node)

	sendMsg(clientId, []byte("您已加入预警节点"))
}

// recProc 从websocket中将消息体拿出，然后进行解析，再进行信息类型判断， 最后将消息发送至目的用户的node中
func recProc(node *Node) {
	for {
		//获取信息
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("读取消息失败", err)
			return
		}

		//这里是简单实现的一种方法
		msg := Message{}
		err = json.Unmarshal(data, &msg)
		if err != nil {
			zap.S().Error("json解析失败", err)
			return
		}

		//注意：不处理管理员上传的位置信息，只针对普通用户
		uDate := dao.GetUserList([]int{node.ClientId}, []string{"id", "email", "username", "ruler"})
		if user, ok := uDate[node.ClientId]; ok && user.Ruler == 1 {
			continue
		}

		//检查用户是否进入预警区域，预警级别：低级
		geofenceInfo, check := dao.CheckPointInAlertRange(int(msg.GeofenceId), model.Point{
			Latitude:  msg.Content.Point.Latitude,
			Longitude: msg.Content.Point.Longitude,
		}, "alert_area")

		//检查是否进入围栏内，进入围栏不预警
		_, onGonFence := dao.CheckPointInAlertRange(int(msg.GeofenceId), model.Point{
			Latitude:  msg.Content.Point.Latitude,
			Longitude: msg.Content.Point.Longitude,
		}, "boundary")

		minDistance := dao.GetMinDistanceMeters(int(msg.GeofenceId), msg.Content.Point)
		msg.Content.MinDistance = minDistance.MinDistanceMeters
		fmt.Println("预警：", geofenceInfo)
		//在预警区域，没有进入围栏内
		if check && !onGonFence {
			alertType := 1
			alertDoc := "未知人员靠近围栏"
			alertClass := "低级"
			//检查当前用户是否长时间处于预警区, 检查缓存，不存在这写入, 时间阈值需要存储到数据库
			//逗留三分钟，警报
			now := time.Now().Unix()
			if now-node.JoinTime > 60*3 {
				alertType = 2
				alertDoc = "围栏附近存在未知人员长期滞留"
				alertClass = "中级"
			}

			//获取多边形质心
			centroid := dao.GetCentroid(msg.GeofenceId)
			fmt.Println("center:", centroid)

			centerLongitude, _ := strconv.ParseFloat(centroid.Longitude, 64)
			centerLatitude, _ := strconv.ParseFloat(centroid.Latitude, 64)
			center := model.PointFloat{
				Longitude: centerLongitude,
				Latitude:  centerLatitude,
			}

			targetLongitude, _ := strconv.ParseFloat(msg.Content.Point.Longitude, 64)
			targetLatitude, _ := strconv.ParseFloat(msg.Content.Point.Latitude, 64)
			target := model.PointFloat{
				Longitude: targetLongitude,
				Latitude:  targetLatitude,
			}

			msg.AlertType = alertType
			content := model.Content{
				Id:          int(msg.GeofenceId),
				AddressName: geofenceInfo.Name + tools.GetDirection(center, target) + "方向",
				AlertTime:   time.Now().Format("2006-01-02 15:04:05"),
				AlertClass:  alertClass,
				Point:       msg.Content.Point,
				AlertDic:    alertDoc,
				MinDistance: minDistance.MinDistanceMeters,
				State:       1,
			}
			msg.Content = content

			pointLocker.Lock()
			manyPointInfoMap[node.ClientId] = pointInfo{
				point:  target,
				geoId:  int(msg.GeofenceId),
				uid:    node.ClientId,
				isReal: true,
			}
			manyMeantimePointMap[msg.GeofenceId] = manyPointInfoMap
			pointLocker.Unlock()

			//获取当前在预警围栏的目标
			loc := make([]*redis.GeoLocation, 0)
			for key, v := range manyMeantimePointMap[msg.GeofenceId] {
				if v.isReal {
					loc = append(loc, &redis.GeoLocation{
						Name:      fmt.Sprintf("point%d", key),
						Longitude: v.point.Longitude,
						Latitude:  v.point.Latitude,
						GeoHash:   time.Now().Unix(),
					})
				}
			}

			key := "alert:" + strconv.Itoa(int(msg.GeofenceId))
			ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
			_, err = global.RDB.GeoAdd(ctx, key, loc...).Result()
			if err != nil {
				zap.S().Error("写入坐标失败", err)
				return
			}
			list, err := global.RDB.GeoRadius(ctx, key, target.Longitude, target.Latitude, &redis.GeoRadiusQuery{
				Radius:    5,    // 5米内
				Unit:      "m",  // 单位为米
				WithCoord: true, // 如果需要获取点的具体坐标可以加上这个
			}).Result()
			if err != nil {
				zap.S().Error("获取数据失败:", err)
				return
			}

			if _, err := global.RDB.Del(ctx, key).Result(); err != nil {
				zap.S().Error("移除key失败:", err)
				return
			}

			//区域内同一时间段出现三人
			if len(list) >= 3 {
				msg.AlertType = 3
				msg.Content.AlertDic = "围栏附近存在大量人员滞留"
				msg.Content.AlertClass = "高级"
			}
		} else { //当定位点走出预警区域后，需要扣除预警人员数
			pointLocker.Lock()
			manyPointInfoMap[node.ClientId] = pointInfo{
				isReal: false,
			}
			manyMeantimePointMap[msg.GeofenceId] = manyPointInfoMap
			pointLocker.Unlock()
		}

		msgStr, err := json.Marshal(msg)
		if err != nil {
			zap.S().Error("Marshal fail", msg.clientId)
			return
		}

		//预警通知
		userIds, _ := tools.ParseIntSliceFromString(geofenceInfo.ManagerIds)

		//获取用户邮件
		userDate := dao.GetUserList(userIds, []string{"id", "email", "username", "ruler"})

		for _, uid := range userIds {
			if userDate[uid].Ruler != 1 { //只发送消息至管理员
				continue
			}
			curUserNode, ok := clientMap[int64(uid)]
			if !ok && msg.AlertType != 0 { //AlertTyp：0为不预警
				zap.S().Info("客户端掉线", msg.clientId)
				email := userDate[uid].Email
				name := userDate[uid].UserName
				link := "http://www.iceymoss.top/alert/list"
				//构造预警邮件内容
				body := fmt.Sprintf("尊敬的管理员%s：\n\n您好！我们在%s监测到了以下重要预警信息，请立即查看并采取适当行动：\n\n%s别预警：%s（经纬度：%s, %s）发现不明人员长时间滞留。\n\n您可以通过以下链接访问预警详情页面:%s，以获取更多信息和可能的行动指南：\n\n点击此处访问预警详情页面\n\n请确保采取必要的预防措施保障安全。\n\n谢谢！\n\n安全监控团队", name, msg.Content.AlertTime, msg.Content.AlertClass, msg.Content.AddressName, msg.Content.Point.Longitude, msg.Content.Point.Latitude, link)
				if err := utils.SendEMail(email, body); err != nil {
					zap.S().Error("邮件推送失败", err)
				}
				continue
			}
			if msg.AlertType != 0 {
				curUserNode.DataQueue <- msgStr //推送消息至管理员界面
			}
		}

		//写入历史
		mod := model.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		lat, _ := strconv.ParseFloat(msg.Content.Point.Latitude, 64)
		lng, _ := strconv.ParseFloat(msg.Content.Point.Longitude, 64)
		ok := dao.CreateMessage(model.Messages{
			Model:       mod,
			GeoID:       uint64(msg.GeofenceId),
			ClientID:    uint64(msg.clientId),
			AddressName: msg.Content.AddressName,
			AlertTime:   msg.Content.AlertTime,
			AlertClass:  msg.Content.AlertClass,
			AlertDic:    msg.Content.AlertDic,
			MinDistance: msg.Content.MinDistance,
			State:       1,
			PointLat:    lat,
			PointLng:    lng,
			Name:        geofenceInfo.Name,
		})
		if !ok {
			zap.S().Error("写入预警数据失败")
		}
	}
}

// sendProc 从node中获取信息并写入websocket中
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Info("写入消息失败", err)
				return
			}
			fmt.Println("数据发送socket成功")
		}

	}
}

func sendMsg(uid int64, msg []byte) {
	var node *Node
	var ok bool
	if node, ok = clientMap[uid]; !ok {
		zap.S().Info("构造通讯节点失败", uid)
		return
	}
	err := node.Conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		zap.S().Error("写入消息失败", err)
		return
	}
	zap.S().Info("数据发送socket成功")
}

func CreateGeofence(geo model.GeofenceModel) error {
	//处理坐标
	uid := ""
	for _, v := range geo.ManagerIDs {
		uid += strconv.Itoa(v) + ","
	}
	//清除末尾逗号
	if len(uid) > 1 {
		uid = uid[:len(uid)-1]
	}

	//处理围栏
	boundary := ""
	for _, v := range geo.Boundary {
		boundary += fmt.Sprintf("%s %s,", v.Lng, v.Lat)
	}
	boundary = boundary[:len(boundary)-1]

	alertArea := ""
	for _, v := range geo.AlertArea {
		alertArea += fmt.Sprintf("%s %s,", v.Lng, v.Lat)
	}
	alertArea = alertArea[:len(alertArea)-1]

	if geo.CityName == "" {
		geo.CityName = "未知"
	}
	if geo.CityCoords == nil {
		geo.CityCoords = &model.Coordinate{
			Lat: "120.234",
			Lng: "31.3434",
		}
	}
	return dao.CreateGeofence(geo, boundary, alertArea, uid)
}

// GetGeofenceById 获取围栏详情
func GetGeofenceById(geoId int) model.GeofenceModel {
	return dao.GetGeofenceById(geoId)
}

func GetGeofenceList() []model.GeofenceModel {
	return dao.GetGeofenceList()
}

func UpdateGeofence(geoId int, name string, des string, state int, alertDist int, managerIds []int) error {
	return dao.UpdateGeofence(geoId, name, des, state, alertDist, managerIds)
}
