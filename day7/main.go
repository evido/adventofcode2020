package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Specification struct {
	color    string
	contents map[string]int
}

func readSpecifications(filename string) ([]Specification, error) {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	specifications := make([]Specification, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			continue
		}

		specification, err := readSpecification(line)
		if err != nil {
			return specifications, err
		}

		specifications = append(specifications, specification)
	}

	return specifications, nil
}

func readSpecification(line string) (Specification, error) {
	words := strings.Split(line, " ")
	target := strings.Join(words[:2], " ")

	specification := Specification{
		color:    target,
		contents: make(map[string]int),
	}
	for ix := 4; ix < len(words); ix += 4 {
		if strings.Join(words[ix:ix+3], " ") == "no other bags." {
			break
		}

		count, err := strconv.ParseInt(words[ix], 10, 32)
		if err != nil {
			return specification, err
		}

		contentColor := strings.Join(words[ix+1:ix+3], " ")
		specification.contents[contentColor] = int(count)
	}

	return specification, nil
}

func canContainBag(specifications []Specification, current string, target string) bool {
	for _, specification := range specifications {
		if specification.color != current {
			continue
		}

		for content := range specification.contents {
			if content == target {
				return true
			}

			if canContainBag(specifications, content, target) {
				return true
			}
		}
	}

	return false
}

func countPossibleBags(specifications []Specification, target string) int {
	possible := 0

	for _, specification := range specifications {
		if canContainBag(specifications, specification.color, target) {
			possible += 1
		}
	}

	return possible
}

func countRequiredBags(specifications []Specification, target string) int {
	total := 0
	for _, specification := range specifications {
		if specification.color != target {
			continue
		}

		for content, count := range specification.contents {
			total += count * (1 + countRequiredBags(specifications, content))
		}
	}
	return total
}

func main() {
	specifications, err := readSpecifications(os.Args[1])

	if err != nil {
		log.Fatalf("Unable to read specifications: %s\n", err)
	}

	log.Printf("Possible bags: %d\n", countPossibleBags(specifications, "shiny gold"))
	log.Printf("Required bags: %d\n", countRequiredBags(specifications, "shiny gold"))
}
