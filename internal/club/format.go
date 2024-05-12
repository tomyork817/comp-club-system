package club

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type FormatError struct {
	message string
}

func NewFormatError(message string) FormatError {
	return FormatError{message: message}
}

func (e FormatError) Error() string {
	return e.message
}

const (
	clientNamePattern = `\w+`
)

func ParseEvent(line string) (IncomingEvent, error) {
	words := strings.Split(line, " ")
	if len(words) < 3 || len(words) > 4 {
		return EmptyEvent{}, NewFormatError("invalid format of event")
	}
	id, err := strconv.Atoi(words[1])
	if err != nil {
		return EmptyEvent{}, NewFormatError("invalid type of event id")
	}
	t, err := time.Parse(TimeFormat, words[0])
	if err != nil {
		return EmptyEvent{}, NewFormatError("invalid time format")
	}
	name := words[2]
	match, err := regexp.MatchString(clientNamePattern, name)
	if !match {
		return EmptyEvent{}, NewFormatError("invalid client name format")
	}
	switch id {
	case clientCameEventId:
		return ClientCameEvent{
			time: t,
			name: name,
		}, nil
	case clientSatAtTableInEventId:
		table, err := strconv.Atoi(words[3])
		if err != nil {
			return EmptyEvent{}, NewFormatError("invalid table number format")
		}
		return ClientSatAtTheTableInEvent{
			time:  t,
			name:  name,
			table: table,
		}, nil
	case clientWaitingEventId:
		return ClientWaitingEvent{
			time: t,
			name: name,
		}, nil
	case clientGoneInEventId:
		return ClientGoneInEvent{
			time: t,
			name: name,
		}, nil
	default:
		return EmptyEvent{}, NewFormatError("invalid id of event")
	}
}
