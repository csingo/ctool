package cHelper

import (
	"sort"
)

// InArrayInt 判断是否在数组中，int
func InArrayInt(data int, array []int) bool {
	sort.Ints(array)
	index := sort.SearchInts(array, data)
	if index < len(array) && array[index] == data {
		return true
	}
	return false
}

// InArrayString 判断是否在数组中，string
func InArrayString(data string, array []string) bool {
	sort.Strings(array)
	index := sort.SearchStrings(array, data)
	if index < len(array) && array[index] == data {
		return true
	}
	return false
}

// InArrayFloat64 判断是否在数组中，float64
func InArrayFloat64(data float64, array []float64) bool {
	sort.Float64s(array)
	index := sort.SearchFloat64s(array, data)
	if index < len(array) && array[index] == data {
		return true
	}
	return false
}
