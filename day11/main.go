package main

import (
	"flag"
	"fmt"
	"slices"

	"github.com/danieljoos/aoc-2024/core"
)

func readInput() []int {
	return core.First(core.IntVals(core.Lines()))
}

func getDigits(v int) []int {
	digits := make([]int, 0, 10)
	for {
		digits = append(digits, v%10)
		if v < 10 {
			break
		}
		v = v / 10
	}
	slices.Reverse(digits)
	return digits
}

func getValue(digits []int) int {
	value := 0
	for _, digit := range digits {
		value *= 10
		value += digit
	}
	return value
}

type CacheKey struct{ s, rem int }
type Cache map[CacheKey]int

func countStonesRecursive(s int, rem int, cache Cache) int {
	if rem <= 0 {
		return 1
	}
	if v, ok := cache[CacheKey{s, rem}]; ok {
		return v
	}
	result := 0
	if s == 0 {
		result += countStonesRecursive(1, rem-1, cache)
	} else if digits := getDigits(s); len(digits)%2 == 0 {
		result += countStonesRecursive(getValue(digits[:len(digits)/2]), rem-1, cache)
		result += countStonesRecursive(getValue(digits[len(digits)/2:]), rem-1, cache)
	} else {
		result += countStonesRecursive(s*2024, rem-1, cache)
	}
	cache[CacheKey{s, rem}] = result
	return result
}

func countStones(maxDepth int) int {
	cache := Cache{}
	sum := 0
	for _, s := range readInput() {
		sum += countStonesRecursive(s, maxDepth, cache)
	}
	return sum
}

func part1() {
	fmt.Println(countStones(25))
}

func part2() {
	fmt.Println(countStones(75))
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
