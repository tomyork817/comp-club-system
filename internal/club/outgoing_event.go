package club

import (
	"fmt"
	"time"
)

type OutgoingEvent interface {
	fmt.Stringer
}

const (
	clientGoneOutEventId       = 11
	clientSatAtTableOutEventId = 12
	errorEventId               = 13
)

type errorMessage = string

const (
	youShallNotPass  errorMessage = "YouShallNotPass"
	notOpenYet       errorMessage = "NotOpenYet"
	placeIsBusy      errorMessage = "PlaceIsBusy"
	clientUnknown    errorMessage = "ClientUnknown"
	iCanWaitNoLonger errorMessage = "ICanWaitNoLonger!"
)

type ClientGoneOutEvent struct {
	time time.Time
	name string
}

func (e ClientGoneOutEvent) String() string {
	return fmt.Sprintf("%s %d %s", e.time.Format(TimeFormat), clientGoneOutEventId, e.name)
}

type ClientSatAtTheTableOutEvent struct {
	time  time.Time
	name  string
	table int
}

func (e ClientSatAtTheTableOutEvent) String() string {
	return fmt.Sprintf("%s %d %s %d", e.time.Format(TimeFormat), clientGoneOutEventId, e.name, e.table)
}

type ErrorEvent struct {
	time    time.Time
	message errorMessage
}

func (e ErrorEvent) String() string {
	return fmt.Sprintf("%s %d %s", e.time.Format(TimeFormat), errorEventId, e.message)
}
