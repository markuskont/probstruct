package probstruct

import (
	"errors"
	"fmt"
	"math"
	"math/bits"

	"github.com/markuskont/probstruct/pkg/hasher"
	"github.com/markuskont/probstruct/pkg/util"
)

const (
	bitness = 64
)

// HyperLogLog implements hyperloglog prob counting algorithm
type HyperLogLog struct {
	m uint32
	p uint
	// each bucket will hold max( count_zeroes + 1 ) in 64bit uint with 4..16 bits already derived
	// thus, uint8 cannot overflow
	buckets []uint8
	alpha   float64

	cardinality uint64

	// select hashing method, in line with bloom and cms implementations
	hash hasher.Algorithm
}

// NewHyperLogLog initializes new HLL struct
// hashing algorithm will default to fnv if left as nil
func NewHyperLogLog(precision uint, hash *hasher.Algorithm) (h *HyperLogLog, err error) {
	if precision < 4 || precision > 16 {
		return nil, errors.New("precision must be integer between 4 and 16")
	}
	h = &HyperLogLog{
		p:    precision,
		m:    1 << precision,
		hash: hasher.Fnv,
	}
	if hash != nil {
		h.hash = *hash
	}
	h.buckets = make([]uint8, h.m)
	// Magic numbers for hash collision correction
	switch h.m {
	case 16:
		h.alpha = 0.673
	case 32:
		h.alpha = 0.697
	case 64:
		h.alpha = 0.709
	default:
		h.alpha = 0.7213 / (1 + 1.079/float64(h.m))
	}
	return h, err
}

// Add is a wrapper for adding new items to leading zero counter
func (h *HyperLogLog) Add(items ...[]byte) *HyperLogLog {
	if items == nil || len(items) == 0 {
		return h
	}
	for _, v := range items {
		h.add64(h.hash.GetBaseHash(v).First())
	}
	return h
}

// Add64 calculates leading zeros from 64bit hash value and updates respective buckets
func (h *HyperLogLog) add64(hash uint64) *HyperLogLog {
	diff := bitness - h.p
	index := hash >> diff
	tail := hash << h.p
	count := uint8(bits.LeadingZeros64(tail)) + 1

	if count > h.buckets[index] {
		h.buckets[index] = count
	}
	return h
}

// Count calculates harmonic mean and updates cardinality
func (h *HyperLogLog) Count() *HyperLogLog {
	Z := float64(0)
	for _, c := range h.buckets {
		if c > 0 {
			Z += float64(1 / math.Pow(float64(2), float64(c)))
		}
	}
	Z = 1 / Z
	count := h.alpha * math.Pow(float64(h.m), 2) * Z
	h.cardinality = uint64(math.Floor(count))
	return h
}

// Merge allows multiple hll objects to be merged into a single structure
func Merge(containers ...*HyperLogLog) (*HyperLogLog, error) {
	if containers == nil || len(containers) < 2 {
		return nil, &ErrNoContainers{msg: "At least 2 HLL containers needed for merge"}
	}
	precision := containers[0].p
	hash := containers[0].hash
	for _, v := range containers {
		if v.p != precision {
			return nil, &ErrPrecisionMismatch{p1: v.p, p2: precision}
		}
		if v.hash != hash {
			return nil, &ErrHashFnMismatch{h1: v.hash, h2: hash}
		}
		precision = v.p
		hash = v.hash
	}
	merged, err := NewHyperLogLog(precision, &hash)
	if err != nil {
		return nil, err
	}
	for _, v := range containers {
		for j, v2 := range v.buckets {
			merged.buckets[j] = util.MaxUint8(v2, merged.buckets[j])
		}
	}
	return merged, nil
}

type ErrNoContainers struct{ msg string }

func (e ErrNoContainers) Error() string { return e.msg }

type ErrPrecisionMismatch struct{ p1, p2 uint }

func (e ErrPrecisionMismatch) Error() string {
	return fmt.Sprintf("HLL precision mismatch %d - %d, cannot merge", e.p1, e.p2)
}

type ErrHashFnMismatch struct{ h1, h2 hasher.Algorithm }

func (e ErrHashFnMismatch) Error() string {
	return fmt.Sprintf("HLL hash function mismatch %d - %d, cannot merge", e.h1, e.h2)
}
