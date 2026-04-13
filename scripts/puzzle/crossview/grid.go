package crossview

import (
	"fmt"
	"slices"
	"strings"
)

const maxSize = 32

const (
	blocked uint8 = iota + maxSize + 1
	unknown
	filled
)

type axis uint8

const (
	rowAxis axis = iota + 1
	colAxis
)

type grid struct {
	size  int
	cells []uint8
}

func (c Crossview) grid() grid {
	return grid{
		size:  c.size,
		cells: c.field,
	}
}

func emptyGrid(size int) grid {
	cells := make([]uint8, 0, size*size)

	maxCell := uint8(size*2 - 2)

	for range size * size {
		cells = append(cells, maxCell)
	}

	return grid{
		size:  size,
		cells: cells,
	}
}

func newGrid(size int, board []uint8) grid {
	return grid{
		size:  size,
		cells: board,
	}
}

func (g grid) at(rowIdx, colIdx int) uint8 {
	return g.cells[rowIdx*g.size+colIdx]
}

func (g grid) set(rowIdx, colIdx int, value uint8) {
	g.cells[rowIdx*g.size+colIdx] = value
}

func (g grid) sees(rowIdx, colIdx int) int {
	if g.at(rowIdx, colIdx) == blocked {
		return 0
	}

	row := g.row(rowIdx)
	col := g.col(colIdx)

	return row.sees(colIdx) + col.sees(rowIdx)
}

func (g grid) row(rowIdx int) line {
	return line{
		grid:  g,
		axis:  rowAxis,
		index: rowIdx,
	}
}

func (g grid) col(colIdx int) line {
	return line{
		grid:  g,
		axis:  colAxis,
		index: colIdx,
	}
}

func (g grid) solved() bool {
	return !slices.Contains(g.cells, unknown)
}

// addBoundary adds a boundary to a cell
// It returns true if it added a boundary somewhere
func (g grid) addBoundary(rowIdx, colIdx int) bool {
	return g.row(rowIdx).addBoundary(colIdx) || g.col(colIdx).addBoundary(rowIdx)
}

// fill fills in any non blocked open cell
// It returns the amount of filled in cells
func (g grid) fill(rowIdx, colIdx int) int {
	return g.row(rowIdx).fill(colIdx) + g.col(colIdx).fill(rowIdx)
}

// open returns the amount of non blocked open cells
func (g grid) open(rowIdx, colIdx int) int {
	return g.row(rowIdx).open(colIdx) + g.col(colIdx).open(rowIdx)
}

// valid check if a grid is valid
// It simply checks if any number sees at maximum that amount
// And that it still has the possibility to reach it requires visible cells
func (g grid) valid() bool {
	for r := range g.size {
		for c := range g.size {
			cell := g.at(r, c)
			if cell >= blocked {
				continue
			}

			sees := g.sees(r, c)
			if uint8(sees) > cell {
				return false
			}

			open := g.open(r, c)
			if uint8(sees+open) < cell {
				return false
			}
		}
	}

	return true
}

func (g grid) string() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "\n Size: %d * %d", g.size, g.size)

	divider := "\n "
	for range g.size {
		divider += "-----"
	}
	divider += "-\n"

	for r := range g.size {
		builder.WriteString(divider)
		for c := range g.size {
			cell := g.at(r, c)
			value := fmt.Sprintf("%02d", cell)
			if cell == blocked {
				value = "XX"
			}
			builder.WriteString(" | " + value)
		}
		builder.WriteString(" |")
	}
	builder.WriteString(divider)

	return builder.String()
}

func validateSize(size int) error {
	if size <= 0 {
		return errInvalidSize
	}

	if size > maxSize {
		return errSizeTooLarge
	}

	return nil
}
