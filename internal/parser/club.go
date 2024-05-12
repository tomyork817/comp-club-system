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

func ReadComputerClub(filename string) (club.ComputerClub, error) {
	file, err := os.Open(filename)
	if err != nil {
		return club.ComputerClub{}, NewParseError("")
	}

	scanner := bufio.NewScanner(file)

	if !scanner.Scan() {
		return club.ComputerClub{}, NewParseError("")
	}
	tablesCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return club.ComputerClub{}, NewParseError(scanner.Text())
	}

	if !scanner.Scan() {
		return club.ComputerClub{}, NewParseError(scanner.Text())
	}
	times := strings.Split(scanner.Text(), " ")
	if len(times) != 2 {
		return club.ComputerClub{}, NewParseError(scanner.Text())
	}
	startTime, err := time.Parse(club.TimeFormat, times[0])
	if err != nil {
		return club.ComputerClub{}, NewParseError(scanner.Text())
	}
	endTime, err := time.Parse(club.TimeFormat, times[1])
	if err != nil {
		return club.ComputerClub{}, NewParseError(scanner.Text())
	}

	if !scanner.Scan() {
		return club.ComputerClub{}, NewParseError(scanner.Text())
	}
	costPerHour, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return club.ComputerClub{}, NewParseError(scanner.Text())
	}

	events := make([]club.IncomingEvent, 0)
	for scanner.Scan() {
		event, err := club.ParseEvent(scanner.Text())
		if err != nil {
			return club.ComputerClub{}, NewParseError(scanner.Text())
		}
		events = append(events, event)
	}

	if err = scanner.Err(); err != nil {
		return club.ComputerClub{}, NewParseError(scanner.Text())
	}

	err = file.Close()
	if err != nil {
		return club.ComputerClub{}, err
	}

	return club.NewComputerClub(tablesCount, startTime, endTime, costPerHour, events), err
}
