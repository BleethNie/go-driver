package main

import (
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/cpu"
	"testing"
)

func TestCPURun(t *testing.T) {
	// 获取CPU温度
	temps, err := cpu.Info()
	if err != nil {
		fmt.Printf("无法获取CPU温度：%s\n", err)
		return
	}

	// 获取计算机ID（可选）
	id, err := machineid.ID()
	if err != nil {
		fmt.Printf("无法获取计算机ID：%s\n", err)
		return
	}

	// 打印CPU温度
	fmt.Printf("计算机ID：%s\n", id)
	for _, temp := range temps {
		fmt.Printf("CPU温度：%.2f°C\n", temp)
	}
}
