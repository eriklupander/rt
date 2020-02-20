package mat

// Translate creates a translation matrix from the identity matrix, setting wx, wy, wz to the
// passed xyz coords.
func Translate(x, y, z float64) Mat4x4 {
	m1 := New4x4() //NewMat4x4(make([]float64, 16))
	//copy(m1.Elems, IdentityMatrix.Elems)
	m1[3] = x
	m1[7] = y
	m1[11] = z
	return m1
}
