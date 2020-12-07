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

func main() {

	goal := 2020
	expenses, err := read_expense_report(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read expense report!\n")
	}

	expenseSet := make(map[int]int)
	for _, expense := range expenses {
		if count, ok := expenseSet[expense]; ok {
			expenseSet[expense] = count + 1
		} else {
			expenseSet[expense] = 1
		}
	}

	for expense, expenseCount := range expenseSet {
		targetExpense := goal - expense
		if targetExpense == expense && expenseCount > 1 {
			fmt.Printf("%d\n", expense*targetExpense)
			break
		}
		if _, ok := expenseSet[targetExpense]; ok {
			fmt.Printf("%d\n", expense*targetExpense)
			break
		}
	}
}
