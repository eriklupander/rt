//+build !noasm
//+build !appengine

package mat

import "unsafe"

//go:noescape
func __MultiplyMatrixByVec64(m, vec4, result unsafe.Pointer)

//func MultiplyByTuplePtr(m *[16]float64, f2 *[4]float64, _f4 *[4]float64)  {
//	__MultiplyMatrixByVec64(unsafe.Pointer(m), unsafe.Pointer(f2), unsafe.Pointer(_f4))
//}
func MultiplyByTuplePtr2(m Mat4x4, f2 Tuple4, _f4 *Tuple4) {
	__MultiplyMatrixByVec64(unsafe.Pointer(&m), unsafe.Pointer(&f2), unsafe.Pointer(_f4))
}
func MultiplyByTuplePtr(m *Mat4x4, f2 *Tuple4, _f4 *Tuple4) {
	__MultiplyMatrixByVec64(unsafe.Pointer(m), unsafe.Pointer(f2), unsafe.Pointer(_f4))
}
