package main

import "fmt"

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
