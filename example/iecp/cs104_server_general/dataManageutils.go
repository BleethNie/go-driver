package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	YxAddress = 1
	YcAddress = 16385
	YkAddress = 24577
)

var YCMap = make(map[int]float32)
var YXMap = make(map[int]int)
var YCTest [length]int
var YXTest [length]int
var AppMap = map[string][length]int{}

const length = 1000

func InitRandomData() {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < length; i++ {
		YCMap[YcAddress+i] = rand.Float32() * 100
		//YCMap[YcAddress+i] = 12.35
	}

	for i := 0; i < length; i++ {
		val := rand.Float32()
		if val > 0.5 {
			YXMap[YxAddress+i] = 0
		} else {
			YXMap[YxAddress+i] = 1
		}
	}

	for i := 0; i < length; i++ {
		YCTest[i] = rand.Intn(length) + YcAddress
	}
	AppMap["yc"] = YCTest

	for i := 0; i < length; i++ {
		YXTest[i] = rand.Intn(length) + YxAddress
	}
	AppMap["yx"] = YXTest

}
func PrintMap() {

	for key, value := range YXTest {
		fmt.Printf("id = %d : ID:=%d  value = %d\n", key, YXTest[key], YXMap[value])
	}
}
