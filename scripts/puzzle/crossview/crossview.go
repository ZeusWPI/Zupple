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
	size  int
	field []uint8
}

func New(size int) (*Crossview, error) {
	if err := validateSize(size); err != nil {
		return nil, err
	}

	c := &Crossview{
		size: size,
	}

	if err := c.generate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Crossview) String() string {
	return c.grid().string()
}
