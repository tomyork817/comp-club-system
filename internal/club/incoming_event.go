package club

import (
	"fmt"
	"time"
)

type IncomingEvent interface {
	OutgoingEvent
	execute(info info, clubState *clubState) OutgoingEvent
}

const (
	clientCameEventId         = 1
	clientSatAtTableInEventId = 2
	clientWaitingEventId      = 3
	clientGoneInEventId       = 4
)

type EmptyEvent struct{}

func (e EmptyEvent) String() string {
	return ""
}

func (e EmptyEvent) execute(info info, state *clubState) OutgoingEvent {
	return e
}

type ClientCameEvent struct {
	time time.Time
	name string
}

func (e ClientCameEvent) execute(info info, clubState *clubState) OutgoingEvent {
	_, ok := clubState.clients[e.name]
	if ok {
		return ErrorEvent{
			time:    e.time,
			message: youShallNotPass,
		}
	}
	if info.startTime.After(e.time) || info.endTime.Before(e.time) {
		return ErrorEvent{
			time:    e.time,
			message: notOpenYet,
		}
	}
	clubState.clients[e.name] = &client{
		name:  e.name,
		state: inClub,
		table: nil,
	}
	return EmptyEvent{}
}

func (e ClientCameEvent) String() string {
	return fmt.Sprintf("%s %d %s", e.time.Format(TimeFormat), clientCameEventId, e.name)
}

type ClientSatAtTheTableInEvent struct {
	time  time.Time
	name  string
	table int
}

func (e ClientSatAtTheTableInEvent) execute(info info, clubState *clubState) OutgoingEvent {
	_, ok := clubState.clients[e.name]
	if !ok {
		return ErrorEvent{
			time:    e.time,
			message: clientUnknown,
		}
	}

	client := clubState.clients[e.name]
	switch client.state {
	case gone:
		return ErrorEvent{
			time:    e.time,
			message: clientUnknown,
		}
	case atTable:
		if e.table == client.table.id {
			return ErrorEvent{
				time:    e.time,
				message: placeIsBusy,
			}
		}
		if clubState.tables[e.table].state == busy {
			return ErrorEvent{
				time:    e.time,
				message: placeIsBusy,
			}
		}
		client := clubState.clients[e.name]
		table := client.table
		table.state = empty
		table.client = nil
		clubState.overallCost += countCost(table.startTime, e.time, info.costPerHour)

		client.table = clubState.tables[e.table]
		client.table.client = client
		client.table.startTime = e.time
		client.table.state = busy

		return EmptyEvent{}
	case inClub:
		if clubState.tables[e.table].state == busy {
			return ErrorEvent{
				time:    e.time,
				message: placeIsBusy,
			}
		}
		client := clubState.clients[e.name]
		client.table = clubState.tables[e.table]
		client.table.client = client
		client.table.startTime = e.time
		client.table.state = busy

		return EmptyEvent{}
	default:
		return EmptyEvent{}
	}
}

func (e ClientSatAtTheTableInEvent) String() string {
	return fmt.Sprintf("%s %d %s %d", e.time.Format(TimeFormat), clientSatAtTableInEventId, e.name, e.table)
}

type ClientWaitingEvent struct {
	time time.Time
	name string
}

func (e ClientWaitingEvent) execute(info info, clubState *clubState) OutgoingEvent {
	_, ok := clubState.clients[e.name]
	if !ok {
		return ErrorEvent{
			time:    e.time,
			message: clientUnknown,
		}
	}

	client := clubState.clients[e.name]
	switch client.state {
	case inClub:
		for _, t := range clubState.tables {
			if t.state == empty {
				return ErrorEvent{
					time:    e.time,
					message: iCanWaitNoLonger,
				}
			}
		}
		if len(clubState.waitQueue) > info.tablesCount {
			client.state = gone
			return ClientGoneOutEvent{
				time: e.time,
				name: client.name,
			}
		}
		client.state = waiting
		clubState.waitQueue = append(clubState.waitQueue, client)
		return EmptyEvent{}
	default:
		return EmptyEvent{}
	}
}

func (e ClientWaitingEvent) String() string {
	return fmt.Sprintf("%s %d %s", e.time.Format(TimeFormat), clientWaitingEventId, e.name)
}

type ClientGoneInEvent struct {
	time time.Time
	name string
}

func (e ClientGoneInEvent) execute(info info, clubState *clubState) OutgoingEvent {
	_, ok := clubState.clients[e.name]
	if !ok {
		return ErrorEvent{
			time:    e.time,
			message: clientUnknown,
		}
	}

	client := clubState.clients[e.name]
	switch client.state {
	case atTable:
		table := client.table
		client.table = nil
		client.state = gone

		table.client = nil
		table.state = empty
		clubState.overallCost += countCost(table.startTime, e.time, info.costPerHour)
		if len(clubState.waitQueue) == 0 {
			return EmptyEvent{}
		}

		client = clubState.waitQueue[0]
		clubState.waitQueue = clubState.waitQueue[1:]
		table.startTime = e.time
		table.state = busy
		table.client = client

		return ClientSatAtTheTableOutEvent{
			time:  e.time,
			name:  client.name,
			table: table.id,
		}
	default:
		return EmptyEvent{}
	}
}

func (e ClientGoneInEvent) String() string {
	return fmt.Sprintf("%s %d %s", e.time.Format(TimeFormat), clientGoneInEventId, e.name)
}

type ClientGoneEndTimeEvent struct {
	name string
}

func (e ClientGoneEndTimeEvent) execute(info info, clubState *clubState) OutgoingEvent {
	_, ok := clubState.clients[e.name]
	if !ok {
		return ErrorEvent{
			time:    info.endTime,
			message: clientUnknown,
		}
	}

	client := clubState.clients[e.name]
	switch client.state {
	case atTable:
		table := client.table
		client.table = nil
		client.state = gone

		table.client = nil
		table.state = empty
		clubState.overallCost += countCost(table.startTime, info.endTime, info.costPerHour)

		return ClientGoneOutEvent{
			time: info.endTime,
			name: client.name,
		}
	case inClub:
		client.state = gone
		return ClientGoneOutEvent{
			time: info.endTime,
			name: client.name,
		}
	case waiting:
		client.state = gone
		return ClientGoneOutEvent{
			time: info.endTime,
			name: client.name,
		}
	default:
		return EmptyEvent{}
	}
}

func (e ClientGoneEndTimeEvent) String() string {
	return ""
}
