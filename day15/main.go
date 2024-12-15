package main

import (
	"flag"
	"fmt"
	"slices"

	"github.com/danieljoos/aoc-2024/core"
)

type Direction int

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

type Position struct{ X, Y int }
type Warehouse struct {
	Walls, Boxes []Position
	Robot        Position
	Moves        []Direction
}

type PositionWide struct{ P1, P2 Position }
type WarehouseWide struct {
	Boxes []PositionWide
	Walls []Position
	Robot Position
	Moves []Direction
}

func readWarehouse() *Warehouse {
	warehouse := &Warehouse{}
	section := 0
	y := 0
	for line := range core.Lines() {
		if line == "" {
			section++
			continue
		}
		switch section {
		case 0:
			for x, c := range line {
				switch c {
				case '#':
					warehouse.Walls = append(warehouse.Walls, Position{x, y})
				case 'O':
					warehouse.Boxes = append(warehouse.Boxes, Position{x, y})
				case '@':
					warehouse.Robot = Position{x, y}
				}
			}
			y++
		case 1:
			for _, c := range line {
				switch c {
				case '^':
					warehouse.Moves = append(warehouse.Moves, DirectionUp)
				case '>':
					warehouse.Moves = append(warehouse.Moves, DirectionRight)
				case 'v':
					warehouse.Moves = append(warehouse.Moves, DirectionDown)
				case '<':
					warehouse.Moves = append(warehouse.Moves, DirectionLeft)
				}
			}
		}
	}
	return warehouse
}

func widenWarehouse(w *Warehouse) *WarehouseWide {
	ww := &WarehouseWide{
		Walls: make([]Position, 0, len(w.Walls)*2),
		Boxes: make([]PositionWide, 0, len(w.Boxes)),
	}
	for _, wall := range w.Walls {
		ww.Walls = append(ww.Walls,
			Position{X: wall.X * 2, Y: wall.Y},
			Position{X: wall.X*2 + 1, Y: wall.Y},
		)
	}
	for _, box := range w.Boxes {
		p := PositionWide{
			P1: Position{X: box.X * 2, Y: box.Y},
			P2: Position{X: box.X*2 + 1, Y: box.Y},
		}
		ww.Boxes = append(ww.Boxes, p)
	}
	ww.Robot.X = w.Robot.X * 2
	ww.Robot.Y = w.Robot.Y
	ww.Moves = w.Moves
	return ww
}

func (w *Warehouse) MoveRobot() {
outer:
	for _, m := range w.Moves {
		pNext := w.Robot.Moved(m)
		if isWall := slices.Contains(w.Walls, pNext); isWall {
			continue
		}

		if boxIdx := slices.Index(w.Boxes, pNext); boxIdx != -1 {
			box := &w.Boxes[boxIdx]
			boxes := []*Position{box}
			for {
				pBox := box.Moved(m)
				if slices.Contains(w.Walls, pBox) {
					continue outer
				}
				if boxIdx := slices.Index(w.Boxes, pBox); boxIdx != -1 {
					box = &w.Boxes[boxIdx]
					boxes = append(boxes, box)
				} else {
					break
				}
			}
			for _, b := range boxes {
				pBox := b.Moved(m)
				b.X = pBox.X
				b.Y = pBox.Y
			}
		}

		w.Robot.X = pNext.X
		w.Robot.Y = pNext.Y
	}
}

func (w *WarehouseWide) MoveRobot() {
	hits := func(p Position, pw PositionWide) bool {
		return pw.P1 == p || pw.P2 == p
	}
outer:
	for _, m := range w.Moves {
		pNext := w.Robot.Moved(m)
		if isWall := slices.Contains(w.Walls, pNext); isWall {
			continue
		}

		if boxIdx := slices.IndexFunc(w.Boxes, func(box PositionWide) bool {
			return hits(pNext, box)
		}); boxIdx != -1 {
			box := &w.Boxes[boxIdx]
			boxes := map[*PositionWide]struct{}{box: {}}
			next := []PositionWide{*box}
			for len(next) > 0 {
				nextNext := []PositionWide{}
				for _, p := range next {
					pNextBox := p.Moved(m)
					if slices.Contains(w.Walls, pNextBox.P1) || slices.Contains(w.Walls, pNextBox.P2) {
						continue outer
					}
					for i, b := range core.Filter2(slices.All(w.Boxes), func(i int, b PositionWide) bool {
						return (hits(pNextBox.P1, b) || hits(pNextBox.P2, b)) && !(hits(p.P1, b) || hits(p.P2, b))
					}) {
						boxes[&w.Boxes[i]] = struct{}{}
						nextNext = append(nextNext, b)
					}
				}
				next = nextNext
			}
			for b := range boxes {
				pBox := b.Moved(m)
				b.P1 = pBox.P1
				b.P2 = pBox.P2
			}
		}

		w.Robot.X = pNext.X
		w.Robot.Y = pNext.Y
	}
}

func (p Position) Moved(dir Direction) Position {
	switch dir {
	case DirectionUp:
		p.Y--
	case DirectionRight:
		p.X++
	case DirectionDown:
		p.Y++
	case DirectionLeft:
		p.X--
	}
	return p
}

func (p Position) GPS() int {
	return 100*p.Y + p.X
}

func (p PositionWide) Moved(dir Direction) PositionWide {
	p.P1 = p.P1.Moved(dir)
	p.P2 = p.P2.Moved(dir)
	return p
}

func (p PositionWide) GPS() int {
	return p.P1.GPS()
}

func part1() {
	warehouse := readWarehouse()
	warehouse.MoveRobot()
	sum := 0
	for _, b := range warehouse.Boxes {
		sum += b.GPS()
	}
	fmt.Println(sum)
}

func part2() {
	warehouse := widenWarehouse(readWarehouse())
	warehouse.MoveRobot()
	sum := 0
	for _, b := range warehouse.Boxes {
		sum += b.GPS()
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
