package club

import "time"

func addTime(currentTime, startTime, endTime time.Time) time.Time {
	return currentTime.Add(endTime.Sub(startTime))
}
