package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

// Note how the translation x,y,z is added to the point xyz
func TestTranslate(t *testing.T) {
	m1 := Translate(5, -3, 2)
	p := NewPoint(-3, 4, 5)
	p2 := MultiplyByTuple(*m1, *p)
	assert.Equal(t, 2.0, p2.Get(0))
	assert.Equal(t, 1.0, p2.Get(1))
	assert.Equal(t, 7.0, p2.Get(2))
}

// Note how translating by the inverse subtracts the translation from the point
func TestMultiplyByInverseOfTranslationMatrix(t *testing.T) {
	m1 := Translate(5, -3, 2)
	p := NewPoint(-3, 4, 5)
	inverseMatrix := Inverse(m1)

	p2 := MultiplyByTuple(*inverseMatrix, *p)
	assert.Equal(t, -8.0, p2.Get(0))
	assert.Equal(t, 7.0, p2.Get(1))
	assert.Equal(t, 3.0, p2.Get(2))
}

// Remember, the w (as in xyzw) is always 0 on a vector, causing no effect when a multiplication of a translation matrix occurrs.
func TestVectorNotAffectedByTranslation(t *testing.T) {
	m1 := Translate(5, -3, 2)
	v := NewVector(-3, 4, 5)
	p2 := MultiplyByTuple(*m1, *v)

	assert.True(t, TupleEquals(*p2, *v))
}

func TestApplyInSequence(t *testing.T) {
	p := NewPoint(1, 0, 1)
	rotA := RotateX(math.Pi / 2)
	scaleB := Scale(5, 5, 5)
	translateC := Translate(10, 5, 7)

	p2 := MultiplyByTuple(*rotA, *p)
	p3 := MultiplyByTuple(*scaleB, *p2)
	p4 := MultiplyByTuple(*translateC, *p3)

	assert.True(t, TupleEquals(*p4, *NewPoint(15, 0, 7)))
}

func TestApplyChained(t *testing.T) {
	p := NewPoint(1, 0, 1)
	rotA := RotateX(math.Pi / 2)
	scaleB := Scale(5, 5, 5)
	translateC := Translate(10, 5, 7)

	// Note that we chain with last-first (IIRC, if we rotate first and then translate, the translation will be off by rot)
	m3 := Multiply(Multiply(translateC, scaleB), rotA)

	p4 := MultiplyByTuple(*m3, *p)
	assert.True(t, TupleEquals(*p4, *NewPoint(15, 0, 7)))
}

func TestTranslateSphere(t *testing.T) {
	s := NewSphere()
	t1 := Translate(2, 3, 4)
	SetTransform(s, t1)
	assert.True(t, Equals(*s.Transform, *t1))
}
