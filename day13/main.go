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
	departure int
	busLines  []int
}

func readNote(filename string) (Note, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Note{}, err
	}

	busLines := make([]int, 0)
	lines := strings.Split(string(bytes), "\n")
	if len(lines) < 2 {
		return Note{}, errors.New("Note should have at least 2 lines")
	}

	departure, err := strconv.ParseInt(lines[0], 10, 32)
	if err != nil {
		return Note{}, err
	}

	for _, id := range strings.Split(lines[1], ",") {
		if id == "x" {
			continue
		}

		busLine, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			return Note{}, err
		}

		busLines = append(busLines, int(busLine))
	}

	return Note{
		departure: int(departure),
		busLines:  busLines,
	}, nil
}

type Departure struct {
	time    int
	busLine int
}

func findMinDeparture(departure, busLine int) Departure {
	remainder := departure % busLine
	return Departure{
		busLine: busLine,
		time:    departure + (busLine - remainder),
	}
}

func main() {
	note, err := readNote(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read note: %s\n", err)
	}

	minDeparture := findMinDeparture(note.departure, note.busLines[0])
	for ix := 1; ix < len(note.busLines); ix += 1 {
		departure := findMinDeparture(note.departure, note.busLines[ix])
		if departure.time < minDeparture.time {
			minDeparture = departure
		}
	}

	log.Printf("Departure: %+v\n", minDeparture)
	log.Printf("Score: %d\n", (minDeparture.time-note.departure)*minDeparture.busLine)
}
