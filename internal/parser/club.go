package parser

import (
	"bufio"
	"comp-club-system/internal/club"
	"os"
	"strconv"
	"strings"
	"time"
)

type ParseError struct {
	message string
}

func (e ParseError) Error() string {
	return e.message
}

func NewParseError(message string) ParseError {
	return ParseError{message: message}
}

func ReadComputerClub(filename string) (*club.ComputerClub, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return nil, NewParseError("")
	}
	tablesCount, err := strconv.Atoi(scanner.Text())
	if err != nil || tablesCount <= 0 {
		return nil, NewParseError(scanner.Text())
	}

	if !scanner.Scan() {
		return nil, NewParseError(scanner.Text())
	}
	times := strings.Split(scanner.Text(), " ")
	if len(times) != 2 {
		return nil, NewParseError(scanner.Text())
	}
	startTime, err := time.Parse(club.TimeFormat, times[0])
	if err != nil {
		return nil, NewParseError(scanner.Text())
	}
	endTime, err := time.Parse(club.TimeFormat, times[1])
	if err != nil || !endTime.After(startTime) {
		return nil, NewParseError(scanner.Text())
	}

	if !scanner.Scan() {
		return nil, NewParseError(scanner.Text())
	}
	costPerHour, err := strconv.Atoi(scanner.Text())
	if err != nil || costPerHour <= 0 {
		return nil, NewParseError(scanner.Text())
	}

	info := club.NewInfo(tablesCount, startTime, endTime, costPerHour)

	events := make([]club.IncomingEvent, 0)
	prevTime, _ := time.Parse(club.TimeFormat, "00:00")
	for scanner.Scan() {
		event, err := club.ParseEvent(scanner.Text(), info, prevTime)
		if err != nil {
			return nil, NewParseError(scanner.Text())
		}
		prevTime = event.Time()
		events = append(events, event)
	}

	if err = scanner.Err(); err != nil {
		return nil, NewParseError(scanner.Text())
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return club.NewComputerClub(info, events), err
}
