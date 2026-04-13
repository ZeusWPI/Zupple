package main

import (
	"fmt"

	"puzzle/binary"
	"puzzle/crossview"
)

const size = 16

func main() {
	// cross()
	bin()
}

func cross() {
	c, err := crossview.New(size)
	if err != nil {
		panic(err)
	}

	fmt.Println(c)
}

func bin() {
	b, err := binary.New(size)
	if err != nil {
		panic(err)
	}

	fmt.Println(b)

	mask, err := b.Puzzle(binary.Expert)
	if err != nil {
		panic(err)
	}

	fmt.Println(b.PuzzleString(mask))
}
