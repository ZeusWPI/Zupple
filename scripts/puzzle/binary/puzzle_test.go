package binary

import (
	"strings"
	"testing"
)

func TestMaskMatchesDifficulty(t *testing.T) {
	b := Binary{
		Size: 4,
		Field: []uint8{
			1, 0, 0, 1,
			0, 0, 1, 1,
			1, 1, 0, 0,
			0, 1, 1, 0,
		},
	}

	g := newGrid(b.Size, b.Field)

	easyMask := []uint8{
		0, 1, 1, 1,
		1, 0, 1, 1,
		1, 0, 0, 1,
		0, 1, 1, 0,
	}

	mediumMask := []uint8{
		0, 1, 1, 1,
		0, 1, 1, 0,
		1, 0, 0, 0,
		0, 1, 1, 1,
	}

	if !g.maskMatchesDifficulty(easyMask, Easy) {
		t.Error("easy mask should match easy difficulty")
	}

	if g.maskMatchesDifficulty(easyMask, Medium) {
		t.Error("easy mask should not match medium difficulty")
	}

	if !g.maskMatchesDifficulty(mediumMask, Medium) {
		t.Error("medium mask should match medium difficulty")
	}

	if g.maskMatchesDifficulty(mediumMask, Easy) {
		t.Error("medium mask should not match easy difficulty")
	}
}

func TestGenerateMask(t *testing.T) {
	b := Binary{
		Size: 4,
		Field: []uint8{
			1, 0, 0, 1,
			0, 0, 1, 1,
			1, 1, 0, 0,
			0, 1, 1, 0,
		},
	}

	mask, err := b.Puzzle(Easy)
	if err != nil {
		t.Fatalf("unexpected error generating easy mask: %v", err)
	}

	if filledRatio(mask) < defaultFilledPercentageByDifficulty[Easy] {
		t.Fatalf("expected filled ratio >= %f, got %f", defaultFilledPercentageByDifficulty[Easy], filledRatio(mask))
	}

	if filledRatio(mask) != 10.0/16.0 {
		t.Fatalf("expected easy ratio to stay as close as possible to the minimum, got %f", filledRatio(mask))
	}

	if !newGrid(b.Size, b.Field).maskMatchesDifficulty(mask, Easy) {
		t.Fatal("generated mask should match easy difficulty")
	}

	mask, err = b.Puzzle(Medium, 0.55)
	if err != nil {
		t.Fatalf("unexpected error generating medium mask: %v", err)
	}

	if filledRatio(mask) < 0.55 {
		t.Fatalf("expected filled ratio >= 0.55, got %f", filledRatio(mask))
	}

	if filledRatio(mask) != 9.0/16.0 {
		t.Fatalf("expected medium ratio to stay as close as possible to the minimum, got %f", filledRatio(mask))
	}

	if !newGrid(b.Size, b.Field).maskMatchesDifficulty(mask, Medium) {
		t.Fatal("generated mask should match medium difficulty")
	}
}

func TestPuzzleString(t *testing.T) {
	b := Binary{
		Size: 4,
		Field: []uint8{
			1, 0, 0, 1,
			0, 0, 1, 1,
			1, 1, 0, 0,
			0, 1, 1, 0,
		},
	}

	mask := []uint8{
		1, 0, 1, 0,
		0, 1, 0, 1,
		1, 0, 1, 0,
		0, 1, 0, 1,
	}

	out, err := b.PuzzleString(mask)
	if err != nil {
		t.Fatalf("unexpected error rendering puzzle: %v", err)
	}

	if strings.Contains(out, "2") {
		t.Fatal("rendered puzzle should not contain unknown markers")
	}

	if !strings.Contains(out, "|   |") {
		t.Fatal("rendered puzzle should show spaces for hidden cells")
	}
}
