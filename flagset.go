package lfe

import "fmt"

// FlagSetValue is the SDK representation of one packed FlagSet logical value.
// Applications identify flags by position; the SDK owns shifting and packing.
type FlagSetValue struct {
	lo      uint64
	hi      uint64
	invalid bool
}

// FlagSet converts one raw 0/1 application value at position into a FlagSet value.
// Positions are zero-based and limited to 0..127.
func FlagSet(position uint8, raw uint8) FlagSetValue {
	if position > 127 || raw > 1 {
		return FlagSetValue{invalid: true}
	}
	if raw == 0 {
		return FlagSetValue{}
	}
	if position < 64 {
		return FlagSetValue{lo: uint64(1) << position}
	}
	return FlagSetValue{hi: uint64(1) << (position - 64)}
}

// Flags composes multiple FlagSet values into one value without exposing packing.
func Flags(values ...FlagSetValue) FlagSetValue {
	var out FlagSetValue
	for _, value := range values {
		out.lo |= value.lo
		out.hi |= value.hi
		out.invalid = out.invalid || value.invalid
	}
	return out
}

func (v FlagSetValue) validate() error {
	if v.invalid {
		return fmt.Errorf("lfe: FlagSet position must be 0..127 and raw value must be 0 or 1")
	}
	return nil
}

func (v FlagSetValue) cmBatchValue() batchValue {
	var err error
	if v.invalid {
		err = fmt.Errorf("lfe: FlagSet position must be 0..127 and raw value must be 0 or 1")
	}
	return batchValue{kind: 2, hi: v.hi, lo: v.lo, invalid: err}
}
