package binary

import (
	"slices"
	"testing"
)

func TestCompress(t *testing.T) {
	b := Binary{
		Size: 2,
		Field: []uint8{
			0, 1,
			1, 0,
		},
	}

	compressed := b.Compress()
	correct := []byte{6}

	if !slices.Equal(compressed, correct) {
		t.Errorf("Expected %+v, got %+v", correct, compressed)
	}

	b = Binary{
		Size: 4,
		Field: []uint8{
			1, 0, 0, 1,
			0, 1, 1, 0,
			0, 1, 0, 1,
			1, 0, 1, 0,
		},
	}

	compressed = b.Compress()
	correct = []byte{105, 90}

	if !slices.Equal(compressed, correct) {
		t.Errorf("Expected %+v, got %+v", correct, compressed)
	}
}

func TestDecompress(t *testing.T) {
	compressed := []byte{6}
	b, err := NewCompressed(2, compressed)
	if err != nil {
		t.Errorf("Unexpected error when creating a compressed puzzle %v", err)
	}

	field := []uint8{
		0, 1,
		1, 0,
	}

	if !slices.Equal(b.Field, field) {
		t.Errorf("Expected %+v, got %+v", field, b.Field)
	}

	compressed = []byte{105, 90}
	b, err = NewCompressed(4, compressed)
	if err != nil {
		t.Errorf("Unexpected error when creating a compressed puzzle %v", err)
	}

	field = []uint8{
		1, 0, 0, 1,
		0, 1, 1, 0,
		0, 1, 0, 1,
		1, 0, 1, 0,
	}

	if !slices.Equal(b.Field, field) {
		t.Errorf("Expected %+v, got %+v", field, b.Field)
	}
}

func TestCompressDecompress(t *testing.T) {
	// Test a bunch of random puzzles to see if they compare
	sizes := map[int]int{
		2:  10,
		4:  20,
		6:  50,
		8:  80,
		10: 100,
	}

	for size, amount := range sizes {
		for range amount {
			b1, err := New(size)
			if err != nil {
				t.Errorf("Unexpected error when creating a puzzle %v", err)
			}

			compressed := b1.Compress()

			b2, err := NewCompressed(size, compressed)
			if err != nil {
				t.Errorf("Unexpected error when creating a compressed puzzle %v", err)
			}

			if !slices.Equal(b1.Field, b2.Field) {
				t.Errorf("Expected %+v, got after compression %v", b1.Field, b2.Field)
			}
		}
	}
}
