package club

import (
	"fmt"
	"time"
)

const (
	TimeFormat = "15:04"
)

type ComputerClub struct {
	info      info
	clubState *clubState

	events         []IncomingEvent
	outgoingEvents []OutgoingEvent
}

type info struct {
	tablesCount int
	startTime   time.Time
	endTime     time.Time
	costPerHour int
}

type clubState struct {
	tables      map[int]*table
	clients     map[string]*client
	waitQueue   []*client
	overallCost int
}

// NewComputerClub is a constructor of ComputerClub
func NewComputerClub(tablesCount int, startTime, endTime time.Time, costPerHour int, events []IncomingEvent) ComputerClub {
	info := info{tablesCount, startTime, endTime, costPerHour}
	tables := make(map[int]*table, tablesCount)
	for i := 1; i <= tablesCount; i++ {
		tables[i] = &table{
			id:        i,
			startTime: time.Time{},
			state:     empty,
			client:    nil,
		}
	}
	clients := make(map[string]*client)
	waitQueue := make([]*client, 0)
	outgoingEvents := make([]OutgoingEvent, 0)
	state := clubState{tables, clients, waitQueue, 0}
	return ComputerClub{
		info:           info,
		clubState:      &state,
		events:         events,
		outgoingEvents: outgoingEvents,
	}
}

func (c *ComputerClub) RunIncomingEvents() {
	for _, event := range c.events {
		outEvent := event.execute(c.info, c.clubState)
		c.outgoingEvents = append(c.outgoingEvents, event)
		switch outEvent.(type) {
		case EmptyEvent:
			continue
		default:
			c.outgoingEvents = append(c.outgoingEvents, outEvent)
		}
	}
	for _, client := range c.clubState.clients {
		event := ClientGoneEndTimeEvent{
			name: client.name,
		}
		outEvent := event.execute(c.info, c.clubState)
		switch outEvent.(type) {
		case EmptyEvent:
			continue
		default:
			c.outgoingEvents = append(c.outgoingEvents, outEvent)
		}
	}
}

func (c *ComputerClub) Print() {
	fmt.Println(c.info.startTime.Format(TimeFormat))
	for _, event := range c.outgoingEvents {
		fmt.Println(event.String())
	}
	fmt.Println(c.info.endTime.Format(TimeFormat))
}
