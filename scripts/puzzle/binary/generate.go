package binary

import (
	"errors"
	"math/rand"
)

type lineState struct {
	counts cellCounts
	bits   uint32 // Limits size to 32 * 32
}

func (s *lineState) add(value uint8, bitIdx int) {
	s.counts.add(value)
	if value == 1 {
		s.bits |= uint32(1) << bitIdx
	}
}

func (s *lineState) remove(value uint8, bitIdx int) {
	switch value {
	case 0:
		s.counts.zeros--
	case 1:
		s.counts.ones--
		s.bits &^= uint32(1) << bitIdx
	}
}

// generate Generates a new puzzle
func (b *Binary) generate() error {
	// Reset the field
	b.Field = make([]uint8, b.Size*b.Size)

	// Get all possible rows
	rows := generateRows(b.Size)
	if len(rows) == 0 {
		return errors.New("no possible rows")
	}

	// Shuffle rows
	rand.Shuffle(len(rows), func(i, j int) {
		rows[i], rows[j] = rows[j], rows[i]
	})

	usedRows := make([]bool, len(rows))
	states := make([]lineState, b.Size)
	g := newGrid(b.Size, b.Field)

	// Try to fill in field with rows
	ok := b.fillRows(g, rows, usedRows, states, 0)
	if !ok {
		return errors.New("unable to fill field")
	}

	return nil
}

func (b *Binary) fillRows(g grid, candidates [][]uint8, usedRows []bool, states []lineState, rowIdx int) bool {
	if rowIdx == g.size {
		return true
	}

	// Keep trying every unused candidate until one leads to a solution
	for idx, candidate := range candidates {
		// Already used
		if usedRows[idx] {
			continue
		}

		if !canPlaceRow(g, states, rowIdx, candidate) {
			continue
		}

		g.setRow(rowIdx, candidate)
		applyRow(states, candidate, rowIdx) // Keeps track of columns
		usedRows[idx] = true

		ok := true
		// If it is the last row check for column uniqueness
		// Row uniques is guaranteed by how the rows are generated
		if rowIdx == g.size-1 && !columnsUnique(states) {
			ok = false
		}

		// Go deeper
		if ok && b.fillRows(g, candidates, usedRows, states, rowIdx+1) {
			return true
		}

		// Unable to complete the field
		// Try another row
		usedRows[idx] = false
		undoRow(states, candidate, rowIdx)
	}

	return false
}

// canPlaceRow checks if placing the row won't cause any immediate conflicts
func canPlaceRow(g grid, states []lineState, rowIdx int, candidate []uint8) bool {
	size := g.size
	half := size / 2
	remainingRows := size - rowIdx - 1

	for col, state := range states {
		bit := candidate[col]
		counts := state.counts
		counts.add(bit)

		if !counts.canBalanceWithRemaining(half, remainingRows) {
			return false
		}

		// No triplets
		if rowIdx >= 2 {
			a := g.at(rowIdx-2, col)
			prev := g.at(rowIdx-1, col)
			if bit == a && bit == prev {
				return false
			}
		}
	}

	return true
}

// applyRow updates the affected columns state if a row is added
func applyRow(states []lineState, candidate []uint8, rowIdx int) {
	for col, bit := range candidate {
		states[col].add(bit, rowIdx)
	}
}

// undoRow updates the affected columns state if a row is removed
func undoRow(states []lineState, candidate []uint8, rowIdx int) {
	for col, bit := range candidate {
		states[col].remove(bit, rowIdx)
	}
}

// columnsUnique checks if every column is unique
func columnsUnique(states []lineState) bool {
	// The bit state to determine if a column is equal is maintained by applyRow and undoRow
	seen := make(map[uint32]bool, len(states))
	for _, s := range states {
		if _, ok := seen[s.bits]; ok {
			return false
		}
		seen[s.bits] = true
	}

	return true
}

// generateRows generates all possible rows for a given size
// It only generates valid and unique rows
func generateRows(size int) [][]uint8 {
	half := size / 2
	out := make([][]uint8, 0)

	buf := make([]uint8, size)

	var gen func(col, zeros, ones int)
	gen = func(col, zeros, ones int) {
		// Max zeros and ones == half of the size
		if zeros > half || ones > half {
			return
		}

		// Is reaching half even still possible?
		remaining := size - col
		if zeros+remaining < half || ones+remaining < half {
			return
		}

		// No triplets
		if col >= 3 {
			if buf[col-3] == buf[col-2] && buf[col-2] == buf[col-1] {
				return
			}
		}

		// Full row
		if col == size {
			if zeros == half && ones == half {
				cpy := make([]uint8, size)
				copy(cpy, buf)
				out = append(out, cpy)
			}
			return
		}

		buf[col] = 0
		gen(col+1, zeros+1, ones)

		buf[col] = 1
		gen(col+1, zeros, ones+1)
	}

	gen(0, 0, 0)
	return out
}
