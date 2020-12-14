package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type State int

const (
	FLOOR State = iota
	EMPTY
	OCCUPIED
)

type Board struct {
	state [][]State
}

func (board *Board) CountOccupiedNeighboursv2(i, j int) int {
	directions := [][]int{
		{0, 1},
		{1, 1},
		{1, 0},
		{1, -1},
		{0, -1},
		{-1, -1},
		{-1, 0},
		{-1, 1},
	}

	occupied := 0
	for _, direction := range directions {
		if board.SeesOccupiedSeat(i, j, direction) {
			occupied += 1
		}
	}

	return occupied
}

func (board *Board) SeesOccupiedSeat(i, j int, direction []int) bool {
	for n := 1; ; n += 1 {
		t_i := i + n*direction[0]
		t_j := j + n*direction[1]

		if t_i < 0 || t_i >= len(board.state) {
			break
		}

		if t_j < 0 || t_j >= len(board.state[t_i]) {
			break
		}

		switch board.state[t_i][t_j] {
		case FLOOR:
			continue
		case EMPTY:
			return false
		case OCCUPIED:
			return true
		}
	}

	return false
}

func (board *Board) CountOccupiedNeighbours(i, j int) int {
	occupied := 0
	for di := -1; di <= 1; di += 1 {
		if i+di < 0 || i+di >= len(board.state) {
			continue
		}

		for dj := -1; dj <= 1; dj += 1 {
			if di == 0 && dj == 0 {
				continue
			}

			if j+dj < 0 || j+dj >= len(board.state[i]) {
				continue
			}

			if board.IsOccupied(i+di, j+dj) {
				occupied += 1
			}
		}
	}
	return occupied
}

func (board *Board) SimulateSeat(i, j int) State {
	occupied := board.CountOccupiedNeighbours(i, j)

	switch board.state[i][j] {
	case FLOOR:
		return FLOOR
	case EMPTY:
		if occupied > 0 {
			return EMPTY
		}
		return OCCUPIED
	case OCCUPIED:
		if occupied > 3 {
			return EMPTY
		}
		return OCCUPIED
	default:
		return board.state[i][j]
	}
}

func (board *Board) SimulateSeatv2(i, j int) State {
	occupied := board.CountOccupiedNeighboursv2(i, j)

	switch board.state[i][j] {
	case FLOOR:
		return FLOOR
	case EMPTY:
		if occupied > 0 {
			return EMPTY
		}
		return OCCUPIED
	case OCCUPIED:
		if occupied > 4 {
			return EMPTY
		}
		return OCCUPIED
	default:
		return board.state[i][j]
	}
}

func (board *Board) Simulate() int {

	modifications := 0
	newState := make([][]State, len(board.state))
	for i := 0; i < len(board.state); i += 1 {
		newState[i] = make([]State, len(board.state[i]))
		for j := 0; j < len(board.state[i]); j += 1 {
			newState[i][j] = board.SimulateSeatv2(i, j)
			if newState[i][j] != board.state[i][j] {
				modifications += 1
			}
		}
	}

	board.state = newState
	return modifications
}

func (board *Board) IsOccupied(i, j int) bool {
	return board.state[i][j] == OCCUPIED
}

func readBoard(filename string) (Board, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Board{}, err
	}

	state := make([][]State, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			break
		}

		row := make([]State, len(line))
		for ix, col := range line {
			switch col {
			case 'L':
				row[ix] = EMPTY
				break
			case '.':
				row[ix] = FLOOR
				break
			default:
				return Board{}, errors.New("Unknown state code")
			}
		}

		state = append(state, row)
	}

	return Board{
		state: state,
	}, nil
}

func countOccupiedSeats(board *Board) int {
	count := 0
	for i := 0; i < len(board.state); i += 1 {
		for j := 0; j < len(board.state[i]); j += 1 {
			if board.IsOccupied(i, j) {
				count += 1
			}
		}
	}
	return count
}

func main() {
	board, err := readBoard(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read board: %s\n", err)
	}

	for board.Simulate() > 0 {
	}
	log.Printf("Occupied seats: %d\n", countOccupiedSeats(&board))
}
