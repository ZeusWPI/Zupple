package crossview

import (
	"testing"
)

func TestSolveSaturationValid(t *testing.T) {
	board := []uint8{
		unknown, 1, unknown, 3,
		5, 4, 5, 6,
		2, unknown, 3, 4,
		5, 3, 5, 6,
	}

	if !solveWithRules(4, board, solveMode{allowSaturation: true}) {
		t.Error("Board should be solvable with only saturation")
	}
}

func TestSolveSaturationInValid(t *testing.T) {
	board := []uint8{
		blocked, 1, blocked, 3,
		5, unknown, unknown, unknown,
		unknown, blocked, 3, unknown,
		unknown, 3, 5, 6,
	}

	if solveWithRules(4, board, solveMode{allowSaturation: true}) {
		t.Error("Board should not be solvable with only saturation")
	}
}

func TestSolveExactFillValid(t *testing.T) {
	board := []uint8{
		blocked, 1, blocked, 3,
		5, unknown, unknown, unknown,
		unknown, blocked, 3, unknown,
		unknown, 3, 5, 6,
	}

	if !solveWithRules(4, board, solveMode{allowExactFill: true}) {
		t.Error("Board should be solvable with only exact fill")
	}
}

func TestSolveExactFillInValid(t *testing.T) {
	board := []uint8{
		blocked, blocked, blocked, 3,
		blocked, unknown, blocked, 3,
		1, blocked, 2, 4,
		4, 3, 4, 6,
	}

	if solveWithRules(4, board, solveMode{allowExactFill: true}) {
		t.Error("Board should not be solvable with only exact fill")
	}
}

func TestSolveSingleDirectionValid(t *testing.T) {
	board := []uint8{
		blocked, 1, blocked, 3,
		unknown, unknown, unknown, 6,
		unknown, blocked, 3, 4,
		5, 3, 5, 6,
	}

	if !solveWithRules(4, board, solveMode{allowSingleDirection: true}) {
		t.Error("Board should be solvable with only single direction")
	}
}

func TestSolveSingleDirectionInValid(t *testing.T) {
	board := []uint8{
		blocked, blocked, blocked, 3,
		blocked, unknown, blocked, 3,
		1, blocked, 2, 4,
		4, 3, 4, 6,
	}

	if solveWithRules(4, board, solveMode{allowSingleDirection: true}) {
		t.Error("Board should not be solvable with only single direction")
	}
}

func TestSolveSeesNothingValid(t *testing.T) {
	board := []uint8{
		blocked, blocked, blocked, 3,
		blocked, unknown, blocked, 3,
		1, blocked, 2, 4,
		4, 3, 4, 6,
	}

	if !solveWithRules(4, board, solveMode{allowSeesNothing: true}) {
		t.Error("Board should be solvable with only sees nothing")
	}
}

func TestSolveSeesNothingInValid(t *testing.T) {
	board := []uint8{
		blocked, 1, blocked, 3,
		5, 4, 5, 6,
		2, unknown, 3, 4,
		5, 3, 5, 6,
	}

	if solveWithRules(4, board, solveMode{allowSeesNothing: true}) {
		t.Error("Board should not be solvable with only sees nothing")
	}
}

func TestSolveCandidateElimininationValid(t *testing.T) {
	board := []uint8{
		blocked, 1, blocked, 3,
		5, 4, 5, 6,
		2, unknown, 3, 4,
		5, 3, 5, 6,
	}

	if !solveWithRules(4, board, solveMode{allowCandidateElimination: true}) {
		t.Error("Board should be solvable with only candidate elimination")
	}
}

func TestSolveCandidateEliminationInValid(t *testing.T) {
}

func TestSolveForcedPrefixValid(t *testing.T) {
}

func TestSolveForcedPrefixInValid(t *testing.T) {
}

func TestSolveForcePrefixGlobalValid(t *testing.T) {
}

func TestSolveForcedPrefixGlobalInValid(t *testing.T) {
}

func TestSolveBacktrackingValid(t *testing.T) {
}

func TestSolveBacktrackingInValid(t *testing.T) {
}

func TestSolveUniqueValid(t *testing.T) {
}

func TestSolveUniqueInvalid(t *testing.T) {
}
