package core

import (
	"bufio"
	"iter"
	"os"
)

func Lines() iter.Seq[string] {
	return func(yield func(string) bool) {
		file, err := os.OpenFile("input.txt", os.O_RDONLY, 0)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			yield(scanner.Text())
		}
	}
}
