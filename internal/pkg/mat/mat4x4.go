package mat

type Mat4x4 [16]float64

func New4x4() Mat4x4 {
	return Mat4x4{1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func NewMat4x4(elems []float64) Mat4x4 {
	return Mat4x4{elems[0], elems[1], elems[2], elems[3],
		elems[4], elems[5], elems[6], elems[7], elems[8], elems[9],
		elems[10], elems[11], elems[12], elems[13], elems[14], elems[15]}
}
func (m Mat4x4) Get(row int, col int) float64 {
	return m[(row*4)+col]
}
