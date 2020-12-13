package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Machine struct {
	accumulator         int
	instruction_pointer int
}

type Instruction interface {
	process(machine *Machine)
}

type NoOperation struct {
	argument int
}

func (operation NoOperation) process(machine *Machine) {
	machine.instruction_pointer += 1
}

type Jump struct {
	argument int
}

func (operation Jump) process(machine *Machine) {
	machine.instruction_pointer += operation.argument
}

type Accumulate struct {
	argument int
}

func (operation Accumulate) process(machine *Machine) {
	machine.accumulator += operation.argument
	machine.instruction_pointer += 1
}

func readInstructions(filename string) ([]Instruction, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	instructions := make([]Instruction, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		if len(line) == 0 {
			break
		}

		instruction, err := readInstruction(line)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, instruction)
	}

	return instructions, nil
}

func readInstruction(line string) (Instruction, error) {
	parts := strings.Split(line, " ")

	argument, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return nil, err
	}

	switch parts[0] {
	case "nop":
		return NoOperation{
			argument: int(argument),
		}, nil
	case "acc":
		return Accumulate{
			argument: int(argument),
		}, nil
	case "jmp":
		return Jump{
			argument: int(argument),
		}, nil
	}

	return nil, fmt.Errorf("Unrecognized instruction: %s", line)
}

func findLoop(machine *Machine, instructions []Instruction) bool {
	visited := make(map[int]bool)
	current := machine.instruction_pointer

	for current < len(instructions) {
		if _, ok := visited[current]; !ok {
			visited[current] = true
			instructions[current].process(machine)
			current = machine.instruction_pointer
		} else {
			return true
		}
	}

	return false
}

func (machine *Machine) Reset() {
	machine.instruction_pointer = 0
	machine.accumulator = 0
}

func swapInstruction(instruction Instruction) Instruction {
	switch instruction.(type) {
	case NoOperation:
		return Jump{
			argument: instruction.(NoOperation).argument,
		}
	case Jump:
		return NoOperation{
			argument: instruction.(Jump).argument,
		}
	default:
		return instruction
	}
}

func tryFixLoop(machine *Machine, instructions []Instruction) int {
	for ix, original := range instructions {
		machine.Reset()

		instructions[ix] = swapInstruction(original)
		loops := findLoop(machine, instructions)
		instructions[ix] = original

		if !loops {
			return ix
		}
	}

	return -1
}

func main() {
	instructions, err := readInstructions(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read instructions: %s\n", err)
	}

	machine := Machine{}
	findLoop(&machine, instructions)
	fmt.Printf("Accumulator: %d\n", machine.accumulator)

	machine.Reset()
	fmt.Printf("Fixed instruction: %d\n", tryFixLoop(&machine, instructions))
	fmt.Printf("Fixed accumulator: %d\n", machine.accumulator)
}
