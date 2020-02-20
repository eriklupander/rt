package mat

func Shear(xy, xz, yx, yz, zx, zy float64) Mat4x4 {
	m1 := New4x4() //NewMat4x4(make([]float64, 16))
	// copy(m1, IdentityMatrix)
	m1[1] = xy
	m1[2] = xz
	m1[4] = yx
	m1[6] = yz
	m1[8] = zx
	m1[9] = zy
	return m1
}
func ShearBy(args []float64) Mat4x4 {
	m1 := New4x4() //NewMat4x4(make([]float64, 16))
	// copy(m1, IdentityMatrix)
	m1[1] = args[0]
	m1[2] = args[1]
	m1[4] = args[2]
	m1[6] = args[3]
	m1[8] = args[4]
	m1[9] = args[5]
	return m1
}
