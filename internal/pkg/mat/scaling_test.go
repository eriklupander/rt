package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Note how the point x,y,z is scaled by xyz
func TestScale(t *testing.T) {
	scaleTransform := Scale(2, 3, 4)
	p := NewPoint(-4, 6, 8)
	p2 := MultiplyByTuple(*scaleTransform, *p)

	assert.Equal(t, -8.0, p2.Get(0))
	assert.Equal(t, 18.0, p2.Get(1))
	assert.Equal(t, 32.0, p2.Get(2))
}

func TestScaleVector(t *testing.T) {
	scaleTransform := Scale(2, 3, 4)
	v := NewVector(-4, 6, 8)
	p2 := MultiplyByTuple(*scaleTransform, *v)

	assert.Equal(t, -8.0, p2.Get(0))
	assert.Equal(t, 18.0, p2.Get(1))
	assert.Equal(t, 32.0, p2.Get(2))
}

// Note how scaling by inverse effectively divides v xyz by scale xyz
func TestScaleByInverse(t *testing.T) {
	scaleTransform := Scale(2, 3, 4)
	v := NewVector(-4, 6, 8)
	inv := Inverse(scaleTransform)
	p2 := MultiplyByTuple(*inv, *v)
	assert.Equal(t, -2.0, p2.Get(0))
	assert.Equal(t, 2.0, p2.Get(1))
	assert.Equal(t, 2.0, p2.Get(2))
}

// Reflect (e.g.) scale by negative. Flips sign of X in this case.
func TestReflect(t *testing.T) {
	scaleTransform := Scale(-1, 1, 1)
	p := NewPoint(2, 3, 4)
	p2 := MultiplyByTuple(*scaleTransform, *p)
	assert.Equal(t, -2.0, p2.Get(0))
	assert.Equal(t, 3.0, p2.Get(1))
	assert.Equal(t, 4.0, p2.Get(2))
}
