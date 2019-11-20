package mat

import "math"

func RotateX(radians float64) Mat4x4 {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	m1.Elems[5] = math.Cos(radians)
	m1.Elems[6] = -math.Sin(radians)
	m1.Elems[9] = math.Sin(radians)
	m1.Elems[10] = math.Cos(radians)
	return m1
}

func RotateY(radians float64) Mat4x4 {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	m1.Elems[0] = math.Cos(radians)
	m1.Elems[2] = math.Sin(radians)
	m1.Elems[8] = -math.Sin(radians)
	m1.Elems[10] = math.Cos(radians)
	return m1
}

func RotateZ(radians float64) Mat4x4 {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	m1.Elems[0] = math.Cos(radians)
	m1.Elems[1] = -math.Sin(radians)
	m1.Elems[4] = math.Sin(radians)
	m1.Elems[5] = math.Cos(radians)
	return m1
}
