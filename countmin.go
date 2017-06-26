package main

import (
  "errors"
  "math"
)

// estimate size of count-min-sketch using provided estimates
// width = [ e / epsilon ]
// depth = [ ln( 1 / delta ) ]
// hash = hashing method to use ( <= 1 for murmur, 2 for fnv, else mix of both)
type CountMinSketch struct {
  depth   uint
  width   uint
  count   [][]uint64
  hash    int
}

func InitMinSketchWithEstimate(epsilon, delta float64, h int) (s *CountMinSketch, err error) {
  depth, width := estimateCountMinSize(epsilon, delta)
  if epsilon <= 0 || epsilon >= 1  { return nil, errors.New("CountMinSketch.Init: epsilon must be 0 < eps < 1") }
  if delta <= 0 || delta >= 1  { return nil, errors.New("CountMinSketch.Init: delta must be 0 < eps < 1") }
  s = &CountMinSketch{
    depth:  depth,
    width:  width,
    hash:   h,
  }
  s.count = make([][]uint64, depth)
  for i := uint(0); i < depth; i++ { s.count[i] = make([]uint64, width) }
  return s, err
}

func estimateCountMinSize(epsilon, delta float64) (depth, width uint) {
  depth = uint( math.Ceil( math.Log( 1.0 / delta ) ) )
  width = uint( math.Ceil( math.E / epsilon ) )
  return
}

//func (s *CountMinSketch) Add(data []byte) *CountMinSketch {
//  //defer timeTrack(time.Now(), "bloom add")
//  locations := 1
//  //for i := range locations {
//  //  //if b.bits[locations[i]] == true { fmt.Println("Collision!") }
//  //  b.bits[locations[i]] = true
//  //}
//  return s
//}
//
//func (s *CountMinSketch) genLocs(data []byte) []uint64 {
//
//}
