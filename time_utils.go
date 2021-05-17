package main

import (
	"fmt"
	"time"
)

// MillsToHumanText 毫秒美化
func MillsToHumanText(val int64) string {
	if val < 1000 {
		return fmt.Sprintf("%v毫秒", val)
	}

	valF := float64(val) / 1000.0 // 获取秒数
	if valF < 60.0 {
		return fmt.Sprintf("%.2f秒", valF)
	}

	valF /= 60.0 // 获取分钟数
	if valF < 60.0 {
		return fmt.Sprintf("%.2f分钟", valF)
	}

	valF /= 60.0 // 获取小时数
	if valF < 24.0 {
		return fmt.Sprintf("%.2f小时", valF)
	}

	valF /= 24.0 // 获取天数
	return fmt.Sprintf("%.2f天", valF)
}

const (
	// TimeZoneShangHai 上海时区
	TimeZoneShangHai = "Asia/Shanghai"
	// TimeZoneTokyo 日本时区
	TimeZoneTokyo = "Asia/Tokyo"
	// TimeZoneNewYork 纽约时区
	TimeZoneNewYork = "America/New_York"
	// TimeZoneLosAngeles 洛杉矶时区
	TimeZoneLosAngeles = "America/Los_Angeles"
	// TimeZoneChicago 芝加哥时区
	TimeZoneChicago = "America/Chicago"
	// TimeZoneLondon 伦敦时区
	TimeZoneLondon = "Europe/London"
)

// GetTimeNow 获取当前时间
func GetTimeNow(locationStr string) string {
	now := time.Now()
	location, err := time.LoadLocation(locationStr)
	if err != nil {
		return now.Format("2006-01-02 15:04:05")
	}

	return now.In(location).Format("2006-01-02 15:04:05")
}
