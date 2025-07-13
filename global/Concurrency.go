package global

import (
	"math"
	"runtime"
)

// CalcConcurrency 动态计算并发数
func CalcConcurrency(scancount int) int {
	MinConcurrency := 10   // 最小并发数
	MaxConcurrency := 1000 // 最大并发数
	MaxForSmallJob := 30   // 小任务限制并发数

	if scancount <= 0 {
		return MinConcurrency
	}
	if scancount < MinConcurrency {
		return scancount
	}

	cpu := runtime.NumCPU()
	base := int(math.Log2(float64(scancount)) * float64(cpu) * GrowthFactor)

	// 针对小任务（如 < 500 个），上限限制
	if scancount < 500 && base > MaxForSmallJob {
		base = MaxForSmallJob
	}

	if base < MinConcurrency {
		base = MinConcurrency
	}
	if base > MaxConcurrency {
		base = MaxConcurrency
	}

	return base
}
