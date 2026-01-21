package bloom

import (
	_ "math"
	"testing"
)

// Stage 1: Test Initialization
func TestNewBloomFilter(t *testing.T) {
	m, k := uint(64), uint(3)
	bf := New(m, k)

	expectedBytes := 8
	if len(bf.bitset) != expectedBytes {
		t.Errorf("Expected %d bytes for %d bits, got %d", expectedBytes, m, len(bf.bitset))
	}
}

// Stage 3 & 4: Test Basic Add and Contains
func TestAddAndContains(t *testing.T) {
	bf := New(1024, 3)
	word := []byte("gopher")

	if bf.Contains(word) {
		t.Errorf("Empty filter reported containing word")
	}

	bf.Add(word)

	if !bf.Contains(word) {
		t.Errorf("Filter does not contain word after adding it")
	}
}

// Stage 4: Test for non-existent items
func TestContainsFalseNegatives(t *testing.T) {
	bf := New(1024, 3)
	bf.Add([]byte("apple"))
	bf.Add([]byte("orange"))

	if bf.Contains([]byte("banana")) {
		// Note: With a large m, this should be false.
		// If it's true, it's a rare false positive.
		t.Log("Note: Probabilistic false positive occurred (normal behavior at scale)")
	}
}

// Stage 5: Test Parameter Estimation
func TestEstimateParameters(t *testing.T) {
	n := uint(1000) // expected elements
	p := 0.01       // 1% false positive rate

	m, k := EstimateParameters(n, p)

	// Expected values for n=1000, p=0.01: m ~ 9585, k ~ 7
	if m == 0 || k == 0 {
		t.Errorf("Estimated parameters should not be zero")
	}
	if k < 1 {
		t.Errorf("K should be at least 1")
	}
}

// Stage 6: Test Concurrency (Race Condition Check)
func TestConcurrency(t *testing.T) {
	bf := New(1000, 3)
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(val int) {
			bf.Add([]byte(string(rune(val))))
			bf.Contains([]byte(string(rune(val))))
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
