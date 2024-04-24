package tools

import (
	"fmt"
	"math"
	"serve/model"
	"strconv"
	"strings"
)

// GetDirection 返回两个坐标点之间的方向
func GetDirection(center, target model.PointFloat) string {
	// 计算目标点相对于中心点的经度和纬度差

	latDiff := target.Latitude - center.Latitude
	longDiff := target.Longitude - center.Longitude

	// 使用反正切计算角度，转换为度数
	angle := math.Atan2(longDiff, latDiff) * (180 / math.Pi)

	// 根据角度判断方位
	switch {
	case angle >= -22.5 && angle < 22.5:
		return "北"
	case angle >= 22.5 && angle < 67.5:
		return "东北"
	case angle >= 67.5 && angle < 112.5:
		return "东"
	case angle >= 112.5 && angle < 157.5:
		return "东南"
	case angle >= 157.5 || angle < -157.5:
		return "南"
	case angle >= -157.5 && angle < -112.5:
		return "西南"
	case angle >= -112.5 && angle < -67.5:
		return "西"
	case angle >= -67.5 && angle < -22.5:
		return "西北"
	default:
		return "未知"
	}
}

// ParseIntSliceFromString takes a string in format "{n1,n2,n3,...}" and returns a slice of integers.
func ParseIntSliceFromString(str string) ([]int, error) {
	// 移除字符串中的花括号
	trimmedStr := strings.Trim(str, "{}")

	// 使用逗号分割字符串
	stringSlice := strings.Split(trimmedStr, ",")

	// 创建一个切片用于存储转换后的整数
	var intSlice []int

	// 遍历字符串切片，并将每个字符串转换为整数
	for _, s := range stringSlice {
		num, err := strconv.Atoi(strings.TrimSpace(s)) // 转换前移除任何空白字符
		if err != nil {
			return nil, fmt.Errorf("failed to convert '%v' to integer: %v", s, err)
		}
		intSlice = append(intSlice, num)
	}

	return intSlice, nil
}
