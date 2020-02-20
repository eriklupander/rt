package mat

func Scale(x, y, z float64) Mat4x4 {
	scaleMatrix := IdentityMatrix // NewMat4x4(make([]float64, 16))
	// copy(scaleMatrix.Elems, IdentityMatrix.Elems)
	scaleMatrix[0] = x
	scaleMatrix[5] = y
	scaleMatrix[10] = z
	return scaleMatrix
}
