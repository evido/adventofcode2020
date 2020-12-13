package main

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readData(filename string) ([]int64, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	data := make([]int64, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			break
		}

		el, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}

		data = append(data, el)
	}

	return data, nil
}

func isValid(terms []int64, target int64) bool {
	for ix1, term1 := range terms {
		for ix2, term2 := range terms {
			if ix1 == ix2 {
				continue
			}

			if term1+term2 == target {
				return true
			}
		}
	}

	return false
}

func findInvalidElement(data []int64, preamble, context int) int {
	for ix := preamble; ix < len(data); ix += 1 {
		if !isValid(data[ix-context:ix], data[ix]) {
			return ix
		}
	}

	return -1
}

func findTerms(data []int64, target int64) []int64 {
	lower := 0
	upper := 0
	current := int64(0)

	for upper < len(data) {
		if current > target {
			current -= data[lower]
			lower += 1
		} else if current < target {
			current += data[upper]
			upper += 1
		} else {
			break
		}
	}

	return data[lower:upper]
}

func findWeakness(terms []int64) int64 {
	sort.Slice(terms, func(i, j int) bool {
		return terms[i] < terms[j]
	})

	return terms[0] + terms[len(terms)-1]
}

func main() {
	data, err := readData(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read data: %s\n", err)
	}

	invalidElementIndex := findInvalidElement(data, 25, 25)
	log.Printf("Invalid data element: %d\n", data[invalidElementIndex])

	terms := findTerms(data, data[invalidElementIndex])
	log.Printf("terms: %+v\n", terms)
	log.Printf("weakness: %d\n", findWeakness(terms))
}
