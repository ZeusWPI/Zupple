package binary

// Compress compresses the field to a byte array
// The byte array is constructed by putting each row after another
// and setting the bit to 1 if the value is 1, 0 otherwise.
// This also means that it is least significant bit first
func (b *Binary) Compress() []byte {
	return newGrid(b.Size, b.Field).compress()
}

func (b *Binary) decompress(data []byte) error {
	b.Field = make([]uint8, b.Size*b.Size)
	return newGrid(b.Size, b.Field).loadCompressed(data)
}
