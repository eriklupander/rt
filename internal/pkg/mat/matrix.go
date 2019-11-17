package mat

var IdentityMatrix = &Mat4x4{[]float64{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1},
}

func NewIdentityMatrix() *Mat4x4 {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	return m1
}

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
	return m.Elems[(row*2)+col]
}
func (m Mat3x3) Get(row int, col int) float64 {
	return m.Elems[(row*3)+col]
}
func (m Mat4x4) Get(row int, col int) float64 {
	return m.Elems[(row*4)+col]
}

func Equals(m1, m2 Mat4x4) bool {
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if m1.Get(row, col) != m2.Get(row, col) {
				return false
			}
		}
	}
	return true
}
func Equals3x3(m1, m2 Mat3x3) bool {
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			if m1.Get(row, col) != m2.Get(row, col) {
				return false
			}
		}
	}
	return true
}
func Equals2x2(m1, m2 Mat2x2) bool {
	for row := 0; row < 2; row++ {
		for col := 0; col < 2; col++ {
			if m1.Get(row, col) != m2.Get(row, col) {
				return false
			}
		}
	}
	return true
}

func Multiply(m1 *Mat4x4, m2 *Mat4x4) *Mat4x4 {
	m3 := NewMat4x4(make([]float64, 16))
	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			m3.Elems[(row*4)+col] = multiply4x4(m1, m2, row, col)
		}
	}
	return m3
}

func MultiplyByTuple(m1 Mat4x4, t Tuple4) *Tuple4 {
	t1 := NewTuple4(make([]float64, 4))
	for row := 0; row < 4; row++ {
		t1.Elems[row] = (m1.Get(row, 0) * t.Get(0)) +
			(m1.Get(row, 1) * t.Get(1)) +
			(m1.Get(row, 2) * t.Get(2)) +
			(m1.Get(row, 3) * t.Get(3))
	}
	return t1
}

// Transpose flips rows and cols in the matrix.
func Transpose(m1 Mat4x4) *Mat4x4 {
	m3 := NewMat4x4(make([]float64, 16))
	for col := 0; col < 4; col++ {
		for row := 0; row < 4; row++ {
			m3.Elems[(row*4)+col] = m1.Get(col, row)
		}
	}
	return m3
}

// Determinant2x2 returns A-D minus B-C for a
// [A,B]
// [C,D]
// Matrix.
func Determinant2x2(m1 Mat2x2) float64 {
	return m1.Elems[0]*m1.Elems[3] - m1.Elems[1]*m1.Elems[2]
}

// Determinant3x3 takes the first row of the passed matrix, summing the colvalue * Cofactor of the same col
func Determinant3x3(m1 *Mat3x3) float64 {
	det := 0.0
	for col := 0; col < 3; col++ {
		det = det + m1.Elems[col]*Cofactor3x3(m1, 0, col)
	}
	return det
}

// Determinant4x4
func Determinant4x4(m1 *Mat4x4) float64 {
	det := 0.0
	for col := 0; col < 4; col++ {
		det = det + m1.Elems[col]*Cofactor4x4(m1, 0, col)
	}
	return det
}

// Submatrix3x3 extracts the 2x2 submatrix after deleting row and col from the passed 3x3
func Submatrix3x3(m1 Mat3x3, deleteRow, deleteCol int) *Mat2x2 {
	m3 := NewMat2x2(make([]float64, 4))
	idx := 0
	for row := 0; row < 3; row++ {
		if row == deleteRow {
			continue
		}
		for col := 0; col < 3; col++ {
			if col == deleteCol {
				continue
			}

			m3.Elems[idx] = m1.Get(row, col)
			idx++
		}
	}
	return m3
}

// Submatrix4x4 extracts the 3x3 submatrix after deleting row and col from the passed 4x4
func Submatrix4x4(m1 Mat4x4, deleteRow, deleteCol int) *Mat3x3 {
	m3 := NewMat3x3(make([]float64, 9))
	idx := 0
	for row := 0; row < 4; row++ {
		if row == deleteRow {
			continue
		}
		for col := 0; col < 4; col++ {
			if col == deleteCol {
				continue
			}
			m3.Elems[idx] = m1.Get(row, col)
			idx++
		}
	}
	return m3
}

// Minor3x3 computes the submatrix at row/col and returns the determinant of the computed matrix.
func Minor3x3(m1 *Mat3x3, row, col int) float64 {
	m2 := Submatrix3x3(*m1, row, col)
	return Determinant2x2(*m2)
}

// Minor4x4 computes the submatrix at row/col and returns the determinant of the computed matrix.
func Minor4x4(m1 *Mat4x4, row, col int) float64 {
	m2 := Submatrix4x4(*m1, row, col)
	return Determinant3x3(m2)
}

// Cofactor3x3 may change the sign of the computed minor of the passed matrix
func Cofactor3x3(m1 *Mat3x3, row, col int) float64 {
	minor := Minor3x3(m1, row, col)
	if (row+col)%2 != 0 {
		return -minor
	}
	return minor
}

// Cofactor4x4 may change the sign of the computed minor of the passed matrix
func Cofactor4x4(m1 *Mat4x4, row, col int) float64 {
	minor := Minor4x4(m1, row, col)
	if (row+col)%2 != 0 {
		return -minor
	}
	return minor
}

func IsInvertible(m1 *Mat4x4) bool {
	return Determinant4x4(m1) != 0.0
}

func Inverse(m1 *Mat4x4) *Mat4x4 {
	m3 := NewMat4x4(make([]float64, 16))
	d4 := Determinant4x4(m1)
	for row := 0; row < 4; row++ {

		for col := 0; col < 4; col++ {
			c := Cofactor4x4(m1, row, col)
			// note that "col, row" here, instead of "row, col",
			// accomplishes the transpose operation!
			m3.Elems[col*4+row] = c / d4
		}
	}
	return m3
}

func multiply4x4(m1 *Mat4x4, m2 *Mat4x4, row int, col int) float64 {
	// always row from m1, col from m2
	a0 := m1.Get(row, 0) * m2.Get(0, col)
	a1 := m1.Get(row, 1) * m2.Get(1, col)
	a2 := m1.Get(row, 2) * m2.Get(2, col)
	a3 := m1.Get(row, 3) * m2.Get(3, col)
	return a0 + a1 + a2 + a3
}
