package club

import (
	"math"
	"time"
)

func countCost(startTime, endTime time.Time, costPerHour int) int {
	dur := endTime.Sub(startTime)
	hours := int(math.Ceil(dur.Hours()))
	return hours * costPerHour
}
