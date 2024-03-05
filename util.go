package main

import (
	"math"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func interfaceToFloat(num interface{}) float64 {
	switch v := num.(type) {
	case bool:
		if v {
			return 1
		}
		return 0
	case float32:
		return float64(v)
	case float64:
		return v
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		feetFloat, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return 0
		} else {
			return feetFloat
		}
	default:
		return math.NaN()
	}
}

func interfaceToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	default:
		fVal := interfaceToFloat(val)
		return strconv.FormatFloat(fVal, 'f', -1, 64)
	}
}

func interfaceToBool(val interface{}) bool {
	switch v := val.(type) {
	case bool:
		return v
	default:
		fVal := interfaceToString(val)
		return fVal != ""
	}
}

func Xrange(start, end int) []int {
	arr := make([]int, 0, end-start)
	for i := start; i < end; i++ {
		arr = append(arr, i)
	}
	return arr
}

func MergeMaps(maps ...map[string]EvalFunc) map[string]EvalFunc {
	retMap := make(map[string]EvalFunc, 0)
	for _, m := range maps {
		for k, v := range m {
			retMap[k] = v
		}
	}
	return retMap
}

func roundTo2DecimalPlaces(v float64) float64 {
	return roundToNDecimalPlaces(v, 2)
}

func roundToNDecimalPlaces(v float64, n int) float64 {
	p := math.Pow10(n)
	return math.Floor(v*p) / p
}

func MinV[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}
func MaxV[T constraints.Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}
