package main

import (
	"flag"
	"fmt"
	"iter"
	"slices"

	"github.com/danieljoos/aoc-2024/core"
)

func isSafelyIncreasing(report iter.Seq2[int, int], ignoreIdx int) bool {
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

func part1() {
	isSafe := func(report []int) bool {
		return isSafelyIncreasing(slices.All(report), -1) ||
			isSafelyIncreasing(slices.Backward(report), -1)
	}
	safeReports := slices.Collect(core.Filter(core.IntVals(core.Lines()), isSafe))
	fmt.Println(len(safeReports))
}

func part2() {
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
