package hasher

import (
	"hash/fnv"

	"github.com/spaolacci/murmur3"
)

type BaseHash [2]uint64

func (b BaseHash) Transform(bound, N uint64) []uint64 {
	var locations = make([]uint64, N)
	var i uint64
	for i = 0; i < N; i++ {
		locations[i] = (b.First() + i*b.Second()) % bound
	}
	return locations
}

func (b BaseHash) First() uint64  { return b[0] }
func (b BaseHash) Second() uint64 { return b[1] }

type Algorithm int

const (
	Fnv Algorithm = iota
	Murmur
)

func (a Algorithm) GetBaseHash(items ...[]byte) BaseHash {
	if items == nil || len(items) == 0 {
		return BaseHash{0, 0}
	}
	var h1, h2 uint64
	switch a {
	case Fnv:
		h := fnv.New64()
		for _, item := range items {
			h.Write(item)
		}
		h1 = h.Sum64()
		h.Write([]byte{1})
		h2 = h.Sum64()
	case Murmur:
		h := murmur3.New128()
		for _, item := range items {
			h.Write(item)
		}
		h1, h2 = h.Sum128()
	}
	return BaseHash{h1, h2}
}
