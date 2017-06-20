package main

import (
  "math"
  "hash/fnv"

  "github.com/spaolacci/murmur3"
)

// generate base hash values which will be later used for uniques hashes
func genHashBase(data []byte, b *BloomFilter) (h [2]uint64) {
  switch {
  case b.hash <= 1:
    h = genHashBaseMurmur(data)
  case b.hash == 2:
    h = genHashBaseFnv(data)
  case b.hash >= 3:
    h = genHashBaseCombo(data)
  }
  return
}

// this function uses murmur3 hash
// 128bit integer split into 2 distinct 64bit sections
func genHashBaseMurmur(data []byte) [2]uint64 {
  hasher := murmur3.New128()
  hasher.Write(data)
  h1, h2 := hasher.Sum128()
  return [2]uint64{
    h1, h2,
  }
}
// simple fnv
func genHashBaseFnv(data []byte) [2]uint64 {
  hasher := fnv.New64a()
  hasher.Write(data)
  h1 := hasher.Sum64()
  hasher.Write([]byte{1})
  h2 := hasher.Sum64()
  return [2]uint64{
    h1, h2,
  }
}
// mix of fnv and 64bit murmur3
func genHashBaseCombo(data []byte) [2]uint64 {
  hasher1 := fnv.New64a()
  hasher1.Write(data)
  h1 := hasher1.Sum64()
  hasher2 := murmur3.New64()
  hasher2.Write(data)
  h2 := hasher2.Sum64()
  return [2]uint64{
    h1, h2,
  }
}

// https://www.eecs.harvard.edu/~michaelm/postscripts/rsa2008.pdf
func transformHashes(h1, h2, i, size uint64) uint64 {
  return  ( h1 + i * h2 + uint64(math.Pow(float64(i), 2)) ) % size
}
