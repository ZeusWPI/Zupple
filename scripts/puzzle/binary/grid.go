package binary

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const (
	maxSize       = 32
	unknown uint8 = 2
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

func newGrid(size int, cells []uint8) grid {
	return grid{
		size:  size,
		cells: cells,
	}
}

func (g grid) at(rowIdx, colIdx int) uint8 {
	return g.cells[rowIdx*g.size+colIdx]
}

func (g grid) set(rowIdx, colIdx int, value uint8) {
	g.cells[rowIdx*g.size+colIdx] = value
}

func (g grid) setRow(rowIdx int, values []uint8) {
	start := rowIdx * g.size
	copy(g.cells[start:start+g.size], values)
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

func (g grid) compress() []byte {
	byteCount := (len(g.cells) + 7) / 8
	out := make([]byte, 0, byteCount)

	var buf byte
	bufIdx := 0

	for _, value := range g.cells {
		if value == 1 {
			buf |= 1 << bufIdx
		}

		bufIdx++
		if bufIdx == 8 {
			out = append(out, buf)
			buf = 0
			bufIdx = 0
		}
	}

	if bufIdx > 0 {
		out = append(out, buf)
	}

	return out
}

func (g grid) loadCompressed(data []byte) error {
	cellCount := g.size * g.size
	byteCount := (cellCount + 7) / 8

	if len(data) != byteCount {
		return errInvalidCompressedLength
	}

	count := 0
	for _, value := range data {
		for bitIdx := range 8 {
			if count >= cellCount {
				return nil
			}

			g.cells[count] = uint8((value >> bitIdx) & 1)
			count++
		}
	}

	return nil
}

func (g grid) string(mask ...[]uint8) (string, error) {
	var activeMask []uint8
	if len(mask) > 0 {
		activeMask = mask[0]
		if err := g.validateMask(activeMask); err != nil {
			return "", err
		}
	}

	var builder strings.Builder

	fmt.Fprintf(&builder, "\n Size: %d * %d", g.size, g.size)

	divider := "\n "
	for range g.size {
		divider += "----"
	}
	divider += "-\n"

	for r := range g.size {
		builder.WriteString(divider)
		for c := range g.size {
			cellIdx := r*g.size + c
			value := strconv.Itoa(int(g.at(r, c)))
			if activeMask != nil && activeMask[cellIdx] == 0 {
				value = " "
			}
			builder.WriteString(" | " + value)
		}
		builder.WriteString(" |")
	}
	builder.WriteString(divider)

	return builder.String(), nil
}

type line struct {
	grid  grid
	axis  axis
	index int
}

func (l line) len() int {
	return l.grid.size
}

func (l line) inBounds(pos int) bool {
	return pos >= 0 && pos < l.len()
}

func (l line) at(pos int) uint8 {
	if l.axis == rowAxis {
		return l.grid.at(l.index, pos)
	}

	return l.grid.at(pos, l.index)
}

func (l line) set(pos int, value uint8) {
	if l.axis == rowAxis {
		l.grid.set(l.index, pos, value)
		return
	}

	l.grid.set(pos, l.index, value)
}

func (l line) values() []uint8 {
	values := make([]uint8, l.len())
	for pos := range l.len() {
		values[pos] = l.at(pos)
	}
	return values
}

func (l line) counts() cellCounts {
	var counts cellCounts

	for pos := range l.len() {
		counts.add(l.at(pos))
	}

	return counts
}

func (l line) unknownPositions() []int {
	positions := make([]int, 0)

	for pos := range l.len() {
		if l.at(pos) == unknown {
			positions = append(positions, pos)
		}
	}

	return positions
}

type cellCounts struct {
	zeros    int
	ones     int
	unknowns int
}

func (c *cellCounts) add(value uint8) {
	switch value {
	case 0:
		c.zeros++
	case 1:
		c.ones++
	default:
		c.unknowns++
	}
}

func (c cellCounts) canBalanceWithRemaining(half, remaining int) bool {
	return c.zeros <= half &&
		c.ones <= half &&
		c.zeros+remaining >= half &&
		c.ones+remaining >= half
}

func (c cellCounts) canStillBalance(half int) bool {
	return c.canBalanceWithRemaining(half, c.unknowns)
}

func (c cellCounts) completionValue(half int) (uint8, bool) {
	if c.unknowns == 0 {
		return 0, false
	}

	if c.zeros == half {
		return 1, true
	}

	if c.ones == half {
		return 0, true
	}

	return 0, false
}

func validateSize(size int) error {
	if size <= 0 || size%2 != 0 {
		return errInvalidSize
	}

	if size > maxSize {
		return errSizeTooLarge
	}

	return nil
}

func (g grid) validateMask(mask []uint8) error {
	if len(mask) != g.size*g.size {
		return errInvalidMaskLength
	}

	for _, value := range mask {
		if value != 0 && value != 1 {
			return errInvalidMaskValue
		}
	}

	return nil
}

func linesMatchKnown(partial, full []uint8) bool {
	for idx, value := range partial {
		if value == unknown {
			continue
		}

		if full[idx] == unknown || value != full[idx] {
			return false
		}
	}

	return true
}
