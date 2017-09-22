package probstruct

import (
	"errors"
	"math"
	"sync/atomic"
)

// CountMinSketch is a counting BloomFilter
// depth corresponds to number of distinct hashing funcitons used
// with corresponds to number of counters in each hash set
// https://blog.demofox.org/2015/02/22/count-min-sketch-a-probabilistic-histogram/
type CountMinSketch struct {
	depth uint64
	width uint64
	count [][]uint64
	hash  int
}

// InitMinSketchWithEstimate instantiates a new CMS object with user defined estimate parameters
// width = [ e / epsilon ]
// depth = [ ln( 1 / delta ) ]
// hash = hashing method to use ( <= 1 for murmur, 2 for fnv, else mix of both)
func InitMinSketchWithEstimate(epsilon, delta float64, h int) (s *CountMinSketch, err error) {
	depth, width := estimateCountMinSize(epsilon, delta)
	if epsilon <= 0 || epsilon >= 1 {
		return nil, errors.New("CountMinSketch.Init: epsilon must be 0 < eps < 1")
	}
	if delta <= 0 || delta >= 1 {
		return nil, errors.New("CountMinSketch.Init: delta must be 0 < eps < 1")
	}
	s = &CountMinSketch{
		depth: depth,
		width: width,
		hash:  h,
	}
	s.count = make([][]uint64, depth)
	for i := uint64(0); i < depth; i++ {
		s.count[i] = make([]uint64, width)
	}
	return s, err
}

func estimateCountMinSize(epsilon, delta float64) (depth, width uint64) {
	depth = uint64(math.Ceil(math.Log(1.0 / delta)))
	width = uint64(math.Ceil(math.E / epsilon))
	return
}

func (s *CountMinSketch) genLocs(data []byte) (locations []uint64) {
	locations = make([]uint64, s.depth)
	h := genHashBase(data, s.hash)
	for i := uint64(0); i < uint64(s.depth); i++ {
		locations[i] = transformHashes(h[0], h[1], i, uint64(s.width))
	}
	return
}

// Increment item count in CMS without returning the new estimated value
func (s *CountMinSketch) Increment(data []byte) *CountMinSketch {
	// location = hashing function i < depth
	for i, elem := range s.genLocs(data) {
		atomic.AddUint64(&s.count[i][elem], 1)
		//s.count[i][elem] += 1
	}
	return s
}

// IncrementGetVal is a combination of Increment() and QueryMin() methods that returns new estimation upon adding each element
// deduplicates needed work if estimation has to be compared to threshold
func (s *CountMinSketch) IncrementGetVal(data []byte) (min uint64) {
	// location = hashing function i < depth
	for i, elem := range s.genLocs(data) {
		c := &s.count[i][elem]
		atomic.AddUint64(c, 1)
		if min == 0 || *c < min {
			min = *c
		}
	}
	return
}

// IncrementStringGetVal converts textual input before returning IncrementGetVal()
func (s *CountMinSketch) IncrementStringGetVal(data string) (min uint64) {
	return s.IncrementGetVal([]byte(data))
}

// IncrementString converts textual input before returning Increment()
func (s *CountMinSketch) IncrementString(data string) *CountMinSketch {
	return s.Increment([]byte(data))
}

// QueryMin returns estimated value for item
// smallest count = least collisions, thus most accurate estimation
// if smallest value is zero, the item has not been counted before.
// CMS cannot under-estimate by definition, thus any subsequent checks are waste of CPU cycles
func (s *CountMinSketch) QueryMin(data []byte) (min uint64) {
	for i, elem := range s.genLocs(data) {
		c := s.count[i][elem]
		if c == 1 {
			min = 1
			break
		} else if min == 0 || c < min {
			min = c
		}
	}
	return
}

// QueryString converts textual input before returning Query()
func (s *CountMinSketch) QueryString(data string) uint64 {
	return s.QueryMin([]byte(data))
}

// ReturnCounts is helper to access raw matrix in instance
func (s *CountMinSketch) ReturnCounts() [][]uint64 {
	return s.count
}

// GetDimensions is helper to access matrix dimensions
func (s *CountMinSketch) GetDimensions() (w, d uint64) {
	return s.width, s.depth
}

// AssessFill is helper to return instance load ratio
func (s *CountMinSketch) AssessFill() float64 {
	total := s.width * s.depth
	used := 0
	for _, block := range s.count {
		for _, val := range block {
			if val > 0 {
				used++
			}
		}
	}
	return float64(used) / float64(total)
}
