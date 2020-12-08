package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Policy struct {
	MinCount int
	MaxCount int
	Char     byte
}

type PasswordEntry struct {
	Policy   Policy
	Password string
}

func readEntries(filename string) ([]PasswordEntry, error) {
	entries := make([]PasswordEntry, 0)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			continue
		}

		entry, err := readEntry(line)
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func readEntry(line string) (PasswordEntry, error) {
	lineParts := strings.Split(line, ":")
	var entry PasswordEntry
	if len(lineParts) != 2 {
		return entry, errors.New("entry should be formatted as <policy>: <password>")
	}

	entry.Password = string(lineParts[1][1:])

	policyParts := strings.Split(lineParts[0], " ")
	if len(policyParts) != 2 {
		return entry, errors.New("policy should be <min>-<max> <char>")
	}

	entry.Policy.Char = policyParts[1][0]
	policyCounts := strings.Split(policyParts[0], "-")
	if len(policyCounts) != 2 {
		return entry, errors.New("policy should be <min>-max> <char>")
	}

	minCount, err := strconv.ParseInt(policyCounts[0], 10, 32)
	if err != nil {
		return entry, err
	}

	maxCount, err := strconv.ParseInt(policyCounts[1], 10, 32)
	if err != nil {
		return entry, err
	}

	entry.Policy.MinCount = int(minCount)
	entry.Policy.MaxCount = int(maxCount)

	return entry, nil
}

func countValidEntries(entries []PasswordEntry, isValid func(PasswordEntry) bool) int {
	count := 0
	for _, entry := range entries {
		if isValid(entry) {
			count += 1
		}
	}
	return count
}

func IsValid(entry PasswordEntry) bool {
	charCount := 0
	for _, b := range entry.Password {
		if byte(b) == entry.Policy.Char {
			charCount += 1
		}
	}

	return charCount >= entry.Policy.MinCount && charCount <= entry.Policy.MaxCount
}

func IsValidUpdated(entry PasswordEntry) bool {

	var validFirstPhase bool
	if len(entry.Password) >= entry.Policy.MinCount {
		validFirstPhase = byte(entry.Password[entry.Policy.MinCount-1]) == entry.Policy.Char
	} else {
		validFirstPhase = false
	}

	var validSecondPhase bool
	if len(entry.Password) >= entry.Policy.MaxCount {
		validSecondPhase = byte(entry.Password[entry.Policy.MaxCount-1]) == entry.Policy.Char
	} else {
		validSecondPhase = false
	}

	return validFirstPhase != validSecondPhase
}

func main() {
	entries, err := readEntries("custom_input.txt")
	if err != nil {
		log.Fatalf("Invalid test input!: %s\n", err)
	}

	fmt.Printf("Valid entries: %d\n", countValidEntries(entries, IsValidUpdated))
}
