package binary

import "testing"

func TestSolveTripletValid(t *testing.T) {
	board := []uint8{
		2, 0, 0, 1,
		0, 2, 1, 1,
		1, 2, 2, 0,
		2, 1, 1, 2,
	}

	if !solveWithRules(4, board, solveMode{allowTriplets: true}) {
		t.Error("Board should be solvable with only triplets")
	}
}

func TestSolveTripletInValid(t *testing.T) {
	board := []uint8{
		2, 0, 0, 1,
		0, 2, 2, 1,
		1, 2, 2, 0,
		2, 1, 1, 2,
	}

	if solveWithRules(4, board, solveMode{allowTriplets: true}) {
		t.Error("Board should not be solvable with only triplets")
	}
}

func TestSolveBalanceValid(t *testing.T) {
	board := []uint8{
		2, 0, 0, 1,
		2, 0, 1, 2,
		1, 2, 2, 2,
		2, 1, 1, 0,
	}

	if !solveWithRules(4, board, solveMode{allowBalance: true}) {
		t.Error("Board should be solvable with only balance")
	}
}

func TestSolveBalanceInValid(t *testing.T) {
	board := []uint8{
		2, 0, 0, 1,
		2, 2, 1, 2,
		2, 2, 2, 2,
		2, 1, 1, 0,
	}

	if solveWithRules(4, board, solveMode{allowBalance: true}) {
		t.Error("Board should not be solvable with only balance")
	}
}

func TestSolveUniquenessValid(t *testing.T) {
	board := []uint8{
		1, 0, 0, 1,
		1, 0, 2, 2,
		0, 1, 2, 2,
		0, 1, 0, 1,
	}

	if !solveWithRules(4, board, solveMode{allowUniqueness: true}) {
		t.Error("Board should be solvable with only uniqueness")
	}
}

func TestSolveUniquenessInValid(t *testing.T) {
	board := []uint8{
		1, 0, 0, 1,
		1, 2, 2, 2,
		0, 1, 2, 2,
		0, 1, 0, 1,
	}

	if solveWithRules(4, board, solveMode{allowUniqueness: true}) {
		t.Error("Board should not be solvable with only uniqueness")
	}
}

func TestSolveBacktrackingValid(t *testing.T) {
	board := []uint8{
		2, 2, 0, 1,
		2, 2, 1, 1,
		2, 2, 2, 0,
		2, 2, 2, 0,
	}

	if !solveByBacktracking(newGrid(4, board)) {
		t.Error("Board should be solvable with backtracking")
	}
}

func TestSolveBacktrackingInValid(t *testing.T) {
	board := []uint8{
		1, 2, 0, 1,
		0, 1, 1, 1,
		0, 1, 1, 0,
		1, 1, 0, 0,
	}

	if solveByBacktracking(newGrid(4, board)) {
		t.Error("Board should not be solvable with backtracking")
	}
}

func TestSolveUniqueValid(t *testing.T) {
	board := []uint8{
		1, 2, 0, 1,
		2, 2, 1, 1,
		2, 2, 0, 0,
		2, 2, 2, 0,
	}

	if !solveUnique(4, board) {
		t.Error("Board should have an unique solution")
	}
}

func TestSolveUniqueInValid(t *testing.T) {
	board := []uint8{
		2, 2, 0, 1,
		2, 2, 1, 1,
		2, 2, 2, 0,
		2, 2, 2, 0,
	}

	if solveUnique(4, board) {
		t.Error("Board should not have an unique solution")
	}
}
