package main

import (
	"flag"
	"fmt"
	"iter"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/danieljoos/aoc-2024/core"
)

const (
	Width  = 101
	Height = 103
)

type Robot struct {
	Position [2]int
	Velocity [2]int
}

func readRobots() iter.Seq[*Robot] {
	pattern := regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)
	return func(yield func(*Robot) bool) {
		next, stop := iter.Pull(core.Lines())
		defer stop()
		for {
			line, ok := next()
			if !ok {
				return
			}
			match := pattern.FindStringSubmatch(line)
			if match == nil {
				continue
			}
			px, _ := strconv.Atoi(match[1])
			py, _ := strconv.Atoi(match[2])
			vx, _ := strconv.Atoi(match[3])
			vy, _ := strconv.Atoi(match[4])
			rob := &Robot{Position: [2]int{px, py}, Velocity: [2]int{vx, vy}}
			if !yield(rob) {
				return
			}
		}
	}
}

func drawRobotsImageData(robots []*Robot) string {
	hasRobot := func(x, y int) bool {
		return slices.ContainsFunc(robots, func(r *Robot) bool {
			return r.Position[0] == x && r.Position[1] == y
		})
	}
	data := make([]string, 0, Width*Height)
	for y := range Height {
		for x := range Width {
			if hasRobot(x, y) {
				data = append(data, "0", "0", "0", "255")
			} else {
				data = append(data, "255", "255", "255", "255")
			}
		}
	}
	return "[" + strings.Join(data, ",") + "]"
}

func moveRobots(robots []*Robot, n int) {
	for _, r := range robots {
		r.Position[0] = (Width + r.Position[0] + (r.Velocity[0]*n)%Width) % Width
		r.Position[1] = (Height + r.Position[1] + (r.Velocity[1]*n)%Height) % Height
	}
}

func safetyFactor(robots []*Robot) int {
	w2 := Width / 2
	h2 := Height / 2
	quarters := [4]int{}
	for _, robot := range robots {
		switch {
		case robot.Position[0] < w2 && robot.Position[1] < h2:
			quarters[0]++
		case robot.Position[0] > w2 && robot.Position[1] < h2:
			quarters[1]++
		case robot.Position[0] < w2 && robot.Position[1] > h2:
			quarters[2]++
		case robot.Position[0] > w2 && robot.Position[1] > h2:
			quarters[3]++
		}
	}
	mul := 1
	for _, v := range quarters {
		mul *= v
	}
	return mul
}

func part1() {
	robots := slices.Collect(readRobots())
	moveRobots(robots, 100)
	fmt.Println(safetyFactor(robots))
}

func part2() {
	fmt.Println("Serving on port 8080")
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		n, _ := strconv.Atoi(q.Get("n"))

		robots := slices.Collect(readRobots())
		moveRobots(robots, n)
		for n < 1000000 {
			// Just trial-and-error'ed the safety factor limit...
			// Then click "next" until you find an image with a tree (see image.png)
			if sf := safetyFactor(robots); sf < 100000000 {
				break
			}
			moveRobots(robots, 1)
			n++
		}

		imgData := drawRobotsImageData(robots)
		fmt.Fprintf(w, `
		<html>
		<canvas id="robots" style="zoom:300%%"></canvas>
		<p>n=%d <a href="?n=%d">Next</a></p>
		<script>
			const canvas = document.getElementById('robots');
			const ctx = canvas.getContext('2d');
			const arr = new Uint8ClampedArray(%s);
			const imgData = new ImageData(arr, %d, %d);
			ctx.putImageData(imgData, 0, 0);
		</script
		</html>`, n, n+1, imgData, Width, Height)
	}))
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
