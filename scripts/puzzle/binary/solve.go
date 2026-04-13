package binary

import (
	"slices"
)

type Difficulty int

const (
	Easy Difficulty = iota + 1
	Medium
	Hard
	Expert
)

type solveMode struct {
	allowTriplets     bool // Avoid triplets
	allowBalance      bool // Same amount of 0 and 1 in a row / column
	allowUniqueness   bool // Each row / column is unique
	allowBacktracking bool // Full backtracking
}

// Rules for difficulty
// Easy 		-> No triplets
// Medium		-> Easy + row / column balance
// Hard			-> Medium + row / column unique
// Extreme	-> Hard untill no longer possible then backtracking

// solveDifficulty tries to solveDifficulty a field with the constraints
// of the given difficulty
// It returns the solved board and a bool indicating if it was possible
func solveDifficulty(size int, board []uint8, diff Difficulty) ([]uint8, bool) {
	var rules solveMode

	switch diff {
	case Easy:
		rules = solveMode{
			allowTriplets: true,
		}

	case Medium:
		rules = solveMode{
			allowTriplets: true,
			allowBalance:  true,
		}

	case Hard:
		rules = solveMode{
			allowTriplets:   true,
			allowBalance:    true,
			allowUniqueness: true,
		}

	default:
		rules = solveMode{
			allowTriplets:     true,
			allowBalance:      true,
			allowUniqueness:   true,
			allowBacktracking: true,
		}
	}

	solution := slices.Clone(board)
	solvable := solveWithRules(size, solution, rules)

	return solution, solvable
}

func solveWithRules(size int, board []uint8, mode solveMode) bool {
	g := newGrid(size, board)

	for {
		changed := false

		if mode.allowTriplets && applyTriplet(g) {
			changed = true
		}
		if mode.allowBalance && applyBalance(g) {
			changed = true
		}
		if mode.allowUniqueness && applyUniqueness(g) {
			changed = true
		}

		if g.solved() {
			return true
		}

		if !changed {
			break
		}
	}

	if !mode.allowBacktracking {
		return false
	}

	return solveByBacktracking(g)
}

var tripletOffsets = [][2]int{
	{-1, -2},
	{1, 2},
	{-1, 1},
}

func applyTriplet(g grid) bool {
	if g.size <= 2 {
		return false
	}

	changed := false

	for idx := range g.size {
		if applyTripletsToLine(g.row(idx)) {
			changed = true
		}
		if applyTripletsToLine(g.col(idx)) {
			changed = true
		}
	}

	return changed
}

func applyTripletsToLine(l line) bool {
	changed := false

	for pos := range l.len() {
		if l.at(pos) != unknown {
			continue
		}

		for _, offset := range tripletOffsets {
			other1 := pos + offset[0]
			other2 := pos + offset[1]
			if !l.inBounds(other1) || !l.inBounds(other2) {
				continue
			}

			value1 := l.at(other1)
			value2 := l.at(other2)
			if value1 == unknown || value2 == unknown || value1 != value2 {
				continue
			}

			l.set(pos, 1-value1)

			changed = true
			break
		}
	}

	return changed
}

func applyBalance(g grid) bool {
	changed := false

	for idx := range g.size {
		if applyBalanceToLine(g.row(idx)) {
			changed = true
		}
		if applyBalanceToLine(g.col(idx)) {
			changed = true
		}
	}

	return changed
}

func applyBalanceToLine(l line) bool {
	fillValue, ok := l.counts().completionValue(l.len() / 2)
	if !ok {
		return false
	}

	changed := false
	for _, pos := range l.unknownPositions() {
		l.set(pos, fillValue)
		changed = true
	}

	return changed
}

func applyUniqueness(g grid) bool {
	// This only handles cases where a row / column is unique
	// with another one apart from 1 or 2 unfilled cells

	changed := applyUniquenessOnAxis(g, rowAxis)
	if applyUniquenessOnAxis(g, colAxis) {
		changed = true
	}
	return changed
}

func applyUniquenessOnAxis(g grid, axis axis) bool {
	type partialLine struct {
		line      line
		values    []uint8
		emptyIdxs []int
	}

	fullLines := make([][]uint8, 0, g.size)
	partials := make([]partialLine, 0, g.size)

	for idx := range g.size {
		var current line
		if axis == rowAxis {
			current = g.row(idx)
		} else {
			current = g.col(idx)
		}

		emptyIdxs := current.unknownPositions()
		switch len(emptyIdxs) {
		case 0:
			fullLines = append(fullLines, current.values())
		case 1, 2:
			partials = append(partials, partialLine{
				line:      current,
				values:    current.values(),
				emptyIdxs: emptyIdxs,
			})
		}
	}

	changed := false
	for _, partial := range partials {
		for _, full := range fullLines {
			if !linesMatchKnown(partial.values, full) {
				continue
			}

			for _, idx := range partial.emptyIdxs {
				partial.line.set(idx, 1-full[idx])
			}

			changed = true
			break
		}
	}

	return changed
}

func solveByBacktracking(g grid) bool {
	var solve func(idx int) bool
	solve = func(idx int) bool {
		if idx == g.size*g.size {
			return validBoard(g)
		}

		if g.cells[idx] != unknown {
			return solve(idx + 1)
		}

		if canPlaceCell(g, idx, 0) {
			g.cells[idx] = 0
			if solve(idx + 1) {
				return true
			}
		}

		if canPlaceCell(g, idx, 1) {
			g.cells[idx] = 1
			if solve(idx + 1) {
				return true
			}
		}

		g.cells[idx] = unknown
		return false
	}

	return solve(0)
}

func canPlaceCell(g grid, idx int, value uint8) bool {
	rowIdx := idx / g.size
	colIdx := idx % g.size

	return canPlaceInLine(g.row(rowIdx), colIdx, value) &&
		canPlaceInLine(g.col(colIdx), rowIdx, value)
}

func canPlaceInLine(l line, pos int, value uint8) bool {
	for _, offset := range tripletOffsets {
		other1 := pos + offset[0]
		other2 := pos + offset[1]
		if !l.inBounds(other1) || !l.inBounds(other2) {
			continue
		}

		if l.at(other1) == value && l.at(other2) == value {
			return false
		}
	}

	counts := l.counts()
	counts.unknowns--
	counts.add(value)

	return counts.canStillBalance(l.len() / 2)
}

func validBoard(g grid) bool {
	// Check for uniqueness

	// Row
	for i := range g.size {
		rowI := g.row(i).values()

		for j := i + 1; j < g.size; j++ {
			rowJ := g.row(j).values()

			if slices.Equal(rowI, rowJ) {
				return false
			}
		}
	}

	// Col
	cols := make([][]uint8, 0, g.size)
	for colIdx := range g.size {
		cols = append(cols, g.col(colIdx).values())
	}

	for i := range cols {
		for j := i + 1; j < len(cols); j++ {
			if slices.Equal(cols[i], cols[j]) {
				return false
			}
		}
	}

	return true
}

// solveUnique determines if a field leads to an unique solution
func solveUnique(size int, board []uint8) bool {
	g := newGrid(size, board)

	solutions := 0

	var solve func(idx int)
	solve = func(idx int) {
		if solutions > 1 {
			return
		}

		if idx == g.size*g.size {
			if validBoard(g) {
				solutions++
			}

			return
		}

		if g.cells[idx] != unknown {
			solve(idx + 1)
			return
		}

		if canPlaceCell(g, idx, 0) {
			g.cells[idx] = 0
			solve(idx + 1)
		}

		if canPlaceCell(g, idx, 1) {
			g.cells[idx] = 1
			solve(idx + 1)
		}

		g.cells[idx] = unknown
	}

	solve(0)

	return solutions == 1
}
