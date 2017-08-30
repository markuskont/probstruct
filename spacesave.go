package probstruct

// WIP

//import (
//  "fmt"
//)

type SpaceSaving struct {
  k         uint
  counters  map[string]uint64
}

func InitSpaceSaving(k uint) (s *SpaceSaving, err error) {
  s = &SpaceSaving{
    k:  k,
  }
  return s, err
}
