package binary

import "testing"

func TestNewSizeInvalid(t *testing.T) {
	test := func(size int) {
		_, err := New(size)
		if err == nil {
			t.Errorf("Size %d should not be accepted", size)
		}
	}

	// Too small
	test(1)
	test(0)
	test(-1)

	// Uneven
	test(3)
	test(5)
	test(17)

	// Too big
	test(33)
	test(34)
	test(102)
}

func TestEqualTrue(t *testing.T) {
	b1 := Binary{
		Size: 2,
		Field: []uint8{
			0, 1,
			1, 0,
		},
	}

	b2 := Binary{
		Size: 2,
		Field: []uint8{
			0, 1,
			1, 0,
		},
	}

	if !b1.Equal(b2) {
		t.Error("Should be equal")
	}
}

func TestEqualFalseSize(t *testing.T) {
	b1 := Binary{
		Size: 2,
		Field: []uint8{
			0, 1,
			1, 0,
		},
	}

	b2 := Binary{
		Size: 3,
		Field: []uint8{
			0, 1,
			1, 0,
		},
	}

	if b1.Equal(b2) {
		t.Error("Should not be equal")
	}
}

func TestEqualFalseFieldLength(t *testing.T) {
	b1 := Binary{
		Size: 2,
		Field: []uint8{
			0, 1,
			1, 0,
		},
	}

	b2 := Binary{
		Size: 2,
		Field: []uint8{
			0, 1,
			1, 0,
			1,
		},
	}

	if b1.Equal(b2) {
		t.Error("Should not be equal")
	}
}

func TestEqualFieldValue(t *testing.T) {
	b1 := Binary{
		Size: 2,
		Field: []uint8{
			0, 1,
			1, 0,
		},
	}

	b2 := Binary{
		Size: 2,
		Field: []uint8{
			0, 1,
			1, 1,
		},
	}

	if b1.Equal(b2) {
		t.Error("Should not be equal")
	}
}
