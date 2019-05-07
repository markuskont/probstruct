package countmin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"sync/atomic"

	"github.com/markuskont/probstruct/pkg/hasher"
)

// CountMinSketch is a counting BloomFilter
// depth corresponds to number of distinct hashing funcitons used
// with corresponds to number of counters in each hash set
// https://blog.demofox.org/2015/02/22/count-min-sketch-a-probabilistic-histogram/
type Sketch struct {
	depth uint64
	width uint64
	count [][]uint64
	hash  hasher.Algorithm
}

// NewSketch instantiates a new sketch object
func NewSketch(width, depth uint64, h hasher.Algorithm) (*Sketch, error) {
	s := &Sketch{
		depth: depth,
		width: width,
		hash:  h,
	}
	return s.InitCounters(), nil
}

// NewSketchWithEstimate instantiates a new CMS object with user defined estimate parameters
// width = [ e / epsilon ]
// depth = [ ln( 1 / delta ) ]
// hash = hashing method to use ( <= 1 for murmur, 2 for fnv, else mix of both)
func NewSketchWithEstimate(epsilon, delta float64, h hasher.Algorithm) (*Sketch, error) {
	if epsilon <= 0 || epsilon >= 1 {
		return nil, errors.New("CountMinSketch.Init: epsilon must be 0 < eps < 1")
	}
	if delta <= 0 || delta >= 1 {
		return nil, errors.New("CountMinSketch.Init: delta must be 0 < eps < 1")
	}
	w, d := estimateCountMinSize(epsilon, delta)
	return NewSketch(w, d, h)
}

// InitCounters is a small helper for initializing or resetting the unerlying matrix
func (s *Sketch) InitCounters() *Sketch {
	s.count = make([][]uint64, s.depth)
	for i := uint64(0); i < s.depth; i++ {
		s.count[i] = make([]uint64, s.width)
	}
	return s
}

// Increment item count in CMS without returning the new estimated value
// Uses atomics for efficient thread safety
func (s *Sketch) Increment(data []byte) uint64 {
	// location = hashing function i < depth
	var min uint64
	locations := s.hash.GetBaseHash(data).Transform(s.width, s.depth)
	for i, j := range locations {
		count := atomic.AddUint64(&s.count[i][j], 1)
		if count < min || count == 0 {
			min = count
		}
	}
	return min
}

// IncrementString is a wrapper to avoid excessive typecasting
func (s *Sketch) IncrementString(data string) uint64 {
	return s.Increment([]byte(data))
}

// IncrementAny is a wrapper to handle arbitrary data types
// Encooding errors are currently quietly consumed to maintain return type consistency
func (s *Sketch) IncrementAny(data interface{}) uint64 {
	switch v := data.(type) {
	case []byte:
		return s.Increment(v)
	case string:
		return s.IncrementString(v)
	default:
		b, err := encodeNonByte(data)
		if err == nil {
			return s.Increment(b)
		}
		return 0
	}
}

// Query is nearly identical fot Increment but does not modify anything, only returning the estimation
func (s Sketch) Query(data []byte) uint64 {
	// location = hashing function i < depth
	var min uint64
	locations := s.hash.GetBaseHash(data).Transform(s.width, s.depth)
	for i, j := range locations {
		count := s.count[i][j]
		if count < min || count == 0 {
			min = count
		}
	}
	return min
}

// QueryString is a helper to avoid excessive typecasting
func (s Sketch) QueryString(data string) uint64 {
	return s.Query([]byte(data))
}

// QueryAny is a wrapper to handle arbitrary data types
// Encooding errors are currently quietly consumed to maintain return type consistency
func (s Sketch) QueryAny(data interface{}) uint64 {
	switch v := data.(type) {
	case []byte:
		return s.Query(v)
	case string:
		return s.QueryString(v)
	default:
		b, err := encodeNonByte(data)
		if err == nil {
			return s.Query(b)
		}
		return 0
	}
}

func encodeNonByte(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func estimateCountMinSize(epsilon, delta float64) (depth, width uint64) {
	depth = uint64(math.Ceil(math.Log(1.0 / delta)))
	width = uint64(math.Ceil(math.E / epsilon))
	return
}
