package cHelper

import "math/rand"

// RandomInt 生成随机整数
func RandomInt(min, max int) int {
	if min > max {
		min = 0
	}

	diff := max - min
	num := rand.Intn(diff)

	return min + num
}
