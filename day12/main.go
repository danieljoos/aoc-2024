package main

import (
	"flag"
	"fmt"
	"iter"
	"slices"

	"github.com/danieljoos/aoc-2024/core"
)

type Garden []string
type Position struct{ X, Y int }
type Region map[Position]struct{}
type FenceSide int

const (
	FenceSideRight FenceSide = iota
	FenceSideBottom
	FenceSideLeft
	FenceSideTop
	FenceSideMAX
)

func fillRegion(garden Garden, pos Position) Region {
	width := len(garden[0])
	height := len(garden)
	region := Region{}
	plant := garden[pos.Y][pos.X]
	q := []Position{pos}
	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		if _, has := region[n]; has {
			continue
		}
		if n.X < 0 || n.Y < 0 || n.X >= width || n.Y >= height ||
			garden[n.Y][n.X] != plant {
			continue
		}
		region[n] = struct{}{}
		q = append(q,
			Position{n.X, n.Y + 1},
			Position{n.X, n.Y - 1},
			Position{n.X + 1, n.Y},
			Position{n.X - 1, n.Y},
		)
	}
	return region
}

func (g Garden) Positions() iter.Seq[Position] {
	return func(yield func(Position) bool) {
		for y, row := range g {
			for x := range row {
				if !yield(Position{x, y}) {
					return
				}
			}
		}
	}
}

func (g Garden) Regions() []Region {
	regions := []Region{}
	for pos := range g.Positions() {
		if slices.IndexFunc(regions, func(r Region) bool {
			_, has := r[pos]
			return has
		}) == -1 {
			regions = append(regions, fillRegion(g, pos))
		}
	}
	return regions
}

func (r Region) Has(p Position) bool {
	_, has := r[p]
	return has
}

func (r Region) Perimeter() int {
	permimeter := 0
	for pos := range r {
		for _, p := range []Position{
			{pos.X, pos.Y + 1},
			{pos.X, pos.Y - 1},
			{pos.X + 1, pos.Y},
			{pos.X - 1, pos.Y},
		} {
			if !r.Has(p) {
				permimeter++
			}
		}
	}
	return permimeter
}

func (r Region) OuterFenceIter() iter.Seq2[Position, FenceSide] {
	return func(yield func(Position, FenceSide) bool) {
		var start Position
		side := FenceSideRight
		for pos := range r {
			if !r.Has(Position{pos.X + 1, pos.Y}) && pos.X > start.X {
				start = pos
			}
		}
		p := start
		for {
			var c1, c2 Position
			switch side {
			case FenceSideRight:
				c1 = Position{p.X + 1, p.Y + 1}
				c2 = Position{p.X, p.Y + 1}
			case FenceSideBottom:
				c1 = Position{p.X - 1, p.Y + 1}
				c2 = Position{p.X - 1, p.Y}
			case FenceSideLeft:
				c1 = Position{p.X - 1, p.Y - 1}
				c2 = Position{p.X, p.Y - 1}
			case FenceSideTop:
				c1 = Position{p.X + 1, p.Y - 1}
				c2 = Position{p.X + 1, p.Y}
			}
			if r.Has(c1) {
				side = (FenceSideMAX + side - 1) % FenceSideMAX
				p = c1
			} else if !r.Has(c2) {
				side = (FenceSideMAX + side + 1) % FenceSideMAX
			} else {
				p = c2
			}
			if !yield(p, side) {
				return
			}
			if p == start && side == FenceSideRight {
				break
			}
		}
	}
}

func (r Region) OuterSides() int {
	sides := 0
	last := FenceSideRight
	for _, side := range r.OuterFenceIter() {
		if side != last {
			sides++
		}
		last = side
	}
	return sides
}

func (r Region) Area() int {
	return len(r)
}

func (r Region) Encloses(inner Region) bool {
	for pos, side := range inner.OuterFenceIter() {
		var p Position
		switch side {
		case FenceSideRight:
			p = Position{pos.X + 1, pos.Y}
		case FenceSideBottom:
			p = Position{pos.X, pos.Y + 1}
		case FenceSideLeft:
			p = Position{pos.X - 1, pos.Y}
		case FenceSideTop:
			p = Position{pos.X, pos.Y - 1}
		}
		if !r.Has(p) {
			return false
		}
	}
	return true
}

func part1() {
	garden := Garden(slices.Collect(core.Lines()))
	sum := 0
	for _, r := range garden.Regions() {
		sum += r.Area() * r.Perimeter()
	}
	fmt.Println(sum)
}

func part2() {
	garden := Garden(slices.Collect(core.Lines()))
	regions := garden.Regions()
	outerSides := make([]int, len(regions))
	for i, r := range regions {
		outerSides[i] = r.OuterSides()
	}
	sum := 0
	for i, r1 := range regions {
		sides := outerSides[i]
		for j, r2 := range regions {
			if i == j {
				continue
			}
			if r1.Encloses(r2) {
				sides += outerSides[j]
			}
		}
		sum += sides * r1.Area()
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
