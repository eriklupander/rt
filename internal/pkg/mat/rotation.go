package mat

import "math"

func RotateX(radians float64) Mat4x4 {
	m1 := New4x4() // NewMat4x4(make([]float64, 16))
	//copy(m1.Elems, IdentityMatrix.Elems)
	m1[5] = math.Cos(radians)
	m1[6] = -math.Sin(radians)
	m1[9] = math.Sin(radians)
	m1[10] = math.Cos(radians)
	return m1
}

func RotateY(radians float64) Mat4x4 {
	m1 := New4x4() //NewMat4x4(make([]float64, 16))
	//copy(m1.Elems, IdentityMatrix.Elems)
	m1[0] = math.Cos(radians)
	m1[2] = math.Sin(radians)
	m1[8] = -math.Sin(radians)
	m1[10] = math.Cos(radians)
	return m1
}

func RotateZ(radians float64) Mat4x4 {
	m1 := New4x4() //NewMat4x4(make([]float64, 16))
	// copy(m1.Elems, IdentityMatrix.Elems)
	m1[0] = math.Cos(radians)
	m1[1] = -math.Sin(radians)
	m1[4] = math.Sin(radians)
	m1[5] = math.Cos(radians)
	return m1
}
