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

func intVals(inp iter.Seq[string]) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		next, stop := iter.Pull(inp)
		defer stop()
		for {
			val, ok := next()
			if !ok {
				return
			}
			fields := strings.Fields(val)
			if len(fields) >= 2 {
				iv1, _ := strconv.Atoi(fields[0])
				iv2, _ := strconv.Atoi(fields[1])
				yield(iv1, iv2)
			}
		}
	}
}

func collect2[T any](inp iter.Seq2[T, T]) ([]T, []T) {
	res1, res2 := []T{}, []T{}
	for v1, v2 := range inp {
		res1 = append(res1, v1)
		res2 = append(res2, v2)
	}
	return res1, res2
}

func sorted(inp iter.Seq2[int, int]) iter.Seq2[int, int] {
	res1, res2 := collect2(inp)
	slices.Sort(res1)
	slices.Sort(res2)
	return func(yield func(int, int) bool) {
		for i := range res1 {
			yield(res1[i], res2[i])
		}
	}
}

func part1() {
	res := 0
	for v1, v2 := range sorted(intVals(core.Lines())) {
		if v1 > v2 {
			v1, v2 = v2, v1
		}
		res += v2 - v1
	}
	fmt.Println(res)
}

func part2() {
	left, right := collect2(intVals(core.Lines()))
	counts := map[int]int{}
	for _, v := range right {
		counts[v]++
	}
	res := 0
	for _, v := range left {
		res += v * counts[v]
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
