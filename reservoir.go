package probstruct

import (
	"math/rand"
)

// Reservoir implements simple reservoir filter
type Reservoir struct {
	k      int
	total  uint64
	sample []interface{}
}

// InitReservoir instantiates new Reservoir struct
func InitReservoir(k int) (r *Reservoir, err error) {
	r = &Reservoir{
		k:      k,
		total:  0,
		sample: make([]interface{}, k),
	}
	return r, nil
}

// Add new item to reservoir
func (r *Reservoir) Add(item interface{}) *Reservoir {
	r.total++
	if len(r.sample) < r.k {
		r.sample = append(r.sample, item)
	} else {
		if rand.Float64() < (float64(r.k) / float64(r.total)) {
			r.sample[rand.Intn(r.k)] = item
		}
	}
	return r
}

// GetSample is a helper to return size of sampled data
func (r *Reservoir) GetSample() []interface{} {
	return r.sample
}

// GetK is a helper to return all sampled values
func (r *Reservoir) GetK() int {
	return r.k
}
