package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type GroupResponse struct {
	answers []string
}

func readGroupResponse(lines []string) GroupResponse {
	return GroupResponse{
		answers: lines,
	}
}

func readGroupResponses(filename string) ([]GroupResponse, error) {

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	responses := make([]GroupResponse, 0)
	responseLines := make([]string, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			responses = append(responses, readGroupResponse(responseLines))
			responseLines = make([]string, 0)
		} else {
			responseLines = append(responseLines, line)
		}
	}

	return responses, nil
}

func (response *GroupResponse) Positive() []string {
	questions := make(map[rune]bool)

	for _, response := range response.answers {
		for _, question := range response {
			questions[question] = true
		}
	}

	questionList := make([]string, 0)
	for question := range questions {
		questionList = append(questionList, string(question))
	}

	return questionList
}

func sumPositiveResponsesByGroup(responses []GroupResponse) int {
	sum := 0

	for _, response := range responses {
		sum += len(response.Positive())
	}

	return sum
}

func main() {
	responses, err := readGroupResponses(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read responses: %s\n", err)
	}

	sumByGroup := sumPositiveResponsesByGroup(responses)
	log.Printf("Positive Response By Group: %d\n", sumByGroup)
}
