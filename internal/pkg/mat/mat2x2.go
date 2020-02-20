package mat

type Mat2x2 [4]float64

func NewMat2x2(elems []float64) Mat2x2 {
	return Mat2x2{elems[0], elems[1], elems[2], elems[3]}
}

func (m Mat2x2) Get(row int, col int) float64 {
	return m[(row*2)+col]
}
