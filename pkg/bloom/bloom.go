package bloom

import "hash/fnv"

type BloomFilter struct {
	k      uint
	m      uint
	bitset []byte
}

func NewBloomFilter(k, m uint) *BloomFilter {
	numBytes := (m + 7) / 8

	return &BloomFilter{
		k:      k,
		m:      m,
		bitset: make([]byte, numBytes),
	}
}

func (bf *BloomFilter) Add(data []byte) {
	indices := bf.getIndices(data)

	for idx := range indices {
		targetByte := idx / 8
		targetBit := idx % 8

		bf.bitset[targetByte] |= 1 << targetBit
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	indices := bf.getIndices(data)

	for idx := range indices {
		targetByte := idx / 8
		targetBit := idx % 8

		exists := bf.bitset[targetByte]&(1<<targetBit) != 0

		if !exists {
			return false
		}
	}

	return true
}

func (bf *BloomFilter) getIndices(data []byte) []uint {
	h := fnv.New64a()
	h.Write(data)
	sum := h.Sum64()

	h1 := uint32(sum)
	h2 := uint32(sum >> 32)

	indices := make([]uint, bf.k)
	for i := uint(0); i < bf.k; i++ {
		index := (uint(h1) + i*uint(h2)) % bf.m
		indices[i] = index
	}

	return indices
}
