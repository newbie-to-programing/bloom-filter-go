package bloom

type BloomFilter struct {
	k      int
	bitset []byte
}

func NewBloomFilter(k, m int) *BloomFilter {
	numBytes := (m + 7) / 8

	return &BloomFilter{
		k:      k,
		bitset: make([]byte, numBytes),
	}
}
