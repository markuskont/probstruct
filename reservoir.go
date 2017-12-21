package probstruct

// Reservoir implements simple reservoir filter
type Reservoir struct {
	k      uint
	full   bool
	sample []interface{}
}

// InitReservoir instantiates new Reservoir struct
func InitReservoir(k uint) (r *Reservoir, err error) {
	r = &Reservoir{
		k:      k,
		full:   false,
		sample: make([]interface{}, k),
	}
	return r, nil
}
