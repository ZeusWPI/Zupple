package crossview

func (c *Crossview) generate() error {
	g := emptyGrid(c.Size)

	// Set blocked cells
	for r := range g.size {
		for c := range g.size {
			if chance(0.3) {
				g.set(r, c, blocked)
			}
		}
	}

	// Recalculate all cells
	for r := range c.Size {
		for c := range c.Size {
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

	c.Field = g.cells

	return nil
}
