package main

import (
	"flag"
	"fmt"
	"maps"
	"slices"

	"github.com/danieljoos/aoc-2024/core"
)

type Position struct{ X, Y int }

type none struct{}

type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionBottom
	DirectionLeft
	Direction_MAX
)

type Guard struct {
	Position
	Direction Direction
	IsOff     bool
}

type Field struct {
	Width, Height int
	Guard         Guard
	Obstacles     map[Position]none
}

func readField() *Field {
	field := &Field{Obstacles: map[Position]none{}}
	y := 0
	for line := range core.Lines() {
		field.Width = len(line)
		field.Height++
		for x, r := range line {
			switch r {
			case '^':
				field.Guard = Guard{Position: Position{x, y}, Direction: DirectionUp}
			case '#':
				field.Obstacles[Position{x, y}] = none{}
			}
		}
		y++
	}
	return field
}

func (f *Field) MoveGuard() {
	guard := &f.Guard
	var nextPos Position
	switch guard.Direction {
	case DirectionUp:
		nextPos = Position{guard.X, guard.Y - 1}
	case DirectionRight:
		nextPos = Position{guard.X + 1, guard.Y}
	case DirectionBottom:
		nextPos = Position{guard.X, guard.Y + 1}
	case DirectionLeft:
		nextPos = Position{guard.X - 1, guard.Y}
	}
	if _, hitsObstacle := f.Obstacles[nextPos]; hitsObstacle {
		guard.Direction = (guard.Direction + 1) % Direction_MAX
		return
	}
	guard.Position = nextPos
	guard.IsOff = guard.X < 0 || guard.Y < 0 || guard.X >= f.Width || guard.Y >= f.Height
}

func (f *Field) Clone() *Field {
	cloned := *f
	cloned.Obstacles = maps.Clone(f.Obstacles)
	return &cloned
}

func patrolPositions(field *Field) []Position {
	field = field.Clone()
	visited := map[Position]none{}
	for {
		visited[field.Guard.Position] = none{}
		field.MoveGuard()
		if field.Guard.IsOff {
			break
		}
	}
	return slices.Collect(maps.Keys(visited))
}

func part1() {
	field := readField()
	visited := patrolPositions(field)
	fmt.Println(len(visited))
}

func part2() {
	field := readField()
	loopCount := 0
	for _, pos := range patrolPositions(field) {
		if field.Guard.Position == pos {
			continue
		}

		cloned := field.Clone()
		cloned.Obstacles[pos] = none{}
		visited := map[Guard]none{}
		for {
			if _, loop := visited[cloned.Guard]; loop {
				loopCount++
				break
			}
			if cloned.Guard.IsOff {
				break
			}
			visited[cloned.Guard] = none{}
			cloned.MoveGuard()
		}
	}
	fmt.Println(loopCount)
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
