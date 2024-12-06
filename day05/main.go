package main

import (
	"flag"
	"fmt"
	"iter"
	"slices"
	"strconv"
	"strings"

	"github.com/danieljoos/aoc-2024/core"
)

type PageOrderingRule [2]int
type Update []int

type Input struct {
	Rules   []PageOrderingRule
	Updates []Update
}

func readInput() *Input {
	result := &Input{}
	section := 0
	for line := range core.Lines() {
		if line == "" {
			section++
			continue
		}
		switch section {
		case 0:
			strs := strings.SplitN(line, "|", 2)
			p1, _ := strconv.Atoi(strs[0])
			p2, _ := strconv.Atoi(strs[1])
			result.Rules = append(result.Rules, PageOrderingRule{p1, p2})
		case 1:
			update := Update{}
			for _, str := range strings.Split(line, ",") {
				p, _ := strconv.Atoi(str)
				update = append(update, p)
			}
			result.Updates = append(result.Updates, update)
		}
	}
	return result
}

func (i *Input) isInRightOrder(update Update) bool {
	for _, r := range i.Rules {
		i0 := slices.Index(update, r[0])
		if i0 < 0 {
			continue
		}
		i1 := slices.Index(update, r[1])
		if i1 < 0 {
			continue
		}
		if !(i0 < i1) {
			return false
		}
	}
	return true
}

func (i *Input) SortFunc() func(a, b int) int {
	return func(a, b int) int {
		idx := slices.IndexFunc(i.Rules, func(r PageOrderingRule) bool {
			return (r[0] == a && r[1] == b) || (r[0] == b && r[1] == a)
		})
		if idx >= 0 {
			r := i.Rules[idx]
			if r[0] == a {
				return -1
			} else {
				return 1
			}
		}
		return 0
	}
}

func midSum(updates iter.Seq[Update]) int {
	sum := 0
	for u := range updates {
		sum += u[len(u)/2]
	}
	return sum
}

func part1() {
	input := readInput()
	sum := midSum(core.Filter(slices.Values(input.Updates), input.isInRightOrder))
	fmt.Println(sum)
}

func part2() {
	input := readInput()
	isNotInRightOrder := func(update Update) bool {
		return !input.isInRightOrder(update)
	}
	sortedUpdates := func(vals iter.Seq[Update]) iter.Seq[Update] {
		sf := input.SortFunc()
		return func(yield func(Update) bool) {
			next, stop := iter.Pull(vals)
			defer stop()
			for {
				v, ok := next()
				if !ok {
					return
				}
				sorted := slices.Clone(v)
				slices.SortFunc(sorted, sf)
				yield(sorted)
			}
		}
	}
	sum := midSum(sortedUpdates(core.Filter(slices.Values(input.Updates), isNotInRightOrder)))
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
