package main

import (
  "time"
  "log"
)

func timeTrack(start time.Time, name string) {
  elapsed := time.Since(start)
  log.Printf("%s took %s", name, elapsed)
}

func round(f float64) float64 {
  return math.Floor(f + .5)
}
