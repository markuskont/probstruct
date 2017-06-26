package main

import (
  "errors"
  "math"
)

// d is number of hash functions
// w is size of every hash table
// count is matrix for containing counts
type CountMinSketch struct {
  depth   uint
  width   uint
  count   [][]uint64
}

func InitMinSketchWithEstimate(epsilon, delta float64) (s *CountMinSketch, err error) {
  depth, width := estimateCountMinSize(epsilon, delta)
  if epsilon <= 0 || epsilon >= 1  { return nil, errors.New("CountMinSketch.Init: epsilon must be 0 < eps < 1") }
  if delta <= 0 || delta >= 1  { return nil, errors.New("CountMinSketch.Init: delta must be 0 < eps < 1") }
  s = &CountMinSketch{
    depth:  depth,
    width:  width,
  }
  s.count = make([][]uint64, depth)
  for i := uint(0); i < depth; i++ { s.count[i] = make([]uint64, width) }
  return s, err
}

// estimate size of count-min-sketch using provided estimates
// width = [ e / epsilon ]
// depth = [ ln( 1 / delta ) ]
func estimateCountMinSize(epsilon, delta float64) (depth, width uint) {
  depth = uint( math.Ceil( math.Log( 1.0 / delta ) ) )
  width = uint( math.Ceil( math.E / epsilon ) )
  return
}
