package bloom

import (
	"hash/fnv"
	"math"
	"sync"
)

type BloomFilter struct {
	mu     sync.RWMutex
	k      uint
	m      uint
	bitset []byte
}

func New(m, k uint) *BloomFilter {
	numBytes := (m + 7) / 8

	return &BloomFilter{
		k:      k,
		m:      m,
		bitset: make([]byte, numBytes),
	}
}

func (bf *BloomFilter) Add(data []byte) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	indices := bf.getIndices(data)

	for idx := range indices {
		targetByte := idx / 8
		targetBit := idx % 8

		bf.bitset[targetByte] |= 1 << targetBit
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	bf.mu.RLock()
	defer bf.mu.RUnlock()

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

func EstimateParameters(n uint, p float64) (uint, uint) {
	// m = - (n * ln(p)) / (ln(2)^2)
	mFloat := -float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)
	m := uint(math.Ceil(mFloat))

	// k = (m / n) * ln(2)
	kFloat := (float64(m) / float64(n)) * math.Log(2)
	k := uint(math.Ceil(kFloat))

	return m, k
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
