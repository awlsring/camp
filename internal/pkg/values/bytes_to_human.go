package values

import (
	"fmt"
	"math"
)

func BytesToHumanReadable(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	exp := int(math.Log(float64(b)) / math.Log(unit))
	pre := "KMGTPE"[exp-1 : exp]
	val := float64(b) / math.Pow(unit, float64(exp))
	if val == float64(int64(val)) {
		return fmt.Sprintf("%d %sB", int64(val), pre)
	}
	return fmt.Sprintf("%.1f %sB", val, pre)
}
