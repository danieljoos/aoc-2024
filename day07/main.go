package main

import (
	"flag"
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/danieljoos/aoc-2024/core"
)

type Operator int

const (
	OperatorAdd Operator = iota
	OperatorMultiply
	OperatorConcat
)

type Equation struct {
	Result int
	Values []int
}

func parseEquations(vals iter.Seq[string]) iter.Seq[*Equation] {
	return func(yield func(*Equation) bool) {
		next, stop := iter.Pull(vals)
		defer stop()
		for {
			line, ok := next()
			if !ok {
				return
			}
			eq := &Equation{}
			parts := strings.SplitN(line, ":", 2)
			eq.Result, _ = strconv.Atoi(parts[0])
			fields := strings.Fields(parts[1])
			eq.Values = make([]int, 0, len(fields))
			for _, f := range fields {
				i, _ := strconv.Atoi(f)
				eq.Values = append(eq.Values, i)
			}
			if !yield(eq) {
				return
			}
		}
	}
}

func operatorPermutations(num int, availableOps []Operator) iter.Seq[[]Operator] {
	return func(yield func([]Operator) bool) {
		if num == 0 {
			yield(nil)
			return
		}
		for p := range operatorPermutations(num-1, availableOps) {
			for _, op := range availableOps {
				if !yield(append(append(make([]Operator, 0, num), op), p...)) {
					return
				}
			}
		}
	}
}

func isValidEquation(eq *Equation, availableOps []Operator) bool {
	numOps := len(eq.Values) - 1
	for ops := range operatorPermutations(numOps, availableOps) {
		res := eq.Values[0]
		for i, op := range ops {
			switch op {
			case OperatorAdd:
				res += eq.Values[i+1]
			case OperatorMultiply:
				res *= eq.Values[i+1]
			case OperatorConcat:
				res, _ = strconv.Atoi(strconv.Itoa(res) + strconv.Itoa(eq.Values[i+1]))
			}
		}
		if res == eq.Result {
			return true
		}
	}
	return false
}

func sum(eqs iter.Seq[*Equation]) int64 {
	s := int64(0)
	for e := range eqs {
		s += int64(e.Result)
	}
	return s
}

func part1() {
	availableOps := []Operator{OperatorAdd, OperatorMultiply}
	isValid := func(eq *Equation) bool { return isValidEquation(eq, availableOps) }
	fmt.Println(sum(core.Filter(parseEquations(core.Lines()), isValid)))
}

func part2() {
	availableOps := []Operator{OperatorAdd, OperatorMultiply, OperatorConcat}
	isValid := func(eq *Equation) bool { return isValidEquation(eq, availableOps) }
	fmt.Println(sum(core.Filter(parseEquations(core.Lines()), isValid)))
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
