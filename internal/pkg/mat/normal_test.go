package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNormalOnSphereAtX(t *testing.T) {
	s := NewSphere()
	normalVector := NormalAt(s, 1, 0, 0)
	assert.True(t, TupleEquals(*normalVector, *NewVector(1, 0, 0)))
}
func TestNormalOnSphereAtY(t *testing.T) {
	s := NewSphere()
	normalVector := NormalAt(s, 0, 1, 0)
	assert.True(t, TupleEquals(*normalVector, *NewVector(0, 1, 0)))
}
func TestNormalOnSphereAtZ(t *testing.T) {
	s := NewSphere()
	normalVector := NormalAt(s, 0, 0, 1)
	assert.True(t, TupleEquals(*normalVector, *NewVector(0, 0, 1)))
}
func TestNormalOnSphereAtNonAxial(t *testing.T) {
	s := NewSphere()
	nonAxial := math.Sqrt(3.0) / 3.0
	normalVector := NormalAt(s, nonAxial, nonAxial, nonAxial)
	assert.InEpsilon(t, nonAxial, normalVector.Get(0), Epsilon)
	assert.InEpsilon(t, nonAxial, normalVector.Get(1), Epsilon)
	assert.InEpsilon(t, nonAxial, normalVector.Get(2), Epsilon)
}
