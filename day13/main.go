package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Note struct {
	departure int64
	busLines  []int64
}

func readNote(filename string) (Note, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Note{}, err
	}

	lines := strings.Split(string(bytes), "\n")
	if len(lines) < 2 {
		return Note{}, errors.New("Note should have at least 2 lines")
	}

	departure, err := strconv.ParseInt(lines[0], 10, 64)
	if err != nil {
		return Note{}, err
	}

	busLines := make([]int64, 0)
	for _, id := range strings.Split(lines[1], ",") {
		if id == "x" {
			busLines = append(busLines, 0)
			continue
		}

		busLine, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return Note{}, err
		}

		busLines = append(busLines, busLine)
	}

	return Note{
		departure: departure,
		busLines:  busLines,
	}, nil
}

type Departure struct {
	time    int64
	busLine int64
}

func findMinDeparture(departure int64, busLine int64) Departure {
	remainder := departure % busLine

	if remainder == 0 {
		remainder = busLine
	}

	return Departure{
		busLine: busLine,
		time:    departure + (busLine - remainder),
	}
}

type BusDeparture struct {
	busLine int64
	offset  int64
}

func main() {
	note, err := readNote(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read note: %s\n", err)
	}

	minDeparture := Departure{time: -1}
	for _, busLine := range note.busLines {
		if busLine == 0 {
			continue
		}

		departure := findMinDeparture(note.departure, busLine)
		if minDeparture.time == -1 || departure.time < minDeparture.time {
			minDeparture = departure
		}
	}

	log.Printf("Departure: %+v\n", minDeparture)
	log.Printf("Score: %d\n", (minDeparture.time-note.departure)*minDeparture.busLine)

	departures := make([]BusDeparture, 0)
	for offset, busLine := range note.busLines {
		if busLine != 0 {
			departures = append(departures, BusDeparture{
				offset:  int64(offset),
				busLine: busLine,
			})
		}
	}

	log.Printf("%+v\n", departures)
	timestamp := int64(0)
	inc := departures[0].busLine

	for ix := 1; ix < len(departures); ix += 1 {
		for (timestamp+departures[ix].offset)%departures[ix].busLine != 0 {
			timestamp += inc
		}

		inc *= departures[ix].busLine
	}

	log.Printf("%d\n", timestamp)
}
