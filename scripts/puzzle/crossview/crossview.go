// Package crossview provides crossview puzzles
package crossview

import (
	"errors"
)

var (
	errInvalidSize  = errors.New("size must be a positive even number")
	errSizeTooLarge = errors.New("size can at max be 32")
)

type Crossview struct {
	Size  int
	Field []uint8
}

func New(size int) (*Crossview, error) {
	if err := validateSize(size); err != nil {
		return nil, err
	}

	c := &Crossview{
		Size: size,
	}

	if err := c.generate(); err != nil {
		return nil, err
	}

	return c, nil
}

// Puzzle generates a new puzzle mask for the field
// When no filled percentage is provided, the difficulty default is used
func (c *Crossview) Puzzle(diff Difficulty, filledPercentage ...float64) ([]uint8, error) {
	return newGrid(c.Size, c.Field).generateMask(diff, filledPercentage...)
}

// PuzzleString renders the solved field using the given mask.
func (c *Crossview) PuzzleString(mask []uint8) (string, error) {
	return newGrid(c.Size, c.Field).string(mask)
}

func (c *Crossview) String() string {
	grid := c.grid()
	rendered, err := grid.string()
	if err != nil {
		return ""
	}

	return rendered
}
