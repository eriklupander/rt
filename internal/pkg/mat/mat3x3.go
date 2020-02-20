package mat

type Mat3x3 [9]float64

func NewMat3x3(elems []float64) Mat3x3 {
	return Mat3x3{elems[0], elems[1], elems[2], elems[3], elems[4], elems[5], elems[6], elems[7], elems[8]}
}

func (m Mat3x3) Get(row int, col int) float64 {
	return m[(row*3)+col]
}
