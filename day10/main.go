package main

import (
	"flag"
	"fmt"

	"github.com/danieljoos/aoc-2024/core"
)

type Topomap [][]int
type Position struct{ X, Y int }
type PositionSet map[Position]struct{}

var neighboring = [][2]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}

func (tm Topomap) Width() int {
	return len(tm[0])
}

func (tm Topomap) Height() int {
	return len(tm)
}

func readTopomap() Topomap {
	result := Topomap{}
	for line := range core.Lines() {
		row := make([]int, 0, len(line))
		for _, c := range line {
			row = append(row, int(c-'0'))
		}
		result = append(result, row)
	}
	return result
}

func followTrail(tm Topomap, pos Position, peaks PositionSet) int {
	width := tm.Width()
	height := tm.Height()
	value := tm[pos.Y][pos.X]
	rating := 0
	for _, n := range neighboring {
		p := Position{X: pos.X + n[0], Y: pos.Y + n[1]}
		if p.X < 0 || p.Y < 0 || p.X >= width || p.Y >= height {
			continue
		}
		if val := tm[p.Y][p.X]; val == value+1 {
			if val < 9 {
				rating += followTrail(tm, p, peaks)
				continue
			}
			rating++
			if _, ok := peaks[p]; !ok {
				peaks[p] = struct{}{}
			}
		}
	}
	return rating
}

func part1() {
	tm := readTopomap()
	result := 0
	for y, row := range tm {
		for x, val := range row {
			if val != 0 {
				continue
			}
			peaks := map[Position]struct{}{}
			followTrail(tm, Position{x, y}, peaks)
			result += len(peaks)
		}
	}
	fmt.Println(result)
}

func part2() {
	tm := readTopomap()
	result := 0
	for y, row := range tm {
		for x, val := range row {
			if val != 0 {
				continue
			}
			result += followTrail(tm, Position{x, y}, map[Position]struct{}{})
		}
	}
	fmt.Println(result)
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
