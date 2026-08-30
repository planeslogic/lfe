//go:build !lfe_test_native

package lfe

/*
#cgo LDFLAGS: -L${SRCDIR}/native/target/release -llfe_be_sdk_go_native -Wl,-rpath,${SRCDIR}/native/target/release
*/
import "C"
