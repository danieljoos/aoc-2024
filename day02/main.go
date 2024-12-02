package main

import (
	"flag"
	"fmt"
	"iter"
	"slices"

	"github.com/danieljoos/aoc-2024/core"
)

func part1() {
	isSafe := func(report []int) bool {
		if len(report) < 2 {
			return true
		}
		it := slices.All(report)
		if report[1] < report[0] {
			it = slices.Backward(report)
		}
		prev := -1
		for _, v := range it {
			if prev >= 0 {
				if diff := v - prev; !(diff >= 1 && diff <= 3) {
					return false
				}
			}
			prev = v
		}
		return true
	}
	safeReports := slices.Collect(core.Filter(core.IntVals(core.Lines()), isSafe))
	fmt.Println(len(safeReports))
}

func part2() {
	isSafelyIncreasing := func(report iter.Seq2[int, int], ignoreIdx int) bool {
		prev := -1
		for i, v := range report {
			if i == ignoreIdx {
				continue
			}
			if prev >= 0 {
				if diff := v - prev; !(diff >= 1 && diff <= 3) {
					return false
				}
			}
			prev = v
		}
		return true
	}
	isSafe := func(report []int) bool {
		for i := range len(report) + 1 {
			if isSafelyIncreasing(slices.All(report), i-1) ||
				isSafelyIncreasing(slices.Backward(report), i-1) {
				return true
			}
		}
		return false
	}
	safeReports := slices.Collect(core.Filter(core.IntVals(core.Lines()), isSafe))
	fmt.Println(len(safeReports))
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
