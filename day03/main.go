package main

import (
	"flag"
	"fmt"
	"iter"
	"regexp"
	"strconv"

	"github.com/danieljoos/aoc-2024/core"
)

type instructionKind int

const (
	instructionKind_undefined instructionKind = iota
	instructionKind_mul
	instructionKind_do
	instructionKind_dont
)

type instruction struct {
	kind               instructionKind
	operandA, operandB int
}

func parse(input iter.Seq[string]) iter.Seq[instruction] {
	pat := regexp.MustCompile(`(don't)|(do)|mul\((\d+),(\d+)\)`)
	return func(yield func(instruction) bool) {
		next, stop := iter.Pull(input)
		defer stop()
		for {
			v, ok := next()
			if !ok {
				return
			}
			for _, groups := range pat.FindAllStringSubmatch(v, -1) {
				switch {
				case groups[1] != "":
					yield(instruction{kind: instructionKind_dont})
				case groups[2] != "":
					yield(instruction{kind: instructionKind_do})
				default:
					operandA, _ := strconv.Atoi(groups[3])
					operandB, _ := strconv.Atoi(groups[4])
					yield(instruction{
						kind:     instructionKind_mul,
						operandA: operandA,
						operandB: operandB,
					})
				}
			}
		}
	}
}

func part1() {
	res := 0
	for i := range parse(core.Lines()) {
		if i.kind == instructionKind_mul {
			res += i.operandA * i.operandB
		}
	}
	fmt.Println(res)
}

func part2() {
	res := 0
	enabled := true
	for i := range parse(core.Lines()) {
		switch i.kind {
		case instructionKind_do:
			enabled = true
		case instructionKind_dont:
			enabled = false
		case instructionKind_mul:
			if enabled {
				res += i.operandA * i.operandB
			}
		}
	}
	fmt.Println(res)
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
