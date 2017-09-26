package probstruct

import(
  "errors"
  "math"
  "math/bits"
)

// N = number of buckets in registry
// b = number of significant bits used assigning incoming data
// leading zeros in remaining binary value are used to estimate probability
// buckets = estimators with Length of N
type HLL struct {
  m uint32
  p uint
  // each bucket will hold max( count_zeroes + 1 ) in 64bit uint with 4..16 bits already derived
  // thus, uint8 cannot overflow
  buckets []uint8
  alpha float64

  cardinality uint64

  // some data points are already globally unique unsigned integers, e.g. IPv4 addresses and ports
  // thus, alpha correction for hash collisions is not needed, and we will have to adjust for lenth when calculating buckets
  hashing bool
  // switch between 32 and 64 bit input items
  // by default bitwise operations assumes unsigned 64bit hash input
  // breaks when counting 32bit globally unique values
  bitness uint
  // select hashing method, in line with bloom and cms implementations
  hash int
}

// Init new HLL instance
// user should choose between 4 and 16 significant bits
// each bit increases computational complexity by order of magnitude
// 16 = 65536, 15 = 32768, 14 = 16384, etc
// x << y == x * 2**y
func InitHLL(precision uint, hashing bool, hash int) (h *HLL, err error) {
  if precision < 4 || precision > 16 {
    return nil, errors.New("precision must be integer between 4 and 16")
  }
  h = &HLL{
    p:  precision,
    m:  1 << precision,
    hashing:  hashing,
    hash: hash,
  }
  h.buckets = make([]uint8, h.m)
  if h.hashing == false {
    h.bitness = 32
    h.alpha = 1
  } else {
    h.bitness = 64
    // Magic numbers for hash collision correction
    switch h.m {
    case 16:
      h.alpha = 0.673
    case 32:
      h.alpha = 0.697
    case 64:
      h.alpha = 0.709
    default:
      h.alpha = 0.7213 / ( 1 + 1.079 / float64(h.m) )
    }
  }
  return h, err
}

func (h *HLL) AddString(item string) *HLL {
  return h.AddHash([]byte(item))
}

// Add converts databyte item into uint64 and adds position of first true boolean after bitwise header into appropriate bucket
func (h *HLL) AddHash(item []byte) *HLL {
  hash := genHashBase(item, h.hash)[0]
  return h.Add(hash)
}

func (h *HLL) Add(hash uint64) *HLL {
  diff := h.bitness - h.p
  index := hash >> diff
  tail := hash << h.p
  count := uint8(bits.LeadingZeros64(tail)) + 1

  if count > h.buckets[index] {
    h.buckets[index] = count
  }
  return h
}

func (h *HLL) Count() *HLL {
  Z := float64(0)
  for _, c := range h.buckets {
    if c > 0 {
      Z += float64( 1 / math.Pow( float64(2), float64(c) ) )
    }
  }
  Z = 1 / Z
  count := h.alpha * math.Pow( float64(h.m), 2) * Z
  h.cardinality = uint64( math.Floor(count) )
  return h
}

func (h *HLL) GetCardinality() uint64 {
  return h.cardinality
}

func (h *HLL) GetCounters() uint32 {
  return h.m
}
