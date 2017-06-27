package probstruct

import (
  "fmt"
  "log"
)

func main() {
  bloom, err := NewBloomWithEstimate(300000, 0.01, 1)
  countminsketch, err := InitMinSketchWithEstimate(0.01, 0.01, 1)
  //bloom, err := NewBloomWithEstimate(300000000, 0.01, -1)
  fmt.Println("Initiated bloom filter with m:", bloom.m, "k:", bloom.k, "hash method:", bloom.hash)
  if err != nil { log.Fatal(err) }
  fmt.Println(len(bloom.bits))
  //fmt.Println(genBaseHashes([]byte("test1m3")))
  bloom.AddString("test1m3")
  bloom.Add([]byte("test1m2"))
  bloom.Add([]byte("lalala"))
  //for bit := range bloom.bits {
  //  if bloom.bits[bit] == true {
  //    fmt.Println(bit)
  //  }
  //}
  fmt.Println( bloom.QueryString("test1m3") )
  fmt.Println( bloom.QueryString("test3m3") )
  fmt.Println( "----------------------" )
  fmt.Println(countminsketch.width)
  fmt.Println(countminsketch.depth)
  countminsketch.AddString("test1m3")
  fmt.Println( countminsketch.QueryString("test1m3") )
  //asd := genRandomIntegers(10)
  //fmt.Println(asd)
}
