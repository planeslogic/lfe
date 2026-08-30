package lfe

/*
#include <stdint.h>
*/
import "C"

import "fmt"

type Error struct {
	Op   string
	Code int
}

func (e *Error) Error() string {
	return fmt.Sprintf("lfe: %s failed (%d)", e.Op, e.Code)
}

func nativeError(op string, code C.int32_t) error {
	if code == 0 {
		return nil
	}
	return &Error{Op: op, Code: int(code)}
}

func nativeResolveError(op string, code C.intptr_t) error {
	return &Error{Op: op, Code: int(code)}
}
