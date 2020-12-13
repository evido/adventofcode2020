package main

import (
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
)

const (
	COLUMNS_PER_ROW = 8
)

type BoardingPass struct {
	code string
}

func decode(code string, lc, uc rune) int {
	lower := 0
	upper := int(math.Pow(2, float64(len(code))) - 1)
	for _, c := range code {
		switch c {
		case lc:
			upper = lower + (upper-lower)/2
			break

		case uc:
			lower = lower + (upper-lower)/2 + 1
			break
		default:
			break
		}
	}

	return lower
}

func (pass *BoardingPass) decodeRow() int {
	rowCode := pass.code[:len(pass.code)-3]
	return decode(rowCode, 'F', 'B')
}

func (pass *BoardingPass) decodeColumn() int {
	columnCode := pass.code[len(pass.code)-3:]
	return decode(columnCode, 'L', 'R')
}

func (pass *BoardingPass) seatId() int {
	row := pass.decodeRow()
	column := pass.decodeColumn()

	return row*COLUMNS_PER_ROW + column
}

func readBoardingPasses(filename string) ([]BoardingPass, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	boardingPasses := make([]BoardingPass, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			continue
		}

		pass := readBoardingPass(line)
		boardingPasses = append(boardingPasses, pass)
	}

	return boardingPasses, nil
}

func readBoardingPass(line string) BoardingPass {
	return BoardingPass{
		code: line,
	}
}

func findMaxSeatId(boardingPasses []BoardingPass) int {

	maxSeatId := -1
	for _, pass := range boardingPasses {
		seatId := pass.seatId()
		if seatId > maxSeatId {
			maxSeatId = seatId
		}
	}

	return maxSeatId
}

func findMySeatId(boardingPasses []BoardingPass) int {
	allSeats := make([]bool, findMaxSeatId(boardingPasses))

	for _, pass := range boardingPasses {
		allSeats[pass.seatId()-1] = true
	}

	for ix := 1; ix < len(allSeats)-1; ix += 1 {
		if allSeats[ix-1] && !allSeats[ix] && allSeats[ix+1] {
			return ix + 1
		}
	}

	return -1
}

func main() {
	boardingPasses, err := readBoardingPasses(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read boarding passes: %s\n", err)
	}

	log.Printf("Seat ID: %d\n", findMaxSeatId(boardingPasses))
	log.Printf("My Seat ID: %d\n", findMySeatId(boardingPasses))
}
