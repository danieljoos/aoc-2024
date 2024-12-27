package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	"github.com/danieljoos/aoc-2024/core"
)

const (
	OpcodeADV = iota
	OpcodeBXL
	OpcodeBST
	OpcodeJNZ
	OpcodeBXC
	OpcodeOUT
	OpcodeBDV
	OpcodeCDV
)

var ErrInvalidOperand = errors.New("invalid operand")

type Computer struct {
	RegisterA          int
	RegisterB          int
	RegisterC          int
	Program            []int
	InstructionPointer int
}

func readComputer() *Computer {
	result := &Computer{}
	lines := slices.Collect(core.Lines())
	result.RegisterA, _ = strconv.Atoi(strings.SplitN(lines[0], ": ", 2)[1])
	result.RegisterB, _ = strconv.Atoi(strings.SplitN(lines[1], ": ", 2)[1])
	result.RegisterC, _ = strconv.Atoi(strings.SplitN(lines[2], ": ", 2)[1])
	for _, p := range strings.Split(strings.SplitN(lines[4], ": ", 2)[1], ",") {
		val, _ := strconv.Atoi(p)
		result.Program = append(result.Program, val)
	}
	return result
}

func (c *Computer) Run() (string, error) {
	outs := []string{}
	for c.InstructionPointer < len(c.Program) {
		opcode := c.Program[c.InstructionPointer]
		operand := c.Program[c.InstructionPointer+1]
		switch opcode {
		case OpcodeADV:
			combo, err := c.Combo(operand)
			if err != nil {
				return "", err
			}
			val := 1 << combo
			c.RegisterA = c.RegisterA / val
		case OpcodeBXL:
			val := c.RegisterB ^ operand
			c.RegisterB = val
		case OpcodeBST:
			combo, err := c.Combo(operand)
			if err != nil {
				return "", err
			}
			c.RegisterB = combo & 0b111
		case OpcodeJNZ:
			if c.RegisterA != 0 {
				c.InstructionPointer = operand
				continue
			}
		case OpcodeBXC:
			c.RegisterB = c.RegisterB ^ c.RegisterC
		case OpcodeOUT:
			combo, err := c.Combo(operand)
			if err != nil {
				return "", err
			}
			combo = combo & 0b111
			outs = append(outs, strconv.Itoa(combo))
		case OpcodeBDV:
			combo, err := c.Combo(operand)
			if err != nil {
				return "", err
			}
			val := 1 << combo
			c.RegisterB = c.RegisterA / val
		case OpcodeCDV:
			combo, err := c.Combo(operand)
			if err != nil {
				return "", err
			}
			val := 1 << combo
			c.RegisterC = c.RegisterA / val
		}
		c.InstructionPointer += 2
	}
	return strings.Join(outs, ","), nil
}

func (c *Computer) Combo(operand int) (int, error) {
	if operand >= 0 && operand <= 3 {
		return operand, nil
	}
	switch operand {
	case 4:
		return c.RegisterA, nil
	case 5:
		return c.RegisterB, nil
	case 6:
		return c.RegisterC, nil
	}
	return 0, fmt.Errorf("%w: %d", ErrInvalidOperand, operand)
}

func part1() {
	comp := readComputer()
	fmt.Println(comp.Run())
}

func part2() {
	comp := readComputer()
	programStr := strings.Join(slices.Collect(core.StrVals(slices.Values(comp.Program))), ",")

	runMe := func(a int) string {
		cl := *comp
		cl.RegisterA = a
		res, err := cl.Run()
		if err != nil {
			log.Fatal(err)
		}
		return res
	}

	// The given program adds a digit to the beginning of the output.
	// So we can check the digits of the output starting at the end.
	// There are 8 possible values for each output digit, covering the 3 bits.
	// It's a 8-base system (vs 10 base).

	solutions := []int{0}
	for i := range slices.Backward(comp.Program) {
		nextSolutions := []int{}
		part := programStr[i*2:]      // part of the output to compare to.
		for _, s := range solutions { // solutions generate the output digits correctly up to this point.
			for a := s * 8; a < s*8+8; a++ { // try all 8 possible values for the next digit.
				res := runMe(a)
				if res == part {
					nextSolutions = append(nextSolutions, a)
				}
			}
		}
		solutions = nextSolutions
	}

	fmt.Println(slices.Min(solutions))
}

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "part1":
		part1()
	case "part2":
		part2()
	}
}
