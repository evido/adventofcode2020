package main

import (
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readAdapters(filename string) ([]int, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	adapters := make([]int, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			break
		}

		adapter, err := strconv.ParseInt(line, 10, 32)
		if err != nil {
			return nil, err
		}

		adapters = append(adapters, int(adapter))
	}

	return adapters, nil
}

func findAdapters(adapters []int, source, target int) []int {
	if target-source <= 3 {
		return []int{source, target}
	}

	path := make([]int, 0)
	path = append(path, source)
	for _, adapter := range adapters {
		if adapter > source && adapter-source <= 3 {
			remaining := findAdapters(adapters, adapter, target)
			if len(remaining) != 0 {
				path = append(path, remaining...)
				return path
			}
		}
	}

	return []int{}
}

func countDiff(path []int, diff int) int {
	count := 0
	for ix := 1; ix < len(path); ix += 1 {
		if path[ix]-path[ix-1] == diff {
			count += 1
		}
	}
	return count
}

var optionsCache map[int]int64 = make(map[int]int64)

func countAdaptersOptions(adapters []int, source, target int) int64 {
	if options, ok := optionsCache[source]; ok {
		return options
	}

	if target-source == 3 {
		return 1
	}

	options := int64(0)
	for _, adapter := range adapters {
		if adapter > source && adapter-source <= 3 {
			options += countAdaptersOptions(adapters, adapter, target)
		}
	}

	optionsCache[source] = options
	return options
}
func main() {
	allAdapters, err := readAdapters(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read adapters: %s\n", err)
	}

	sort.Ints(allAdapters)
	target := allAdapters[len(allAdapters)-1] + 3
	adapters := findAdapters(allAdapters, 0, target)

	log.Printf("Adapters: %+v\n", allAdapters)
	log.Printf("Path: %+v\n", adapters)

	oneDiffCount := countDiff(adapters, 1)
	threeDiffCount := countDiff(adapters, 3)

	log.Printf("Diff 1: %d\n", oneDiffCount)
	log.Printf("Diff 3: %d\n", threeDiffCount)
	log.Printf("Solution: %d\n", oneDiffCount*threeDiffCount)
	log.Printf("Options: %d\n", countAdaptersOptions(adapters, 0, target))
}
