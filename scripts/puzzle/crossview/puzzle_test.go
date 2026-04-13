package crossview

import "testing"

func TestMaskSolvableAcceptsHiddenNumberedCells(t *testing.T) {
	g := newGrid(4, []uint8{
		blocked, 1, blocked, 3,
		5, 4, 5, 6,
		2, blocked, 3, 4,
		5, 3, 5, 6,
	})
	mask := []uint8{
		1, 0, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
		1, 1, 1, 1,
	}

	if !g.maskSolvable(mask, Expert) {
		t.Fatal("hidden numbered cells should be solvable as filled cells")
	}
}
