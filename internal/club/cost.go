package club

import "time"

func countCost(startTime, endTime time.Time, costPerHour int) int {
	hours := endTime.Hour() - startTime.Hour()
	if endTime.Minute() != startTime.Minute() {
		hours++
	}
	return hours * costPerHour
}
