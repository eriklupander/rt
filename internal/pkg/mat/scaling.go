package mat

func Scale(x, y, z float64) Mat4x4 {
	scaleMatrix := NewMat4x4(make([]float64, 16))
	copy(scaleMatrix.Elems, IdentityMatrix.Elems)
	scaleMatrix.Elems[0] = x
	scaleMatrix.Elems[5] = y
	scaleMatrix.Elems[10] = z
	return scaleMatrix
}
