package mat

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateMatrix2x2(t *testing.T) {
	m2 := NewMat2x2([]float64{-3, 5, 1, -2})
	assert.Equal(t, -3.0, m2.Get(0, 0))
	assert.Equal(t, 5.0, m2.Get(0, 1))
	assert.Equal(t, 1.0, m2.Get(1, 0))
	assert.Equal(t, -2.0, m2.Get(1, 1))
}

func TestCreateMatrix3x3(t *testing.T) {
	m3 := NewMat3x3([]float64{-3, 5, 0, 1, -2, -7, 0, 1, 1})
	assert.Equal(t, -3.0, m3.Get(0, 0))
	assert.Equal(t, -2.0, m3.Get(1, 1))
	assert.Equal(t, 1.0, m3.Get(2, 2))
}

func TestCreateMatrix4x4(t *testing.T) {
	m4 := NewMat4x4([]float64{1, 2, 3, 4, 5.5, 6.5, 7.5, 8.5, 9, 10, 11, 12, 13.5, 14.5, 15.5, 16.5})

	assert.Equal(t, 1.0, m4.Get(0, 0))
	assert.Equal(t, 4.0, m4.Get(0, 3))
	assert.Equal(t, 5.5, m4.Get(1, 0))
	assert.Equal(t, 7.5, m4.Get(1, 2))
	assert.Equal(t, 11.0, m4.Get(2, 2))
	assert.Equal(t, 13.5, m4.Get(3, 0))
	assert.Equal(t, 15.5, m4.Get(3, 2))
}

func TestCompare4x4(t *testing.T) {
	m1 := NewMat4x4([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2})
	m2 := NewMat4x4([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2})

	assert.True(t, Equals(m1, m2))
}
func TestCompare4x4NotEqual(t *testing.T) {
	m1 := NewMat4x4([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2})
	m2 := NewMat4x4([]float64{2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2, 1})

	assert.False(t, Equals(m1, m2))
}

func TestMultiply4x4(t *testing.T) {
	m1 := NewMat4x4([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2})
	m2 := NewMat4x4([]float64{-2, 1, 2, 3, 3, 2, 1, -1, 4, 3, 6, 5, 1, 2, 7, 8})

	m3 := Multiply(m1, m2)

	assert.Equal(t, 20.0, m3.Get(0, 0))
	assert.Equal(t, 54.0, m3.Get(1, 1))
	assert.Equal(t, 110.0, m3.Get(2, 2))
	assert.Equal(t, 42.0, m3.Get(3, 3))
}

func TestMultiply4x4ByTuple(t *testing.T) {
	m1 := NewMat4x4([]float64{1, 2, 3, 4, 2, 4, 4, 2, 8, 6, 4, 1, 0, 0, 0, 1})

	t1 := [4]float64{1, 2, 3, 1} //NewTuple4([]float64{1, 2, 3, 1})

	t2 := MultiplyByTuple(m1, t1)
	assert.Equal(t, 18.0, t2.Get(0))
	assert.Equal(t, 24.0, t2.Get(1))
	assert.Equal(t, 33.0, t2.Get(2))
	assert.Equal(t, 1.0, t2.Get(3))
}

func TestMultiplyByIdentityMatrix(t *testing.T) {
	m1 := NewMat4x4([]float64{0, 1, 2, 4, 1, 2, 4, 8, 2, 4, 8, 16, 4, 8, 16, 32})
	m3 := Multiply(m1, IdentityMatrix)
	assert.True(t, Equals(m1, m3))
}

func TestMultiplyTupleByIdentityMatrix(t *testing.T) {
	t1 := NewTupleOf(1, 2, 3, 4)
	t3 := MultiplyByTuple(IdentityMatrix, t1)
	assert.True(t, TupleEquals(t1, t3))
}

func TestTransposeMatrix(t *testing.T) {
	m1 := NewMat4x4([]float64{0, 9, 3, 0, 9, 8, 0, 8, 1, 8, 5, 3, 0, 0, 5, 8})
	expected := NewMat4x4([]float64{0, 9, 1, 0, 9, 8, 8, 0, 3, 0, 5, 5, 0, 8, 3, 8})
	m3 := Transpose(m1)
	assert.True(t, Equals(m3, expected))
}

func TestTransposeIdentityMatrix(t *testing.T) {
	m3 := Transpose(IdentityMatrix)
	assert.True(t, Equals(m3, IdentityMatrix))
}

func TestDeterminant2x2(t *testing.T) {
	m1 := NewMat2x2([]float64{1, 5, -3, 2})
	determinant := Determinant2x2(m1)
	assert.Equal(t, 17.0, determinant)
}

func TestSubmatrix3x3(t *testing.T) {
	m1 := NewMat3x3([]float64{1, 5, 0, -3, 2, 7, 0, 6, -3})
	expected := NewMat2x2([]float64{-3, 2, 0, 6})
	m3 := Submatrix3x3(m1, 0, 2)
	assert.True(t, Equals2x2(m3, expected))
}

func TestSubmatrix4x4(t *testing.T) {
	m1 := NewMat4x4([]float64{-6, 1, 1, 6, -8, 5, 8, 6, -1, 0, 8, 2, -7, 1, -1, 1})
	expected := NewMat3x3([]float64{-6, 1, 6, -8, 8, 6, -7, -1, 1})
	m3 := Submatrix4x4(m1, 2, 1)
	assert.True(t, Equals3x3(m3, expected))
}

func TestMinor(t *testing.T) {
	m1 := NewMat3x3([]float64{3, 5, 0, 2, -1, -7, 6, -1, 5})
	m2 := Submatrix3x3(m1, 1, 0)
	determ := Determinant2x2(m2)
	assert.Equal(t, 25.0, determ)

	minor := Minor3x3(m1, 1, 0)
	assert.Equal(t, 25.0, minor)
}

func TestCofactor3x3(t *testing.T) {
	m1 := NewMat3x3([]float64{3, 5, 0, 2, -1, -7, 6, -1, 5})
	minor1 := Minor3x3(m1, 0, 0)
	cofactor1 := Cofactor3x3(m1, 0, 0)

	minor2 := Minor3x3(m1, 1, 0)
	cofactor2 := Cofactor3x3(m1, 1, 0)

	assert.Equal(t, -12.0, minor1)
	assert.Equal(t, -12.0, cofactor1)
	assert.Equal(t, 25.0, minor2)
	assert.Equal(t, -25.0, cofactor2)
}

func TestDeterminant3x3(t *testing.T) {
	m1 := NewMat3x3([]float64{1, 2, 6, -5, 8, -4, 2, 6, 4})
	cf1 := Cofactor3x3(m1, 0, 0)
	cf2 := Cofactor3x3(m1, 0, 1)
	cf3 := Cofactor3x3(m1, 0, 2)
	determinant := Determinant3x3(m1)

	assert.Equal(t, 56.0, cf1)
	assert.Equal(t, 12.0, cf2)
	assert.Equal(t, -46.0, cf3)
	assert.Equal(t, -196.0, determinant)
}

func TestDeterminant4x4(t *testing.T) {
	m1 := NewMat4x4([]float64{-2, -8, 3, 5, -3, 1, 7, 3, 1, 2, -9, 6, -6, 7, 7, -9})
	cf1 := Cofactor4x4(m1, 0, 0)
	cf2 := Cofactor4x4(m1, 0, 1)
	cf3 := Cofactor4x4(m1, 0, 2)
	cf4 := Cofactor4x4(m1, 0, 3)
	determinant := Determinant4x4(m1)

	assert.Equal(t, 690.0, cf1)
	assert.Equal(t, 447.0, cf2)
	assert.Equal(t, 210.0, cf3)
	assert.Equal(t, 51.0, cf4)
	assert.Equal(t, -4071.0, determinant)
}

func TestIsInvertible(t *testing.T) {
	m1 := NewMat4x4([]float64{6, 4, 4, 4, 5, 5, 7, 6, 4, -9, 3, -7, 9, 1, 7, -6})
	determinant := Determinant4x4(m1)
	assert.Equal(t, -2120.0, determinant)

	isInvertible := IsInvertible(m1)
	assert.True(t, isInvertible)
}

func TestIsNotInvertible(t *testing.T) {
	m1 := NewMat4x4([]float64{-4, 2, -2, -3, 9, 6, 2, 6, 0, -5, 1, -5, 0, 0, 0, 0})
	determinant := Determinant4x4(m1)
	assert.Equal(t, 0.0, determinant)

	isInvertible := IsInvertible(m1)
	assert.False(t, isInvertible)
}

func TestInverse(t *testing.T) {
	m1 := NewMat4x4([]float64{-5, 2, 6, -8, 1, -5, 1, 8, 7, 7, -6, -7, 1, -3, 7, 4})
	m3 := Inverse(m1)

	cf1 := Cofactor4x4(m1, 2, 3)
	cf2 := Cofactor4x4(m1, 3, 2)
	determinant := Determinant4x4(m1)
	assert.Equal(t, 532.0, determinant)
	assert.Equal(t, -160.0, cf1)
	assert.Equal(t, 105.0, cf2)

	expected := NewMat4x4([]float64{0.21805, 0.45113, 0.24060, -0.04511,
		-0.80827, -1.45677, -0.44361, 0.52068,
		-0.07895, -0.22368, -0.05263, 0.19737,
		-0.52256, -0.81391, -0.30075, 0.30639})

	for i := 0; i < 16; i++ {
		assert.InEpsilon(t, expected[i], m3[i], Epsilon, fmt.Sprintf("index %d failed: values: %v %v", i, expected[i], m3[i]))
	}
}

func TestInverse2(t *testing.T) {
	m1 := NewMat4x4([]float64{8, -5, 9, 2, 7, 5, 6, 1, -6, 0, 9, 6, -3, 0, -9, -4})
	m3 := Inverse(m1)
	expected := NewMat4x4([]float64{-0.15385, -0.15385, -0.28205, -0.53846,
		-0.07692, 0.12308, 0.02564, 0.03077,
		0.35897, 0.35897, 0.43590, 0.92308,
		-0.69231, -0.69231, -0.76923, -1.92308})

	for i := 0; i < 16; i++ {
		assert.InEpsilon(t, expected[i], m3[i], Epsilon, fmt.Sprintf("index %d failed: values: %v %v", i, expected[i], m3[i]))
	}
}
func TestInverse3(t *testing.T) {
	m1 := NewMat4x4([]float64{
		9, 3, 0, 9,
		-5, -2, -6, -3,
		-4, 9, 6, 4,
		-7, 6, 6, 2})
	m3 := Inverse(m1)

	expected := NewMat4x4([]float64{-0.04074, -0.07778, 0.14444, -0.22222,
		-0.07778, 0.03333, 0.36667, -0.33333,
		-0.02901, -0.14630, -0.10926, 0.12963,
		0.17778, 0.06667, -0.26667, 0.33333})

	for i := 0; i < 16; i++ {
		assert.InEpsilon(t, expected[i], m3[i], Epsilon, fmt.Sprintf("index %d failed: values: %v %v", i, expected[i], m3[i]))
	}
}

func TestMultiplyByInverse(t *testing.T) {
	m1 := NewMat4x4([]float64{
		3, -9, 7, 3,
		3, -8, 2, -9,
		-4, 4, 4, 1,
		-6, 5, -1, 1})

	m2 := NewMat4x4([]float64{
		8, 2, 2, 2,
		3, -1, 7, 0,
		7, 0, 5, 4,
		6, -2, 0, 5})

	m3 := Multiply(m1, m2)

	Equals(Multiply(m3, Inverse(m2)), m1)
}

func TestInvertIdentity(t *testing.T) {
	// Inverse identity matrix seems to do nothing.
	m3 := Inverse(IdentityMatrix)
	fmt.Printf("%+v\n", m3)

	m1 := NewMat4x4([]float64{
		3, -9, 7, 3,
		3, -8, 2, -9,
		-4, 4, 4, 1,
		-6, 5, -1, 1})

	// Multiply by its own inverse gets you identity matrix?
	m2 := Inverse(m1)
	m3 = Multiply(m1, m2)
	fmt.Printf("%+v\n", m3)

	// The transpose + inverse is equal to the inverse + transpose
	t2 := Transpose(m1)
	i3 := Inverse(t2)

	i1 := Inverse(m1)
	t3 := Transpose(i1)

	fmt.Printf("%+v\n", i3)
	fmt.Printf("%+v\n", t3)

	tuple := NewTupleOf(1, 2, 3, 0)
	id1 := NewIdentityMatrix()
	firstTuple := MultiplyByTuple(id1, tuple)

	id1[5] = 7.0
	secondTuple := MultiplyByTuple(id1, tuple)

	fmt.Printf("%+v\n", firstTuple)
	fmt.Printf("%+v\n", secondTuple)
}
