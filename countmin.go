package main

import (
  //"fmt"
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

func (s *CountMinSketch) genLocs(data []byte) (locations []uint64) {
  locations = make([]uint64, s.depth)
  h := genHashBase(data, s.hash)
  for i := uint64(0); i < uint64(s.depth); i++ {
    locations[i] = transformHashes(h[0], h[1], i, uint64(s.width))
  }
  return
}

func (s *CountMinSketch) Add(data []byte) *CountMinSketch {
  locations := s.genLocs(data)
  // location = hashing function i < depth
  for i, elem := range locations {
    //val := &s.count[i][locations[i]]
    s.count[i][elem] += 1
  }
  return s
}

func (s *CountMinSketch) AddString(data string) *CountMinSketch {
  return s.Add([]byte(data))
}

func (s *CountMinSketch) QueryMin(data []byte) (min uint64) {
  locations := s.genLocs(data)
  for i, elem := range locations {
    c := s.count[i][elem]
    // 1 = only 0 can be smaller, but element is not in dataset in this case (bloom false negative logic)
    // no point in iterating further
    if c == 1 {
      min = 1
      break
    } else if min == 0 || c < min {
      min = c
    }
  }
  return
}

func (s *CountMinSketch) QueryString(data string) uint64 {
  return s.QueryMin([]byte(data))
}
