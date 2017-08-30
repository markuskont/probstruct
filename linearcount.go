package probstruct

// WIP

//import (
//  "fmt"
//)

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

func (lc *LinearCounting) ReturnData() []bool {
  return lc.bits
}

func (lc *LinearCounting) ReturnSize() (m uint64) {
  return lc.m
}
