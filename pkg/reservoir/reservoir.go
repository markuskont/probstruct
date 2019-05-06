package reservoir

import (
	"math/rand"
)

// Sample implements simple reservoir filter
type Sample struct {
	k        int
	total    uint64
	switches uint64
	sample   []interface{}
}

// NewSample instantiates new Sample struct
func NewSample(k int) (r *Sample, err error) {
	r = &Sample{
		k:        k,
		total:    0,
		switches: 0,
		sample:   make([]interface{}, k),
	}
	return r, nil
}

// Add new item to reservoir
func (r *Sample) Add(item interface{}) *Sample {
	r.total++
	if len(r.sample) < r.k {
		r.sample = append(r.sample, item)
		return r
	}
	if rand.Float64() < (float64(r.k) / float64(r.total)) {
		r.sample[rand.Intn(r.k)] = item
		r.switches++
	}
	return r
}

// GetSample is a helper to return size of sampled data
func (r *Sample) GetSample() []interface{} { return r.sample }

// GetK is a helper to return all sampled values
func (r *Sample) GetK() int { return r.k }

// GetTotal is a helper to return number of items seen
func (r *Sample) GetTotal() uint64 { return r.total }

// GetSwitches is a helper to return number of items seen
func (r *Sample) GetSwitches() uint64 { return r.switches }
