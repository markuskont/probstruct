package probstruct

// WIP

import (
  "math"
  "fmt"
)

type LinearCounting struct {
  m     uint64
  bits  []bool
  hash  int
}

func InitLinearCounting(m uint64, h int) (lc *LinearCounting, err error) {
  lc = &LinearCounting{
    m:  m,
  }
  lc.bits = make([]bool, m)
  return lc, err
}

func (lc *LinearCounting) Add(data []byte) (location uint64) {
  h := genHashBase(data, lc.hash)
  location = transformHashes(h[0], h[1], 1, lc.m)
  lc.bits[location] = true
  return
}

func (lc *LinearCounting) AddString(data string) uint64 {
  return lc.Add([]byte(data))
}

func (lc *LinearCounting) GetFill(val bool) (m uint64) {
  m = 0
  for _, bit := range lc.bits {
    if bit == val { m += 1 }
  }
  return
}

func (lc *LinearCounting) AssessCardinality() float64 {
  Un := lc.GetFill(false)
  Vn := float64( Un ) / float64( lc.m )
  n := math.Log( Vn )
  return -1 * float64(lc.m) * n
}

func (lc *LinearCounting) ReturnData() []bool {
  return lc.bits
}

func (lc *LinearCounting) ReturnSize() (m uint64) {
  return lc.m
}
