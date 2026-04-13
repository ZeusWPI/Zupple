// Package binary provides binary puzzles
package binary

import (
	"errors"
	"fmt"
	"slices"
)

var (
	errInvalidSize             = errors.New("size must be a positive even number")
	errSizeTooLarge            = errors.New("size can at max be 32")
	errInvalidCompressedLength = errors.New("invalid bytes length for binary size")
)

type Binary struct {
	Size  int
	Field []uint8
}

// New generated a new puzzle
func New(size int) (*Binary, error) {
	if err := validateSize(size); err != nil {
		return nil, err
	}

	binary := &Binary{
		Size: size,
	}

	if err := binary.generate(); err != nil {
		return nil, fmt.Errorf("no grid could be generated %w", err)
	}

	return binary, nil
}

// NewCompressed reconstructs a puzzle based on a compressed version
func NewCompressed(size int, bytes []byte) (*Binary, error) {
	if err := validateSize(size); err != nil {
		return nil, err
	}

	binary := &Binary{
		Size: size,
	}

	if err := binary.decompress(bytes); err != nil {
		return nil, fmt.Errorf("bytes couldn't be decompressed %w", err)
	}

	return binary, nil
}

func (b *Binary) String() string {
	grid := newGrid(b.Size, b.Field)
	rendered, err := grid.string()
	if err != nil {
		return ""
	}

	return rendered
}

// Puzzle generates a new puzzle mask for the field
// When no filled percentage is provided, the difficulty default is used
func (b *Binary) Puzzle(diff Difficulty, filledPercentage ...float64) ([]uint8, error) {
	return newGrid(b.Size, b.Field).generateMask(diff, filledPercentage...)
}

// PuzzleString renders the solved field using the given mask.
func (b *Binary) PuzzleString(mask []uint8) (string, error) {
	return newGrid(b.Size, b.Field).string(mask)
}

func (b *Binary) Equal(b2 Binary) bool {
	if b.Size != b2.Size {
		return false
	}

	if len(b.Field) != len(b2.Field) {
		return false
	}

	return slices.Equal(b.Field, b2.Field)
}
