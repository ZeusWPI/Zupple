package crossview

func (c *Crossview) generate() error {
	g := emptyGrid(c.size)

	// Set blocked cells
	for r := range g.size {
		for c := range g.size {
			if chance(0.3) {
				g.set(r, c, blocked)
			}
		}
	}

	// Recalculate all cells
	for r := range c.size {
		for c := range c.size {
			if g.at(r, c) == blocked {
				continue
			}

			sees := uint8(g.sees(r, c))
			if sees == 0 {
				sees = blocked
			}

			g.set(r, c, sees)
		}
	}

	c.field = g.cells

	return nil
}
