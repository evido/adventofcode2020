package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func read_expense_report(file_name string) ([]int, error) {
	bytes, err := ioutil.ReadFile(file_name)
	if err != nil {
		return nil, err
	}

	expenses := make([]int, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		expense, err := strconv.ParseInt(line, 10, 32)
		if err == nil {
			expenses = append(expenses, int(expense))
		}
	}

	return expenses, nil
}

func find_product_(goal, terms int, expenses map[int]int) int {
	if terms == 1 {
		if count, ok := expenses[goal]; ok && count > 0 {
			return goal
		}
		return -1
	}

	for expense, count := range expenses {
		if count > 0 && expense < goal {
			expenses[expense] -= 1
			solution := find_product_(goal-expense, terms-1, expenses)
			if solution >= 0 {
				return expense * solution
			}
			expenses[expense] += 1
		}
	}

	return -1
}

func find_product(goal, terms int, expenses []int) int {
	expenseSet := make(map[int]int)
	for _, expense := range expenses {
		if count, ok := expenseSet[expense]; ok {
			expenseSet[expense] = count + 1
		} else {
			expenseSet[expense] = 1
		}
	}

	return find_product_(goal, terms, expenseSet)
}

func main() {

	expenses, err := read_expense_report(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read expense report!\n")
	}

	fmt.Printf("Solution: %d\n", find_product(2020, 3, expenses))
}
