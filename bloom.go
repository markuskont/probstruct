package main

import (
  //"fmt"
  "log"
  "math"
  //"time"
)

// n = number of elements in data set
// p = acceptable false positive 0 < p < 1 (no checks atm)
// m = estimated size of bloom filter array
// m = -1 * float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)
// k = num of needed hash functions
// hash = hashing method to use ( <= 1 for murmur, 2 for fnv, else mix of both)
type BloomFilter struct {
  m     uint
  k     uint
  bits  []bool
  hash  int
}

func NewBloomWithEstimate(n uint, p float64, h int) (b *BloomFilter, err error) {
  m, k := estimateBloomSize(n, p)
  b = &BloomFilter{
    m:    m,
    k:    k,
    hash: h,
  }
  b.bits = make([]bool, m)
  return b, err
}

func estimateBloomSize(n uint, p float64) (m, k uint) {
  //m = math.Ceil(( float64(n) * math.Log(p) ) / math.Log(1.0 / math.Pow( 2.0, math.Log(2.0) )))
  size := math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow( math.Log(2.0), 2.0 ))
  k = uint( round(math.Log(2.0) * size / float64(n)) )
  m = uint( size )
  // max size for 32bit integer
  if m > 4294967295 { log.Fatal("Estimated bitarray length ", m, " does not fit in unsigned 32bit integer. Dataset size is ", n, " and confidence is ", p, ". Try lowering.") }
  return
}

func genLocs(data []byte, b *BloomFilter) (locations []uint64) {
  locations = make([]uint64, b.k)
  h := genHashBase(data, b)
  for i := uint64(0); i < uint64(b.k); i++ {
    locations[i] = transformHashes(h[0], h[1], i, uint64(b.m))
  }
  return
}

func (b *BloomFilter) Add(data []byte) *BloomFilter {
  //defer timeTrack(time.Now(), "bloom add")
  locations := genLocs(data, b)
  for i := range locations {
    //if b.bits[locations[i]] == true { fmt.Println("Collision!") }
    b.bits[locations[i]] = true
  }
  return b
}

func (b *BloomFilter) AddString(data string) *BloomFilter {
  return b.Add([]byte(data))
}

func (b *BloomFilter) Query(data []byte) bool {
  //defer timeTrack(time.Now(), "bloom query")
  locations := genLocs(data, b)
  for i := range locations {
    if b.bits[locations[i]] == false {
      // one missing bit is enough to verify non-existence, exit ASAP
      return false
    }
  }
  return true
}

func (b *BloomFilter) QueryString(data string) bool {
  return b.Query([]byte(data))
}
