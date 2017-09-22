package probstruct

// PLACEHOLDER

//import (
//  "fmt"
//)

// SpaceSaving is algorithm to measure ocurrences of top-k elements in data stream
type SpaceSaving struct {
	k        uint
	counters map[string]uint64
}

// InitSpaceSaving instantiates new SpaceSaving object
func InitSpaceSaving(k uint) (s *SpaceSaving, err error) {
	s = &SpaceSaving{
		k: k,
	}
	return s, err
}
