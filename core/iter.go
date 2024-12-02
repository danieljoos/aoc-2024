package core

import (
	"iter"
	"strconv"
	"strings"
)

func Filter[T any](vals iter.Seq[T], pred func(v T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(vals)
		defer stop()
		for {
			val, ok := next()
			if !ok {
				return
			}
			if pred(val) {
				yield(val)
			}
		}
	}
}

func IntVals(lines iter.Seq[string]) iter.Seq[[]int] {
	return func(yield func([]int) bool) {
		next, stop := iter.Pull(lines)
		defer stop()
		for {
			line, ok := next()
			if !ok {
				return
			}
			fields := strings.Fields(line)
			ints := make([]int, 0, len(fields))
			for _, field := range fields {
				i, err := strconv.Atoi(field)
				if err != nil {
					continue
				}
				ints = append(ints, i)
			}
			yield(ints)
		}
	}
}
