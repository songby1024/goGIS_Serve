package dao

import (
	"fmt"
	"go-gis-api/internal/models"
	"testing"
)

func TestCheckPointInAlertRange(t *testing.T) {
	isContant := CheckPointInAlertRange(1, models.Point{
		Longitude: "-74.0060",
		Latitude:  "41.7128",
	})
	fmt.Println(isContant)
}
