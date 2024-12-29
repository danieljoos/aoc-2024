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
				if !yield(val) {
					return
				}
			}
		}
	}
}

func Filter2[K any, T any](vals iter.Seq2[K, T], pred func(i K, v T) bool) iter.Seq2[K, T] {
	return func(yield func(K, T) bool) {
		next, stop := iter.Pull2(vals)
		defer stop()
		for {
			k, val, ok := next()
			if !ok {
				return
			}
			if pred(k, val) {
				if !yield(k, val) {
					return
				}
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
			if !yield(ints) {
				return
			}
		}
	}
}

func StrVals(items iter.Seq[int]) iter.Seq[string] {
	return func(yield func(string) bool) {
		next, stop := iter.Pull(items)
		defer stop()
		for {
			item, ok := next()
			if !ok {
				return
			}
			if !yield(strconv.Itoa(item)) {
				return
			}
		}
	}
}

func First[T any](items iter.Seq[T]) T {
	for item := range items {
		return item
	}
	var result T
	return result
}

func Take[T any](items iter.Seq[T], num int) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(items)
		defer stop()
		for i := 0; i < num; i++ {
			item, ok := next()
			if !ok || !yield(item) {
				return
			}
		}
	}
}

func Take2[T, K any](items iter.Seq2[T, K], num int) iter.Seq2[T, K] {
	return func(yield func(T, K) bool) {
		next, stop := iter.Pull2(items)
		defer stop()
		for i := 0; i < num; i++ {
			item1, item2, ok := next()
			if !ok || !yield(item1, item2) {
				return
			}
		}
	}
}
