package crossview

import (
	"slices"
)

type Difficulty int

type difficulty = Difficulty

const (
	Easy Difficulty = iota + 1
	Medium
	Hard
	Expert
)

type solveMode struct {
	allowSaturation           bool // Add blocked at the end if a cell already sees the required amount
	allowExactFill            bool // If a cell still need X cells and there are only X possible options left
	allowSingleDirection      bool // Only one direction left and the cell still needs more cells
	allowSeesNothing          bool // A cell surrounded by blocked cells
	allowCandidateElimination bool // Check if adding a cell in a valid position would ruin some connected cell
	allowForcedPrefix         bool // "Fill" all directions but one, if there are required cells left add them there
	allowForcedPrefixGlobal   bool // Same as allowImaginaryFill except it also takes into account that the imaginary fill doesn't ruin any other cells
	allowBacktracking         bool // Full backtracking
}

func solveDifficulty(size int, board []uint8, diff Difficulty) ([]uint8, bool) {
	var rules solveMode

	switch diff {
	case Easy:
		rules = solveMode{
			allowSaturation:      true,
			allowExactFill:       true,
			allowSingleDirection: true,
			allowSeesNothing:     true,
		}

	case Medium:
		rules = solveMode{
			allowSaturation:           true,
			allowExactFill:            true,
			allowSingleDirection:      true,
			allowSeesNothing:          true,
			allowCandidateElimination: true,
			allowForcedPrefix:         true,
		}

	case Hard:
		rules = solveMode{
			allowSaturation:           true,
			allowExactFill:            true,
			allowSingleDirection:      true,
			allowSeesNothing:          true,
			allowCandidateElimination: true,
			allowForcedPrefix:         true,
			allowForcedPrefixGlobal:   true,
		}

	default:
		rules = solveMode{
			allowSaturation:           true,
			allowExactFill:            true,
			allowSingleDirection:      true,
			allowSeesNothing:          true,
			allowCandidateElimination: true,
			allowForcedPrefix:         true,
			allowForcedPrefixGlobal:   true,
			allowBacktracking:         true,
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

		if mode.allowSaturation && applySaturation(g) {
			changed = true
		}
		if mode.allowExactFill && applyExactFill(g) {
			changed = true
		}
		if mode.allowSingleDirection && applySingleDirection(g) {
			changed = true
		}
		if mode.allowSeesNothing && applySeesNothing(g) {
			changed = true
		}
		if mode.allowCandidateElimination && applyCandidateElimination(g) {
			changed = true
		}
		if mode.allowForcedPrefix && applyForcedPrefix(g) {
			changed = true
		}
		if mode.allowForcedPrefixGlobal && applyForcedPrefixGlobal(g) {
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

type dir struct {
	axis axis
	left bool
	open int
}

var dirs = [4]dir{
	{axis: rowAxis, left: true},
	{axis: rowAxis, left: false},
	{axis: colAxis, left: true},
	{axis: colAxis, left: false},
}

func getDirs(g grid, rowIdx, colIdx int) []dir {
	open := func(d dir) int {
		var line line
		var idx int
		if d.axis == rowAxis {
			line = g.row(rowIdx)
			idx = colIdx
		} else {
			line = g.col(colIdx)
			idx = rowIdx
		}

		amount := 0
		if d.left {
			amount = line.openLeft(idx)
		} else {
			amount = line.openRight(idx)
		}

		return amount
	}

	dirOpens := make([]dir, 0, len(dirs))
	for _, d := range dirs {
		d.open = open(d)
		dirOpens = append(dirOpens, d)
	}

	slices.SortFunc(dirOpens, func(a, b dir) int { return a.open - b.open })

	return dirOpens
}

// getLine returns the line and index corresponding to a direction
func (d dir) getLine(g grid, rowIdx, colIdx int) (line, int) {
	var line line
	var idx int
	if d.axis == rowAxis {
		line = g.row(rowIdx)
		idx = colIdx
	} else {
		line = g.col(colIdx)
		idx = rowIdx
	}

	return line, idx
}

func applySaturation(g grid) bool {
	changed := false

	for r := range g.size {
		for c := range g.size {
			cell := g.at(r, c)

			// Only look at numbers
			if cell >= blocked {
				continue
			}

			// Cell isn't full yet
			if uint8(g.sees(r, c)) != cell {
				continue
			}

			// Cell is full
			// Add boundaries
			if g.addBoundary(r, c) {
				changed = true
			}
		}
	}

	return changed
}

func applyExactFill(g grid) bool {
	changed := false

	for r := range g.size {
		for c := range g.size {
			cell := g.at(r, c)

			// Only look at numbers
			if cell >= blocked {
				continue
			}

			sees := g.sees(r, c)

			// Cell already full
			if uint8(sees) == cell {
				continue
			}

			open := g.open(r, c)

			// Too many open spaces
			if uint8(open+sees) != cell {
				continue
			}

			g.fill(r, c)
			changed = true
		}
	}

	return changed
}

func applySingleDirection(g grid) bool {
	changed := false

	for r := range g.size {
		for c := range g.size {
			cell := g.at(r, c)
			// We're only interested in numbers
			if cell >= blocked {
				continue
			}

			// Already full
			sees := uint8(g.sees(r, c))
			if sees == cell {
				continue
			}

			dirs := getDirs(g, r, c)

			// Is there only one direction left with open spots?
			if dirs[3].open == 0 || dirs[2].open > 0 {
				continue
			}

			line, idx := dirs[3].getLine(g, r, c)

			if dirs[3].left {
				seesLeft := line.seesLeft(idx)
				line.fillLeft(idx, int(cell-sees)+seesLeft)
			} else {
				seesRight := line.seesRight(idx)
				line.fillRight(idx, int(cell-sees)+seesRight)
			}

			changed = true
		}
	}

	return changed
}

func applySeesNothing(g grid) bool {
	changed := false

	for r := range g.size {
		for c := range g.size {
			cell := g.at(r, c)
			// Only interested in unknown cells
			if cell != unknown {
				continue
			}

			// Check all sides for a block
			// Wet
			row := g.row(r)
			if row.inBounds(c-1) && row.at(c-1) != blocked {
				continue
			}
			if row.inBounds(c+1) && row.at(c+1) != blocked {
				continue
			}

			col := g.col(c)
			if col.inBounds(r-1) && col.at(r-1) != blocked {
				continue
			}
			if col.inBounds(r+1) && col.at(r+1) != blocked {
				continue
			}

			g.set(r, c, blocked)
			changed = true
		}
	}

	return changed
}

func applyCandidateElimination(g grid) bool {
	changed := false

	for r := range g.size {
		for c := range g.size {
			if g.at(r, c) != unknown {
				continue
			}

			// Try filled
			g.set(r, c, filled)

			if !g.valid() {
				g.set(r, c, blocked)
				changed = true

				continue
			}

			// Try blocked
			g.set(r, c, blocked)
			if !g.valid() {
				g.set(r, c, filled)
				changed = true

				continue
			}

			g.set(r, c, unknown)
		}
	}

	return changed
}

func applyForcedPrefix(g grid) bool {
	changed := false

	for r := range g.size {
		for c := range g.size {
			// Skip non numbers
			cell := g.at(r, c)
			if cell >= blocked {
				continue
			}

			sees := uint8(g.sees(r, c))
			// Full
			if sees == cell {
				continue
			}

			dirs := getDirs(g, r, c)

			// If we fill the smallest open directions, do we still have some left for the last one
			otherMax := dirs[0].open + dirs[1].open + dirs[2].open
			cellsLeft := int(cell-sees) - otherMax
			if cellsLeft <= 0 {
				continue
			}

			// We have some left
			// Fill them in
			line, idx := dirs[3].getLine(g, r, c)

			var filled int
			if dirs[3].left {
				// Compensate for the cells already there
				cellsLeft += line.seesLeft(idx)
				filled = line.fillLeft(idx, cellsLeft)
			} else {
				// Compensate for the cells already there
				cellsLeft += line.seesRight(idx)
				filled = line.fillRight(idx, cellsLeft)
			}

			if filled > 0 {
				changed = true
			}
		}
	}

	return changed
}

func applyForcedPrefixGlobal(g grid) bool {
	changed := false

	for r := range g.size {
		for c := range g.size {
			// Skip non numbers
			cell := g.at(r, c)
			if cell >= blocked {
				continue
			}

			// Full
			sees := uint8(g.sees(r, c))
			if sees == cell {
				continue
			}

			dirs := getDirs(g, r, c)

			// If we fill the smallest open directions, do we still have some left for the last one
			otherMax := dirs[0].open + dirs[1].open + dirs[2].open
			cellsLeft := int(cell-sees) - otherMax

			// Keep track of the cells we fill to undo later on
			filledCells := make([][2]int, 0, otherMax)

			// Let start filling in the cells and validating the board every step
			for _, d := range dirs[:3] {
				// Get the correct line and position index
				line, idx := d.getLine(g, r, c)

				// Get the offset from any already filled in cells
				var offset int
				var getPos func(int) int
				if d.left {
					offset = line.seesLeft(idx)
					getPos = func(pos int) int { return idx - pos }
				} else {
					offset = line.seesRight(idx)
					getPos = func(pos int) int { return idx + pos }
				}

				for i := range d.open {
					pos := getPos(i + offset)
					// Don't do anything if the cell is already filled
					if line.at(pos) != unknown {
						continue
					}

					// Try to set it to fill
					line.set(pos, filled)

					// Grid is no longer valid
					// Set it back to unknown and increase the cellsLeft
					// to add to the last direction and move on to the next direction
					if !g.valid() {
						line.set(pos, unknown)
						cellsLeft++
						break
					}

					// Keep track of the cells we need to restore
					var row int
					var col int
					if d.axis == rowAxis {
						row = r
						col = pos
					} else {
						row = pos
						col = c
					}

					filledCells = append(filledCells, [2]int{row, col})
				}
			}

			// We can use the 3 directions to fill everything
			// No cells left to add to the last direction
			if cellsLeft <= 0 {
				continue
			}

			// Add the remaining cells
			line, idx := dirs[3].getLine(g, r, c)

			var filled int
			if dirs[3].left {
				// Compensate for the cells already there
				cellsLeft += line.seesLeft(idx)
				filled = line.fillLeft(idx, cellsLeft)
			} else {
				// Compensate for the cells already there
				cellsLeft += line.seesRight(idx)
				filled = line.fillRight(idx, cellsLeft)
			}

			if filled > 0 {
				changed = true
			}

			// Restore the grid
			for _, filled := range filledCells {
				g.set(filled[0], filled[1], unknown)
			}
		}
	}

	return changed
}

func solveByBacktracking(g grid) bool {
	var solve func(idx int) bool
	solve = func(idx int) bool {
		if idx == g.size*g.size {
			return g.solved() && g.valid()
		}

		if g.cells[idx] != unknown {
			return solve(idx + 1)
		}

		// Try filled in
		g.cells[idx] = filled
		if g.valid() {
			if solve(idx + 1) {
				return true
			}
		}

		// Try blocked
		g.cells[idx] = blocked
		if g.valid() {
			if solve(idx + 1) {
				return true
			}
		}

		g.cells[idx] = unknown
		return false
	}

	return solve(0)
}

func solveUnique(g grid) bool {
	solutions := 0

	var solve func(idx int)
	solve = func(idx int) {
		if solutions > 1 {
			return
		}

		if idx == g.size*g.size {
			if g.solved() && g.valid() {
				solutions++
			}

			return
		}

		if g.cells[idx] != unknown {
			solve(idx + 1)
			return
		}

		// Try filled in
		g.cells[idx] = filled
		if g.valid() {
			solve(idx + 1)
		}

		// Try blocked
		g.cells[idx] = blocked
		if g.valid() {
			solve(idx + 1)
		}

		g.cells[idx] = unknown
	}

	solve(0)

	return solutions == 1
}
