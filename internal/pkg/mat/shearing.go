package mat

func Shear(xy, xz, yx, yz, zx, zy float64) Mat4x4 {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	m1.Elems[1] = xy
	m1.Elems[2] = xz
	m1.Elems[4] = yx
	m1.Elems[6] = yz
	m1.Elems[8] = zx
	m1.Elems[9] = zy
	return m1
}
func ShearBy(args []float64) Mat4x4 {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	m1.Elems[1] = args[0]
	m1.Elems[2] = args[1]
	m1.Elems[4] = args[2]
	m1.Elems[6] = args[3]
	m1.Elems[8] = args[4]
	m1.Elems[9] = args[5]
	return m1
}
