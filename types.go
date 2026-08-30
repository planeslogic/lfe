package lfe

type Operator uint8

const (
	Eq Operator = iota
	Neq
	Gt
	Gte
	Lt
	Lte
)

type DateTime struct {
	Year   uint16
	Month  uint8
	Day    uint8
	Hour   uint8
	Minute uint8
	Second uint8
}

func signedMagnitude(value int64) (negative bool, magnitude uint64) {
	if value < 0 {
		return true, uint64(-(value + 1)) + 1
	}
	return false, uint64(value)
}
