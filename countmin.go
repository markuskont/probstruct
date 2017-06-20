package main

import (
  "errors"
)

// d is number of hash functions
// w is size of every hash table
// seeds contain unique values which combined with hashing function produce d distinct sketches for each input
type CountMinSketch struct {
  d       uint
  w       uint
  count   [][]uint64
}

func InitMinSketch(d uint, w uint) (s *CountMinSketch, err error) {
  if d < 1 || w < 1 { return nil, errors.New("CountMinSketch.Init: d and w must be >= 1") }
  s = &CountMinSketch{
    d:      d,
    w:      w,
  }
  s.count = make([][]uint64, d)
  for i := uint(0); i < d; i++ { s.count[i] = make([]uint64, w) }
  return s, nil
}
