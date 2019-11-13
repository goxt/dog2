package util

import (
	"math"
)

/**
 * 保留两位小数
 */
func Round(f float64) float64 {
	return math.Round(f*100) / 100
}
