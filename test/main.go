package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func test(num []int, currentPage int, pageSize int) []int {
	// var res []int
	if currentPage == 1 {
		return num[0:pageSize]
	}

	// return num[currentPage*(len(res)/pageSize)-1 : currentPage*(len(res)/pageSize)-1+pageSize]
	return num[(currentPage-1)*pageSize : (currentPage-1)*pageSize+pageSize]
}

func main() {
	// num := []int{1, 2, 3, 4, 5, 6}
	// currentPage := 2
	// pageSize := 2
	// // 1,2,3  4,5,6
	// fmt.Println(test(num, currentPage, pageSize))
	str := "1,2"
	data := strings.Split(str, ",")
	var data2 []int
	for _, x := range data {
		temp, _ := strconv.Atoi(x)
		data2 = append(data2, temp)
	}
	fmt.Println(strings.Split(str, ","), reflect.TypeOf(data))
	fmt.Println(data2, reflect.TypeOf(data2))
}
