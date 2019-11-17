package mat

import (
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

	assert.True(t, Equals(*m1, *m2))
}
func TestCompare4x4NotEqual(t *testing.T) {
	m1 := NewMat4x4([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2})
	m2 := NewMat4x4([]float64{2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2, 1})

	assert.False(t, Equals(*m1, *m2))
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

	t1 := NewTuple4([]float64{1, 2, 3, 1})

	t2 := MultiplyByTuple(*m1, *t1)
	assert.Equal(t, 18.0, t2.Get(0))
	assert.Equal(t, 24.0, t2.Get(1))
	assert.Equal(t, 33.0, t2.Get(2))
	assert.Equal(t, 1.0, t2.Get(3))
}

func TestMultiplyByIdentityMatrix(t *testing.T) {
	m1 := NewMat4x4([]float64{0, 1, 2, 4, 1, 2, 4, 8, 2, 4, 8, 16, 4, 8, 16, 32})
	m3 := Multiply(m1, IdentityMatrix)
	assert.True(t, Equals(*m1, *m3))
}

func TestMultiplyTupleByIdentityMatrix(t *testing.T) {
	t1 := NewTuple4([]float64{1, 2, 3, 4})
	t3 := MultiplyByTuple(*IdentityMatrix, *t1)
	assert.True(t, TupleEquals(*t1, *t3))
}

func TestTransposeMatrix(t *testing.T) {
	m1 := NewMat4x4([]float64{0,9,3,0, 9,8,0,8, 1,8,5,3, 0,0,5,8})
	expected := NewMat4x4([]float64{0,9,1,0, 9,8,8,0, 3,0,5,5, 0,8,3,8 })
	m3 := Transpose(*m1)
	assert.True(t, Equals(*m3, *expected))
}

func TestTransportIdentityMatrix(t *testing.T) {
	m3 := Transpose(*IdentityMatrix)
	assert.True(t, Equals(*m3, *IdentityMatrix))
}

func TestDeterminant2x2(t *testing.T) {
	m1 := NewMat2x2([]float64{1, 5, -3, 2})
	determinant := Determinant2x2(*m1)
	assert.Equal(t, 17.0, determinant)
}

func TestSubmatrix3x3(t *testing.T) {
	m1 := NewMat3x3([]float64{1,5,0,  -3,2,7, 0,6,-3})
	expected := NewMat2x2([]float64{-3, 2, 0, 6})
	m3 := Submatrix3x3(*m1,0, 2)
	assert.True(t, Equals2x2(*m3, *expected))
}