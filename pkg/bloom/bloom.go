package bloom

import (
	"math"

	"github.com/markuskont/probstruct/pkg/hasher"
	"github.com/markuskont/probstruct/pkg/util"
)

// Filter is bitvector of length m elements
// items will be hashed to integers with k non-cryptographic functions
// boolean values in corresponding positions will be flipped
type Filter struct {
	m    uint64
	k    uint64
	bits []bool

	hsh hasher.Algorithm
}

// InitFilterWithEstimate instantiates a new bloom filter with user defined estimate parameters
// hash = hashing method to use ( <= 1 for murmur, 2 for fnv)
// n = number of elements in data set
// p = acceptable false positive 0 < p < 1 (no checks atm)
func NewFilterWithEstimate(n uint, p float64, h hasher.Algorithm) (b *Filter, err error) {
	m, k := estimateBloomSize(n, p)
	if b, err = NewFilter(m, k, h); err != nil {
		return nil, err
	}
	return b, nil
}

// InitFilter instantiates a new bloom filter with static length and hash function number
func NewFilter(m, k uint64, h hasher.Algorithm) (b *Filter, err error) {
	b = &Filter{
		m:   m,
		k:   k,
		hsh: h,
	}
	b.bits = make([]bool, b.m)
	return b, nil
}

// m = estimated size of bloom filter array
// m = -1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)
// k = num of needed hash functions
func estimateBloomSize(n uint, p float64) (m, k uint64) {
	size := math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(math.Log(2.0), 2.0))
	k = uint64(util.RoundFloat64(math.Log(2.0) * size / float64(n)))
	m = uint64(size)
	return
}

// Add method adds new element to bloom filter
func (b *Filter) Add(data ...[]byte) *Filter {
	if data == nil || len(data) == 0 {
		return b
	}
	locations := b.hsh.GetBaseHash(data...).Transform(b.m, b.k)
	if len(locations) == 0 {
		return b
	}
	for _, elem := range locations {
		b.bits[elem] = true
	}
	return b
}

// AddString is a helper to require need for type casting
func (b *Filter) AddString(data string) *Filter {
	return b.Add([]byte(data))
}

// AddStringSlice is a helper like AddString but assumes arbitrary number of strings
func (b *Filter) AddStringSlice(data ...string) *Filter {
	if data == nil || len(data) == 0 {
		return b
	}
	casted := make([][]byte, len(data))
	for i, v := range data {
		casted[i] = []byte(v)
	}
	return b.Add(casted...)
}

// Query returns the presence boolean of item from filter
// one missing bit is enough to verify non-existence
func (b Filter) Query(data []byte) bool {
	if data == nil || len(data) == 0 {
		return false
	}
	locations := b.hsh.GetBaseHash(data).Transform(b.m, b.k)
	for _, elem := range locations {
		if b.bits[elem] == false {
			return false
		}
	}
	return true
}

// QueryString is a helper to avoid excessive typecasting
func (b Filter) QueryString(data string) bool {
	return b.Query([]byte(data))
}

// QuerySlice is a wrapper for getting existence of multiple items with one go
func (b Filter) QuerySlice(data ...[]byte) []bool {
	if data == nil || len(data) == 0 {
		return []bool{}
	}
	res := make([]bool, len(data))
	for i, v := range data {
		res[i] = b.Query(v)
	}
	return res
}

// VecLen is a helper to get length of bloom, mostly for testing
func (b Filter) VecLen() uint64 {
	return b.m
}

// HashNo is a helper to get number of hash functions
func (b Filter) HashNo() uint64 {
	return b.k
}
