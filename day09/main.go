package main

import (
	"flag"
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/danieljoos/aoc-2024/core"
)

type DiskMap []int

type File struct {
	id         int
	length     int
	prev, next *File
}

func ReadDiskMap() DiskMap {
	for line := range core.Lines() {
		result := make([]int, 0, len(line))
		for _, v := range line {
			intval := int(v - '0')
			result = append(result, intval)
		}
		return result
	}
	return nil
}

func (dm DiskMap) String() string {
	var buf strings.Builder
	for _, v := range dm {
		if v < 0 {
			buf.WriteRune('.')
		} else {
			buf.WriteRune(rune(v + '0'))
		}
	}
	return buf.String()
}

func (dm DiskMap) Decompress() DiskMap {
	var decompressed DiskMap
	for i, v := range dm {
		var id int
		if i%2 > 0 {
			id = -1
		} else {
			id = int(i / 2)
		}
		for range v {
			decompressed = append(decompressed, id)
		}
	}
	return decompressed
}

func (dm DiskMap) Checksum() int64 {
	checksum := int64(0)
	for i, v := range dm {
		if v < 0 {
			continue
		}
		checksum += int64(i * v)
	}
	return checksum
}

func (dm DiskMap) Files() *File {
	var root, curr *File
	for i, v := range dm {
		f := &File{length: v, prev: curr}
		if i%2 > 0 {
			f.id = -1
		} else {
			f.id = i / 2
		}
		if root == nil {
			root = f
		} else {
			curr.next = f
		}
		curr = f
	}
	return root
}

func (f *File) Iter() iter.Seq[*File] {
	return func(yield func(*File) bool) {
		if f == nil {
			return
		}
		for {
			if f == nil {
				return
			}
			if !yield(f) {
				return
			}
			f = f.next
		}
	}
}

func (f *File) DecompressedDiskMap() DiskMap {
	result := DiskMap{}
	for f := range f.Iter() {
		for range f.length {
			result = append(result, f.id)
		}
	}
	return result
}

func part1() {
	dm := ReadDiskMap().Decompress()
	for i := range slices.Backward(dm) {
		free := slices.Index(dm, -1)
		if free >= i {
			break
		}
		dm[i], dm[free] = dm[free], dm[i]
	}
	fmt.Println(dm.Checksum())
}

func part2() {
	files := ReadDiskMap().Files()
	for _, f1 := range slices.Backward(slices.Collect(files.Iter())) {
		if f1.id < 0 {
			continue
		}
		for f2 := range files.Iter() {
			if f1 == f2 {
				break
			}
			if f2.id < 0 && f2.length >= f1.length {
				f2.length -= f1.length

				// insert new free node
				free := &File{id: -1, length: f1.length, prev: f1.prev, next: f1}
				if free.prev != nil {
					free.prev.next = free
				}
				f1.prev = free

				// disconnect node
				if f1.prev != nil {
					f1.prev.next = f1.next
				}
				if f1.next != nil {
					f1.next.prev = f1.prev
				}

				// insert before free node
				f1.prev = f2.prev
				f1.next = f2
				f2.prev = f1
				if f1.prev != nil {
					f1.prev.next = f1
				}
				break
			}
		}
	}
	fmt.Println(files.DecompressedDiskMap().Checksum())
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
