package binary

import (
	"errors"
	"math"
	"math/rand"
	"slices"
)

var (
	errInvalidDifficulty       = errors.New("invalid difficulty")
	errInvalidFilledPercentage = errors.New("filled percentage must be between 0 and 1")
	errTooManyFilledArguments  = errors.New("at most one filled percentage can be provided")
	errInvalidMaskLength       = errors.New("invalid mask length for binary size")
	errInvalidMaskValue        = errors.New("mask values must be 0 or 1")
	errUnreachableMask         = errors.New("unable to generate a mask that matches the requested difficulty and filled percentage")
)

var defaultFilledPercentageByDifficulty = map[Difficulty]float64{
	Easy:   0.60,
	Medium: 0.55,
	Hard:   0.45,
	Expert: 0.35,
}

func (d Difficulty) valid() bool {
	switch d {
	case Easy, Medium, Hard, Expert:
		return true
	default:
		return false
	}
}

func (d Difficulty) lower() (Difficulty, bool) {
	switch d {
	case Medium:
		return Easy, true
	case Hard:
		return Medium, true
	case Expert:
		return Hard, true
	default:
		return 0, false
	}
}

// generateMask returns a mask for a grid with a difficulty and optional filled percentage
func (g grid) generateMask(diff Difficulty, filledPercentage ...float64) ([]uint8, error) {
	if !diff.valid() {
		return nil, errInvalidDifficulty
	}

	minimumFilledPercentage, err := resolveFilledPercentage(diff, filledPercentage)
	if err != nil {
		return nil, err
	}

	cellCount := g.size * g.size
	minimumFilledCells := minFilledCells(cellCount, minimumFilledPercentage)

	var bestMask []uint8
	bestFilled := cellCount + 1

	// Try multiple times to get an optimimum problem
	for range 128 {
		mask := fullMask(cellCount)
		filled := cellCount

		for _, idx := range rand.Perm(cellCount) {
			if filled-1 < minimumFilledCells {
				break
			}

			mask[idx] = 0
			if !g.maskSolvable(mask, diff) {
				mask[idx] = 1
				continue
			}

			filled--
		}

		// Test if it is solvable with a lower difficulty
		if !g.maskMatchesDifficulty(mask, diff) {
			continue
		}

		if filled < bestFilled {
			bestMask = slices.Clone(mask)
			bestFilled = filled
		}

		// If the fill count == the desired fill count then we can stop searching
		if filled == minimumFilledCells {
			return mask, nil
		}
	}

	if bestMask == nil {
		return fullMask(cellCount), errUnreachableMask
	}

	return bestMask, nil
}

func (g grid) maskMatchesDifficulty(mask []uint8, diff Difficulty) bool {
	lower, ok := diff.lower()
	if !ok {
		return g.maskSolvable(mask, diff)
	}

	return !g.maskSolvable(mask, lower)
}

func (g grid) maskSolvable(mask []uint8, diff Difficulty) bool {
	board, err := g.maskedBoard(mask)
	if err != nil {
		return false
	}

	solution, solvable := solveDifficulty(g.size, board, diff)
	if !solvable {
		return false
	}
	if !solveUnique(g.size, board) {
		return false
	}

	return slices.Equal(g.cells, solution)
}

// maskedBoard creates a board with values set unknown according to a mask
func (g grid) maskedBoard(mask []uint8) ([]uint8, error) {
	if err := g.validateMask(mask); err != nil {
		return nil, err
	}

	board := make([]uint8, len(g.cells))
	for idx, visible := range mask {
		if visible == 1 {
			board[idx] = g.cells[idx]
			continue
		}

		board[idx] = unknown
	}

	return board, nil
}

// filledRatio returns the percentage that a mask is set to 1
func filledRatio(mask []uint8) float64 {
	filled := 0
	for _, value := range mask {
		if value == 1 {
			filled++
		}
	}

	return float64(filled) / float64(len(mask))
}

// fullMask returns a mask with all cells set to 1
func fullMask(cellCount int) []uint8 {
	mask := make([]uint8, cellCount)
	for idx := range mask {
		mask[idx] = 1
	}

	return mask
}

func resolveFilledPercentage(diff Difficulty, filledPercentage []float64) (float64, error) {
	if len(filledPercentage) > 1 {
		return 0, errTooManyFilledArguments
	}

	if len(filledPercentage) == 0 {
		return defaultFilledPercentageByDifficulty[diff], nil
	}

	value := filledPercentage[0]
	if value < 0 || value > 1 {
		return 0, errInvalidFilledPercentage
	}

	return value, nil
}

func minFilledCells(cellCount int, filledPercentage float64) int {
	return int(math.Ceil(filledPercentage * float64(cellCount)))
}
