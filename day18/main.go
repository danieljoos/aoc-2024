package main

import (
	"flag"
	"fmt"
	"iter"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/danieljoos/aoc-2024/core"
)

type Position struct{ X, Y int }
type PositionSet map[Position]struct{}

type Vertex struct {
	Position
	Distance int
	Visited  bool
}

type MemorySpace []Vertex

func (p Position) NeighboringPositions() []Position {
	return []Position{
		{p.X - 1, p.Y},
		{p.X + 1, p.Y},
		{p.X, p.Y - 1},
		{p.X, p.Y + 1},
	}
}

func (ps PositionSet) Has(p Position) bool {
	_, ok := ps[p]
	return ok
}

func newMemorySpace(width, height int) MemorySpace {
	ms := make(MemorySpace, width*height)
	for y := range height {
		for x := range width {
			ms[y*width+x] = Vertex{Position: Position{x, y}, Distance: math.MaxInt}
		}
	}
	return ms
}

func (ms MemorySpace) Reset() {
	for i := range ms {
		ms[i].Distance = math.MaxInt
		ms[i].Visited = false
	}
}

func corruptedPositions() iter.Seq2[Position, struct{}] {
	return func(yield func(Position, struct{}) bool) {
		next, stop := iter.Pull(core.Lines())
		defer stop()
		for {
			val, ok := next()
			if !ok {
				return
			}
			fields := strings.SplitN(val, ",", 2)
			x, _ := strconv.Atoi(fields[0])
			y, _ := strconv.Atoi(fields[1])
			if !yield(Position{x, y}, struct{}{}) {
				return
			}
		}
	}
}

func shortestPath(corrupted PositionSet, ms MemorySpace, target Position) int {
	ms.Reset()
	ms[0].Distance = 0
	for {
		// Find vertex with min distance:
		var u *Vertex
		for i := range ms {
			v := &ms[i]
			if !v.Visited && (u == nil || v.Distance < u.Distance) {
				u = v
			}
		}
		if u == nil {
			break
		}
		u.Visited = true

		if u.Position == target {
			return u.Distance
		}

		if u.Distance < 0 || u.Distance == math.MaxInt {
			// no path, we can stop here
			return -1
		}

		// Find neighbors that are still in q:
		neighboringPos := u.NeighboringPositions()
		count := 0
		for i := range ms {
			n := &ms[i]
			if count >= len(neighboringPos) ||
				!slices.Contains(neighboringPos, n.Position) ||
				corrupted.Has(n.Position) {
				continue
			}
			alt := u.Distance + 1
			if alt < n.Distance {
				n.Distance = alt
			}
			count++
		}
	}
	return -1
}

func part1() {
	cp := PositionSet{}
	maps.Insert(cp, core.Take2(corruptedPositions(), 1024))
	ms := newMemorySpace(71, 71)
	fmt.Println(shortestPath(cp, ms, Position{70, 70}))
}

func part2() {
	cp := PositionSet{}
	cpIter := corruptedPositions()
	maps.Insert(cp, core.Take2(cpIter, 1024))
	ms := newMemorySpace(71, 71)
	target := Position{70, 70}
	for k, v := range cpIter {
		fmt.Println(k, v)
		cp[k] = v
		sp := shortestPath(cp, ms, target)
		if sp < 0 || sp == math.MaxInt {
			fmt.Printf("no path: %d,%d\n", k.X, k.Y)
			return
		}
	}
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
