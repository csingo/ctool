package cHelper

import (
	"fmt"
	"gitee.com/csingo/ctool/core/cHelper/constants/randomType"
	"math/rand"
	"strconv"
	"strings"
)

// IsNumber 字符串是否为数字
func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// IsInt 字符串是否为整数
func IsInt(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

// IsBool 字符串是否为布尔值
func IsBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

// ToFloat32 字符串转32位浮点数
func ToFloat32(s string) float32 {
	res, _ := strconv.ParseFloat(s, 32)
	return float32(res)
}

// ToFloat64 字符串转64位浮点数
func ToFloat64(s string) float64 {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}

// ToBool 字符串转布尔值
func ToBool(s string) bool {
	res, _ := strconv.ParseBool(s)
	return res
}

// ToInt8 字符串转int8
func ToInt8(s string) int8 {
	res, _ := strconv.ParseInt(s, 10, 8)
	return int8(res)
}

// ToInt 字符串转int
func ToInt(s string) int {
	res, _ := strconv.Atoi(s)
	return res
}

// ToInt32 字符串转int32
func ToInt32(s string) int32 {
	res, _ := strconv.ParseInt(s, 10, 32)
	return int32(res)
}

// ToInt64 字符串转int64
func ToInt64(s string) int64 {
	res, _ := strconv.ParseInt(s, 10, 64)
	return res
}

// ToUint8 字符串转uint8
func ToUint8(s string) uint8 {
	res, _ := strconv.ParseInt(s, 10, 8)
	return uint8(res)
}

// ToUint 字符串转uint
func ToUint(s string) uint {
	res, _ := strconv.Atoi(s)
	return uint(res)
}

// ToUint32 字符串转uint32
func ToUint32(s string) uint32 {
	res, _ := strconv.ParseInt(s, 10, 32)
	return uint32(res)
}

// ToUint64 字符串转uint64
func ToUint64(s string) uint64 {
	res, _ := strconv.ParseInt(s, 10, 64)
	return uint64(res)
}

// ToString 所有类型转字符串
func ToString(num interface{}) string {
	res := fmt.Sprintf("%v", num)
	return res
}

// Ucfirst 首字母大写
func Ucfirst(s string) string {
	if len(s) < 1 {
		return ""
	}

	data := []rune(s)
	if data[0] >= 97 && data[0] <= 122 {
		data[0] = data[0] - 32
	}

	return string(data)
}

func RandomStr(length int, typ randomType.Value) string {
	var base = []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", // 10
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", // 26,36
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", // 26,62
		"~", "!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "=", "-", "{", "}", "[", "]", ":", ";", "<", ">", "?", "/", ".", ",", "|", // 28,90
	}

	var res string
	var randomBase []string
	var randomBaseLen int
	switch typ {
	default:
		fallthrough
	case randomType.All:
		randomBase = base
	case randomType.Num:
		randomBase = base[0:9]
	case randomType.Char:
		randomBase = base[10:61]
	case randomType.LowerChar:
		randomBase = base[10:35]
	case randomType.UpperChar:
		randomBase = base[36:61]
	case randomType.Symbol:
		randomBase = base[62:89]
	case randomType.NumAndChar:
		randomBase = base[0:61]
	case randomType.CharAndSymbol:
		randomBase = base[10:89]
	}
	randomBaseLen = len(randomBase)

	if length < 1 {
		return ""
	}

	for i := 0; i < length; i++ {
		num := rand.Intn(randomBaseLen - 1)
		char := randomBase[num]
		res = res + char
	}

	return res
}

// ReplaceAllFromMap 替换字符串
func ReplaceAllFromMap(src string, data map[string]string) string {
	for k, v := range data {
		src = strings.ReplaceAll(src, k, v)
	}

	return src
}

// ReplaceAllIfNotContain 字符串不存在则替换字符串, 0:判断字符串, 1:被替换字符串, 2:替换字符串
func ReplaceAllIfNotContain(src string, data [][]string) string {
	for _, v := range data {
		if len(v) != 3 {
			continue
		}
		if !strings.Contains(src, v[0]) {
			src = strings.ReplaceAll(src, v[1], v[2])
		}

	}

	return src
}
