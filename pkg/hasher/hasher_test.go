package hasher

import (
	"bytes"
	"flag"
	"testing"
)

var err error
var DEBUG bool

func init() {
	flag.BoolVar(&DEBUG, "debug", false, "Enable additional information logging")
	flag.Parse()
}

var data = []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas vulputate erat non rhoncus pretium. Nullam pulvinar auctor neque at aliquet. Donec aliquet ante consectetur blandit lacinia. Quisque in risus nunc. Aliquam pretium molestie faucibus. Morbi nec nibh sit amet mi mollis rutrum. Aenean fringilla varius egestas. Ut consectetur neque ac ipsum vulputate tempor.

Fusce augue augue, aliquet et velit eu, posuere aliquam neque. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Integer et convallis lorem. Proin ut volutpat eros. Donec vestibulum tincidunt nisi, et facilisis mauris dignissim id. Mauris augue enim, pretium ac orci sed, tincidunt ullamcorper odio. Vivamus at nisl metus.

Vivamus id ultricies nisi. Pellentesque lorem nulla, vestibulum sed elit at, sollicitudin cursus eros. Aliquam vitae lectus erat. Morbi hendrerit mollis venenatis. Donec porttitor mi eget felis dignissim, in vulputate mi porta. Phasellus non scelerisque mi. Donec nisl magna, accumsan vitae lacus quis, sagittis ultrices risus.

Maecenas faucibus suscipit eros, eu placerat ipsum blandit nec. Nulla id leo quis mi elementum porttitor non quis magna. Sed vitae nunc a erat scelerisque aliquet eu in tortor. Phasellus odio ante, scelerisque in magna ac, auctor aliquam urna. Sed ut diam id risus consequat egestas a nec lorem. Vestibulum tempus libero a lectus rhoncus, quis venenatis enim fermentum. Vivamus bibendum nisi ac sodales cursus. Sed vitae dui imperdiet, aliquet mi eu, interdum lectus.

Praesent elit odio, vestibulum sed sagittis et, vehicula id neque. Etiam sit amet massa eget quam commodo placerat vitae ut nibh. Fusce urna sapien, efficitur ac vulputate sit amet, blandit vulputate mi. Aenean at mollis justo. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Ut ullamcorper, sem id viverra maximus, nisi diam rutrum leo, non viverra erat augue commodo nunc. Suspendisse tincidunt venenatis mauris quis venenatis. Cras sit amet tincidunt ipsum. Integer scelerisque, odio laoreet lobortis sagittis, nulla ipsum varius odio, in pretium tortor sem id sapien. Vestibulum ac magna cursus, aliquam urna et, sollicitudin sapien.`)
var input = bytes.Split(data, []byte(" "))

func BenchmarkFnvHashing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fnv.GetBaseHash(input[0])
	}
}
func BenchmarkMurmurHashing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Murmur.GetBaseHash(input[0])
	}
}
func BenchmarkFnvHashingLoop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fnv.GetBaseHash(input...)
	}
}
func BenchmarkMurmurHashingLoop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Murmur.GetBaseHash(input...)
	}
}

var max = uint64(300000)
var uniq = uint64(10)

func BenchmarkFnvTransform(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fnv.GetBaseHash(input[0]).Transform(max, uniq)
	}
}
func BenchmarkMurmurTransform(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Murmur.GetBaseHash(input[0]).Transform(max, uniq)
	}
}
func BenchmarkFnvLoopTransform(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Fnv.GetBaseHash(input...).Transform(max, uniq)
	}
}
func BenchmarkMurmurLoopTransform(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Murmur.GetBaseHash(input...).Transform(max, uniq)
	}
}
