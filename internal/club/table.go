package club

import "time"

type tableState = string

const (
	empty tableState = "EMPTY"
	busy  tableState = "BUSY"
)

type table struct {
	id          int
	startTime   time.Time
	state       tableState
	client      *client
	overallTime time.Time
	gain        int
}
