package crossview

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

// sees returns the total amount of connected filled cells to the pos
func (l line) sees(pos int) int {
	return l.seesLeft(pos) + l.seesRight(pos)
}

// seesLeft returns the amount of connected filled cells to the left of the pos
func (l line) seesLeft(pos int) int {
	if l.at(pos) == blocked {
		return 0
	}

	amount := 0

	for i := pos - 1; i >= 0; i-- {
		cell := l.at(i)
		if cell == blocked || cell == unknown {
			break
		}

		amount++
	}

	return amount
}

// seesRight returns the amount of connected filled cells to the right of the pos
func (l line) seesRight(pos int) int {
	if l.at(pos) == blocked {
		return 0
	}

	amount := 0

	for i := pos + 1; i < l.len(); i++ {
		cell := l.at(i)
		if cell == blocked || cell == unknown {
			break
		}

		amount++
	}

	return amount
}

// open returns the total amount of non blocked open cells
func (l line) open(pos int) int {
	return l.openLeft(pos) + l.openRight(pos)
}

// openLeft returns the amount of non blocked open cells to the left of the position
func (l line) openLeft(pos int) int {
	if l.at(pos) == blocked {
		return 0
	}

	amount := 0

	for i := pos - 1; i >= 0; i-- {
		cell := l.at(i)
		if cell == blocked {
			break
		}
		if cell != unknown {
			continue
		}

		amount++
	}

	return amount
}

// openRight returns the amount of non blocked open cells to the right of the position
func (l line) openRight(pos int) int {
	if l.at(pos) == blocked {
		return 0
	}

	amount := 0

	for i := pos + 1; i < l.len(); i++ {
		cell := l.at(i)
		if cell == blocked {
			break
		}
		if cell != unknown {
			continue
		}

		amount++
	}

	return amount
}

// addBoundary changes the first encountered unknown cell to a blocked cell on both sides
// It returns true if it added a boundary
func (l line) addBoundary(pos int) bool {
	if l.at(pos) == blocked {
		return false
	}

	changed := false

	// Before
	for i := pos - 1; i >= 0; i-- {
		cell := l.at(i)
		// There already is a blocked cell
		if cell == blocked {
			break
		}
		// Not an unknown cell
		if cell != unknown {
			continue
		}

		l.set(i, blocked)
		changed = true

		break
	}

	// After
	for i := pos + 1; i < l.len(); i++ {
		cell := l.at(i)
		if cell == blocked {
			break
		}
		if cell != unknown {
			continue
		}

		l.set(i, blocked)
		changed = true

		break
	}

	return changed
}

// fill fills any unknown cell
// It keeps going until it encounters a blocked cell
// Use fillLeft and fillRight for fine control
// It returns the amount of filled in cells
func (l line) fill(pos int) int {
	return l.fillLeft(pos, 0) + l.fillRight(pos, 0)
}

// fill fills any unknown cell to the left of the position
// It keeps going until the max is reached
// If the max is 0 then it keeps going until there's a blocked cell
// It does not overwrite any already non unknown cell but it does count it towards the max
// It returns the amount of filled in cells
func (l line) fillLeft(pos int, max int) int {
	if l.at(pos) == blocked {
		return 0
	}

	current := 0
	filledIn := 0

	for i := pos - 1; i >= 0; i-- {
		if max != 0 && current == max {
			break
		}

		cell := l.at(i)
		// Blocked cell, stop
		if cell == blocked {
			break
		}

		current++

		// Cell is already filled
		// Add it to the count
		if cell != unknown {
			continue
		}

		// Set as filled
		l.set(i, filled)
		filledIn++
	}

	return filledIn
}

// fill fills any unknown cell to the right of the position
// It keeps going until the max is reached or it encounters a blocked cell
// If the max is 0 then it keeps going until there's a blocked cell
// It does not overwrite any already non unknown cell but it does count it towards the max
// It returns the amount of filled in cells
func (l line) fillRight(pos int, max int) int {
	if l.at(pos) == blocked {
		return 0
	}

	current := 0
	filledIn := 0

	for i := pos + 1; i < l.len(); i++ {
		if max != 0 && current == max {
			break
		}

		cell := l.at(i)
		// Blocked cell, stop
		if cell == blocked {
			break
		}

		current++

		// Cell is already filled
		// Add it to the count
		if cell != unknown {
			continue
		}

		// Set as filled
		l.set(i, filled)
		filledIn++
	}

	return filledIn
}
