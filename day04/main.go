package main

import (
	"flag"
	"fmt"

	"github.com/danieljoos/aoc-2024/core"
)

type Field [][]rune

func readField() Field {
	result := Field{}
	y := 0
	for line := range core.Lines() {
		result = append(result, []rune{})
		for _, c := range line {
			result[y] = append(result[y], c)
		}
		y++
	}
	return result
}

func part1() {
	field := readField()
	count := 0
	xmasLen := len("XMAS")
	isXMAS := func(vals ...rune) bool {
		return len(vals) == 4 && vals[0] == 'X' && vals[1] == 'M' && vals[2] == 'A' && vals[3] == 'S'
	}

	for _, row := range field {
		// LR
		for x := range len(row) - xmasLen + 1 {
			if isXMAS(row[x], row[x+1], row[x+2], row[x+3]) {
				count++
			}
		}

		// RL
		for x := len(row) - 1; x >= xmasLen-1; x-- {
			if isXMAS(row[x], row[x-1], row[x-2], row[x-3]) {
				count++
			}
		}
	}

	for x := range len(field[0]) {
		// TB
		for y := range len(field) - xmasLen + 1 {
			if isXMAS(field[y][x], field[y+1][x], field[y+2][x], field[y+3][x]) {
				count++
			}
		}

		// BT
		for y := len(field) - 1; y >= xmasLen-1; y-- {
			if isXMAS(field[y][x], field[y-1][x], field[y-2][x], field[y-3][x]) {
				count++
			}
		}
	}

	for y := range len(field) - xmasLen + 1 {
		// TLBR
		for x := range len(field[y]) - xmasLen + 1 {
			if isXMAS(field[y][x], field[y+1][x+1], field[y+2][x+2], field[y+3][x+3]) {
				count++
			}
		}

		// TRBL
		for x := len(field[y]) - 1; x >= 3; x-- {
			if isXMAS(field[y][x], field[y+1][x-1], field[y+2][x-2], field[y+3][x-3]) {
				count++
			}
		}
	}

	for y := len(field) - 1; y >= 3; y-- {
		// BLTR
		for x := range len(field[y]) - xmasLen + 1 {
			if isXMAS(field[y][x], field[y-1][x+1], field[y-2][x+2], field[y-3][x+3]) {
				count++
			}
		}

		// BRTL
		for x := len(field[y]) - 1; x >= 3; x-- {
			if isXMAS(field[y][x], field[y-1][x-1], field[y-2][x-2], field[y-3][x-3]) {
				count++
			}
		}
	}

	fmt.Println(count)
}

func part2() {
	field := readField()
	count := 0
	masLen := len("MAS")
	isMAS := func(vals ...rune) bool {
		return len(vals) == 3 && vals[0] == 'M' && vals[1] == 'A' && vals[2] == 'S'
	}

	for y := range len(field) - masLen + 1 {
		for x := range len(field[y]) - masLen + 1 {
			diagTLBR := isMAS(field[y][x], field[y+1][x+1], field[y+2][x+2])
			diagBLTR := isMAS(field[y+2][x], field[y+1][x+1], field[y][x+2])
			diagTRBL := isMAS(field[y][x+2], field[y+1][x+1], field[y+2][x])
			diagBRTL := isMAS(field[y+2][x+2], field[y+1][x+1], field[y][x])
			if (diagTLBR || diagBRTL) && (diagBLTR || diagTRBL) {
				count++
			}
		}
	}

	fmt.Println(count)
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
