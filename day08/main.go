package main

import (
	"flag"
	"fmt"
	"maps"
	"slices"

	"github.com/danieljoos/aoc-2024/core"
)

type Position struct{ X, Y int }

type CityMap struct {
	Width, Height int
	Antennas      map[rune][]Position
}

func (m *CityMap) IsOnMap(p Position) bool {
	return p.X >= 0 && p.X < m.Width && p.Y >= 0 && p.Y < m.Height
}

func readMap() *CityMap {
	result := &CityMap{Antennas: map[rune][]Position{}}
	for line := range core.Lines() {
		for x, c := range line {
			if c == '.' {
				continue
			}
			result.Antennas[c] = append(result.Antennas[c], Position{x, result.Height})
		}
		result.Width = len(line)
		result.Height++
	}
	return result
}

func part1() {
	antinodes := map[Position]struct{}{}
	cityMap := readMap()
	for _, antennas := range cityMap.Antennas {
		for i, a1 := range antennas {
			for _, a2 := range antennas[i+1:] {
				diffX := a2.X - a1.X
				diffY := a2.Y - a1.Y
				antinodes[Position{X: a2.X + diffX, Y: a2.Y + diffY}] = struct{}{}
				antinodes[Position{X: a1.X - diffX, Y: a1.Y - diffY}] = struct{}{}
			}
		}
	}
	an := slices.Collect(core.Filter(maps.Keys(antinodes), cityMap.IsOnMap))
	fmt.Println(len(an))
}

func part2() {
	antinodes := map[Position]struct{}{}
	cityMap := readMap()
	fillAntiNodes := func(origin Position, diffX, diffY int) {
		antinodes[origin] = struct{}{}
		for {
			p := Position{X: origin.X + diffX, Y: origin.Y + diffY}
			if !cityMap.IsOnMap(p) {
				return
			}
			origin = p
			antinodes[p] = struct{}{}
		}
	}
	for _, antennas := range cityMap.Antennas {
		for i, a1 := range antennas {
			for _, a2 := range antennas[i+1:] {
				diffX := a2.X - a1.X
				diffY := a2.Y - a1.Y
				fillAntiNodes(a2, diffX, diffY)
				fillAntiNodes(a1, -1*diffX, -1*diffY)
			}
		}
	}
	fmt.Println(len(antinodes))
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
