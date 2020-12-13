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

func (response *GroupResponse) Responses() map[string]int {
	questions := make(map[string]int)

	for _, response := range response.answers {
		for _, question := range response {
			questions[string(question)] += 1
		}
	}

	return questions
}

func (response *GroupResponse) Size() int {
	return len(response.answers)
}

func (response *GroupResponse) UnanymousResponseCount() int {
	unanymousCount := 0
	size := response.Size()

	for _, count := range response.Responses() {
		if count == size {
			unanymousCount += 1
		}
	}

	return unanymousCount
}

func (response *GroupResponse) PositiveResponseCount() int {
	return len(response.Responses())
}

func sumPositiveResponsesByGroup(responses []GroupResponse) int {
	sum := 0

	for _, response := range responses {
		sum += response.PositiveResponseCount()
	}

	return sum
}

func sumUnanymousPositiveResponseByGroup(responses []GroupResponse) int {
	sum := 0
	for _, response := range responses {
		sum += response.UnanymousResponseCount()
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

	unanymousByGroup := sumUnanymousPositiveResponseByGroup(responses)
	log.Printf("Unanymous responses by group: %d\n", unanymousByGroup)
}
