package probstruct

import (
	"math"
)

// BloomFilter is bitvector of length m elements
// items will be hashed to integers with k non-cryptographic functions
// boolean values in corresponding positions will be flipped
type BloomFilter struct {
	m    uint64
	k    uint64
	bits []bool
	hash int
}

// InitBloomWithEstimate instantiates a new bloom filter with user defined estimate parameters
// n = number of elements in data set
// p = acceptable false positive 0 < p < 1 (no checks atm)
// hash = hashing method to use ( <= 1 for murmur, 2 for fnv, else mix of both)
func InitBloomWithEstimate(n uint, p float64, h int) (b *BloomFilter, err error) {
	m, k := estimateBloomSize(n, p)
	b = &BloomFilter{
		m:    m,
		k:    k,
		hash: h,
	}
	b.bits = make([]bool, m)
	return b, err
}

// m = estimated size of bloom filter array
// m = -1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)
// k = num of needed hash functions
func estimateBloomSize(n uint, p float64) (m, k uint64) {
	size := math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2.0), 2.0))
	k = uint64(round(math.Log(2.0) * size / float64(n)))
	m = uint64(size)
	return
}

// integer values
func (b *BloomFilter) genLocs(data []byte) (locations []uint64) {
	locations = make([]uint64, b.k)
	h := genHashBase(data, b.hash)
	for i := uint64(0); i < b.k; i++ {
		locations[i] = transformHashes(h[0], h[1], i, b.m)
	}
	return
}

// Add method adds new element to bloom filter
func (b *BloomFilter) Add(data []byte) *BloomFilter {
	for _, elem := range b.genLocs(data) {
		b.bits[elem] = true
	}
	return b
}

// AddString converts item to []byte and returns Add()
func (b *BloomFilter) AddString(data string) *BloomFilter {
	return b.Add([]byte(data))
}

// Query returns the presence boolean of item from filter
// one missing bit is enough to verify non-existence
func (b *BloomFilter) Query(data []byte) bool {
	for _, elem := range b.genLocs(data) {
		if b.bits[elem] == false {
			return false
		}
	}
	return true
}

// QueryString converts item to []byte and returns Query()
func (b *BloomFilter) QueryString(data string) bool {
	return b.Query([]byte(data))
}

// AssessFill is helper to return instance load ratio
func (b *BloomFilter) AssessFill() float64 {
	total := b.m
	used := 0
	for _, loc := range b.bits {
		if loc == true {
			used++
		}
	}
	return float64(used) / float64(total)
}
