//+build !noasm
//+build !appengine

package mat

import "unsafe"

//go:noescape
func __CrossProduct(vec1, vec2, result unsafe.Pointer)

func CrossProduct(f1, f2, out *Tuple4) {
	__CrossProduct(unsafe.Pointer(f1), unsafe.Pointer(f2), unsafe.Pointer(out))
}