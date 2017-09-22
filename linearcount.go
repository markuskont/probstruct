package probstruct

// WIP

import (
	"math"
	//"fmt"
)

// LinearCounting is used to measure dataset cardinality
type LinearCounting struct {
	m    uint64
	bits []bool
	hash int
}

// InitLinearCounting instantiates new object
func InitLinearCounting(m uint64, h int) (lc *LinearCounting, err error) {
	lc = &LinearCounting{
		m: m,
	}
	lc.bits = make([]bool, m)
	return lc, err
}

// Add adds new element
func (lc *LinearCounting) Add(data []byte) (location uint64) {
	h := genHashBase(data, lc.hash)
	location = transformHashes(h[0], h[1], 1, lc.m)
	lc.bits[location] = true
	return
}

// AddString converts textual input and returns Add()
func (lc *LinearCounting) AddString(data string) uint64 {
	return lc.Add([]byte(data))
}

// GetFill is used to estimate load
// True and False values can be estimated as required
func (lc *LinearCounting) GetFill(val bool) (m uint64) {
	m = 0
	for _, bit := range lc.bits {
		if bit == val {
			m++
		}
	}
	return
}

// AssessCardinality is used to measure dataset cardinality depending on load factor
func (lc *LinearCounting) AssessCardinality() float64 {
	Un := lc.GetFill(false)
	Vn := float64(Un) / float64(lc.m)
	n := math.Log(Vn)
	return -1 * float64(lc.m) * n
}

// ReturnData is helper that returns raw bitvector
func (lc *LinearCounting) ReturnData() []bool {
	return lc.bits
}

// ReturnSize is helper to return the length of bitvector
func (lc *LinearCounting) ReturnSize() (m uint64) {
	return lc.m
}
