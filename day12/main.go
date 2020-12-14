package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type ActionCode int

const (
	NORTH ActionCode = iota
	SOUTH
	EAST
	WEST
	LEFT
	RIGHT
	FORWARD
)

type Action struct {
	code  ActionCode
	value int
}

func readActions(filename string) ([]Action, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	actions := make([]Action, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			break
		}

		action, err := readAction(line)
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func readAction(line string) (Action, error) {

	code := line[0]
	action := Action{}

	switch code {
	case 'N':
		action.code = NORTH
		break
	case 'S':
		action.code = SOUTH
		break
	case 'E':
		action.code = EAST
		break
	case 'W':
		action.code = WEST
		break
	case 'L':
		action.code = LEFT
		break
	case 'R':
		action.code = RIGHT
		break
	case 'F':
		action.code = FORWARD
		break
	default:
		return action, fmt.Errorf("Unknown error code: %b\n", code)
	}

	value, err := strconv.ParseInt(line[1:], 10, 32)
	if err != nil {
		return action, err
	}
	action.value = int(value)

	return action, nil
}

type Navigation struct {
	direction [2]int
	position  [2]int
}

func (navigation *Navigation) ApplyAction(action *Action) {
	direction := [2]int{0, 0}
	switch action.code {
	case NORTH:
		direction[1] = 1
		break
	case SOUTH:
		direction[1] = -1
		break
	case EAST:
		direction[0] = 1
		break
	case WEST:
		direction[0] = -1
		break
	case LEFT:
		for ix := 0; ix < action.value/90; ix += 1 {
			newXMagnitude := -navigation.direction[1]
			navigation.direction[1] = navigation.direction[0]
			navigation.direction[0] = newXMagnitude
		}
		break
	case RIGHT:
		for ix := 0; ix < action.value/90; ix += 1 {
			newYMagnitude := -navigation.direction[0]
			navigation.direction[0] = navigation.direction[1]
			navigation.direction[1] = newYMagnitude
		}
		break
	case FORWARD:
		direction = navigation.direction
		break
	}

	navigation.position[0] += action.value * direction[0]
	navigation.position[1] += action.value * direction[1]
}

func (navigation *Navigation) Distance() int {
	return int(math.Abs(float64(navigation.position[0])) +
		math.Abs(float64(navigation.position[1])))
}

func NewNavigation() Navigation {
	return Navigation{
		direction: [2]int{1, 0},
		position:  [2]int{0, 0},
	}
}

func main() {
	actions, err := readActions(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read actions: %s\n", err)
	}

	navigation := NewNavigation()
	for _, action := range actions {
		navigation.ApplyAction(&action)
		log.Printf("Navigation: %+v\n", navigation)
	}

	log.Printf("Navigation: %+v\n", navigation)
	log.Printf("Distance: %d\n", navigation.Distance())
}
