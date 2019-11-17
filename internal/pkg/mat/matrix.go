package mat



type Mat2x2 struct {
	Elems []float64
}
func NewMat2x2(elems []float64) *Mat2x2 {
	return &Mat2x2{Elems: elems}
}

type Mat3x3 struct {
	Elems []float64
}
func NewMat3x3(elems []float64) *Mat3x3 {
	return &Mat3x3{Elems: elems}
}

type Mat4x4 struct {
	Elems []float64
}
func NewMat4x4(elems []float64) *Mat4x4 {
	return &Mat4x4{Elems: elems}
}

func (m Mat2x2) Get(row int, col int) float64 {
	return m.Elems[(row * 2) + col]
}
func (m Mat3x3) Get(row int, col int) float64 {
	return m.Elems[(row * 3) + col]
}
func (m Mat4x4) Get(row int, col int) float64 {
	return m.Elems[(row * 4) + col]
}

func Equals(m1, m2 Mat4x4) bool {
	for row:=0; row < 4; row++ {
		for col:=0; col < 4; col++ {
			if m1.Get(row, col) != m2.Get(row, col) {
				return false
			}
		}
	}
	return true
}


func Multiply(m1 *Mat4x4, m2 *Mat4x4) *Mat4x4 {
	m3 := NewMat4x4(make([]float64, 16))
	for row:=0; row < 4; row++ {
		for col:=0; col < 4; col++ {
			m3.Elems[(row * 4) + col] = multiply4x4(m1, m2, row, col)
		}
	}
	return m3
}

func MultiplyByTuple(m1 Mat4x4, t Tuple4) *Tuple4 {
	t1 := NewTuple4(make([]float64, 4))
	for row:=0; row < 4; row++ {
		t1.Elems[row] = (m1.Get(row, 0) * t.Get(0)) +
			(m1.Get(row, 1) * t.Get(1)) +
			(m1.Get(row, 2) * t.Get(2)) +
			(m1.Get(row, 3) * t.Get(3))
	}
	return t1
}

func multiply4x4(m1 *Mat4x4, m2 *Mat4x4, row int, col int) float64 {
	// always row from m1, col from m2
	a0 := m1.Get(row, 0) * m2.Get(0, col)
	a1 := m1.Get(row, 1) * m2.Get(1, col)
	a2 := m1.Get(row, 2) * m2.Get(2, col)
	a3 := m1.Get(row, 3) * m2.Get(3, col)
	return a0 + a1 + a2 + a3
}