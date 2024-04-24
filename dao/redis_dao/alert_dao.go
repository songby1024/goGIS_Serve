package redis_dao

//import (
//	"github.com/redis/go-redis/v9"
//	"go-gis-api/internal/common"
//	"go-gis-api/internal/global"
//	"go.uber.org/zap"
//)
//
//type AlertGeoServer struct {
//}
//
//func (d *AlertGeoServer) CreateGeo(key string, geoList []*redis.GeoLocation) error {
//	ctx, cancel := common.NewContext()
//	defer cancel()
//
//	_, err := global.RDB.GeoAdd(ctx, key, geoList...).Result()
//	if err != nil {
//		zap.S().Error(err)
//		return err
//	}
//	return nil
//}
//
//func (d *AlertGeoServer) GetPointDistant(key, geo []string) {
//	//ctx, cancel := common.NewContext()
//	//defer cancel()
//	//
//	//global.RDB.GeoDist(ctx, key, )
//}
