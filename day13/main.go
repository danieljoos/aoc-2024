package main

import (
	"flag"
	"fmt"
	"iter"
	"regexp"
	"strconv"

	"github.com/danieljoos/aoc-2024/core"
)

type Position struct{ X, Y int64 }
type Machine struct{ ButtonA, ButtonB, Prize Position }

func readMachines() iter.Seq[*Machine] {
	return func(yield func(*Machine) bool) {
		pattern := regexp.MustCompile(`^(.+): X[+=](\d+), Y[+=](\d+)$`)
		curr := &Machine{}
		for line := range core.Lines() {
			matches := pattern.FindAllStringSubmatch(line, 1)
			if len(matches) != 1 {
				continue
			}
			m := matches[0]
			x, _ := strconv.ParseInt(m[2], 10, 64)
			y, _ := strconv.ParseInt(m[3], 10, 64)
			switch m[1] {
			case "Button A":
				curr.ButtonA = Position{x, y}
			case "Button B":
				curr.ButtonB = Position{x, y}
			case "Prize":
				curr.Prize = Position{x, y}
				if !yield(curr) {
					return
				}
				curr = &Machine{}
			}
		}
	}
}

func (m *Machine) Costs(limit int64) int64 {
	b := ((m.ButtonA.X * m.Prize.Y) - (m.ButtonA.Y * m.Prize.X)) / ((m.ButtonA.X * m.ButtonB.Y) - (m.ButtonA.Y * m.ButtonB.X))
	a := (m.Prize.X - (b * m.ButtonB.X)) / m.ButtonA.X
	testPos := Position{
		X: (a * m.ButtonA.X) + (b * m.ButtonB.X),
		Y: (a * m.ButtonA.Y) + (b * m.ButtonB.Y),
	}
	if testPos != m.Prize || (limit > 0 && (a > limit || b > limit)) {
		return -1
	}
	cost := (a * 3) + b
	return cost
}

func part1() {
	sum := int64(0)
	for m := range readMachines() {
		if costs := m.Costs(100); costs != -1 {
			sum += costs
		}
	}
	fmt.Println(sum)
}

func part2() {
	sum := int64(0)
	for m := range readMachines() {
		m.Prize.X += 10000000000000
		m.Prize.Y += 10000000000000
		if costs := m.Costs(-1); costs != -1 {
			sum += costs
		}
	}
	fmt.Println(sum)
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
