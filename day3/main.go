package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

type Field struct {
	Template [][]bool
}

func readField(filename string) (Field, error) {
	var field Field

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return field, err
	}

	lines := strings.Split(string(bytes), "\n")
	if len(lines) == 0 {
		return field, errors.New("Input file should not be empty")
	}

	fieldWidth := len(lines[0])
	if fieldWidth == 0 {
		return field, errors.New("Input file rows should have length greater than 0")
	}
	field.Template = make([][]bool, len(lines)-1)
	for rowIndex, line := range lines {
		if len(line) == 0 {
			break
		}

		if len(line) != fieldWidth {
			return field, fmt.Errorf("All lines should have equal width: %s", line)
		}

		field.Template[rowIndex] = make([]bool, fieldWidth)
		for columnIndex, c := range line {
			switch c {
			case '.':
				field.Template[rowIndex][columnIndex] = false
				break
			case '#':
				field.Template[rowIndex][columnIndex] = true
				break
			default:
				return field, errors.New("Invalid field template character")
			}
		}
	}

	return field, nil
}

func (field *Field) HasTree(x, y int) bool {
	return field.Template[y][x%len(field.Template[0])]
}

func countTrees(field Field, dx, dy int) int {
	x := 0
	y := 0

	trees := 0
	for y < len(field.Template) {

		if field.HasTree(x, y) {
			trees += 1
		}

		x += dx
		y += dy
	}

	return trees
}

func main() {
	field, err := readField("custom_input.txt")
	if err != nil {
		log.Fatalf("Unable to read field: %s\n", err)
	}

	slopes := [][]int{
		[]int{1, 1},
		[]int{3, 1},
		[]int{5, 1},
		[]int{7, 1},
		[]int{1, 2},
	}

	total := 1
	for _, slope := range slopes {
		total *= countTrees(field, slope[0], slope[1])
	}

	fmt.Printf("Total: %d\n", total)
}
