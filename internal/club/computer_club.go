package club

import (
	"fmt"
	"time"
)

const (
	TimeFormat = "15:04"
)

type ComputerClub struct {
	info      Info
	clubState *clubState

	events         []IncomingEvent
	outgoingEvents []OutgoingEvent
}

type Info struct {
	tablesCount int
	startTime   time.Time
	endTime     time.Time
	costPerHour int
}

func NewInfo(tablesCount int, startTime, endTime time.Time, costPerHour int) Info {
	return Info{
		tablesCount: tablesCount,
		startTime:   startTime,
		endTime:     endTime,
		costPerHour: costPerHour,
	}
}

type clubState struct {
	tables    map[int]*table
	clients   map[string]*client
	waitQueue []*client
}

func NewComputerClub(info Info, events []IncomingEvent) *ComputerClub {
	tables := make(map[int]*table, info.tablesCount)
	for i := 1; i <= info.tablesCount; i++ {
		tables[i] = &table{
			id:          i,
			startTime:   time.Time{},
			state:       empty,
			client:      nil,
			overallTime: time.Time{},
			gain:        0,
		}
	}
	clients := make(map[string]*client)
	waitQueue := make([]*client, 0)
	outgoingEvents := make([]OutgoingEvent, 0)
	state := clubState{tables, clients, waitQueue}
	return &ComputerClub{
		info:           info,
		clubState:      &state,
		events:         events,
		outgoingEvents: outgoingEvents,
	}
}

func (c *ComputerClub) UpdateEvents(events []IncomingEvent) {
	var newEvents []IncomingEvent
	copy(newEvents, events)
	c.events = newEvents
}

func (c *ComputerClub) RunIncomingEvents() {
	c.outgoingEvents = make([]OutgoingEvent, 0)
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
	for i := 1; i <= c.info.tablesCount; i++ {
		fmt.Printf("%d %d %s\n", i, c.clubState.tables[i].gain, c.clubState.tables[i].overallTime.Format(TimeFormat))
	}
}
