package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNormalOnSphereAtX(t *testing.T) {
	s := NewSphere()
	normalVector := NormalAt(s, NewPoint(1, 0, 0))
	assert.True(t, TupleEquals(normalVector, NewVector(1, 0, 0)))
}
func TestNormalOnSphereAtY(t *testing.T) {
	s := NewSphere()
	normalVector := NormalAt(s, NewPoint(0, 1, 0))
	assert.True(t, TupleEquals(normalVector, NewVector(0, 1, 0)))
}
func TestNormalAtPointOnSphereAtZ(t *testing.T) {
	s := NewSphere()
	normalVector := NormalAt(s, NewPoint(0, 0, 1))
	assert.True(t, TupleEquals(normalVector, NewVector(0, 0, 1)))
}
func TestNormalOnSphereAtNonAxial(t *testing.T) {
	s := NewSphere()
	nonAxial := math.Sqrt(3.0) / 3.0
	normalVector := NormalAt(s, NewPoint(nonAxial, nonAxial, nonAxial))
	assert.InEpsilon(t, nonAxial, normalVector.Get(0), Epsilon)
	assert.InEpsilon(t, nonAxial, normalVector.Get(1), Epsilon)
	assert.InEpsilon(t, nonAxial, normalVector.Get(2), Epsilon)
}
func TestNormalIsNormalized(t *testing.T) {
	s := NewSphere()
	nonAxial := math.Sqrt(3.0) / 3.0
	normalVector := NormalAt(s, NewPoint(nonAxial, nonAxial, nonAxial))
	normalizedNormalVector := Normalize(normalVector)
	assert.True(t, TupleEquals(normalVector, normalizedNormalVector))
}
func TestComputeNormalOnTranslatedSphere(t *testing.T) {

	s := NewSphere()
	s.SetTransform(Translate(0, 1, 0))
	normalVector := NormalAt(s, NewPoint(0, 1.70711, -0.70711))
	assert.Equal(t, 0.0, normalVector.Get(0))
	assert.InEpsilon(t, 0.70711, normalVector.Get(1), Epsilon)
	assert.InEpsilon(t, -0.70711, normalVector.Get(2), Epsilon)
}

func TestComputeNormalOnTransformedSphere(t *testing.T) {

	s := NewSphere()
	m1 := Multiply(Scale(1, 0.5, 1), RotateZ(math.Pi/5.0))
	s.SetTransform(m1)
	normalVector := NormalAt(s, NewPoint(0, math.Sqrt(2)/2, -math.Sqrt(2)/2))
	assert.Equal(t, 0.0, normalVector.Get(0))
	assert.InEpsilon(t, 0.97014, normalVector.Get(1), Epsilon)
	assert.InEpsilon(t, -0.24254, normalVector.Get(2), Epsilon)
}

// Reflecting a vector approaching at 45°
func TestReflectRay(t *testing.T) {
	v := NewVector(1, -1, 0)
	normal := NewVector(0, 1, 0) // straight up
	reflectV := Reflect(v, normal)
	assert.True(t, TupleEquals(NewVector(1, 1, 0), reflectV))
}
func TestReflectRaySlanted(t *testing.T) {
	v := NewVector(0, -1, 0) // Pointing straight down
	fortyFive := math.Sqrt(2) / 2.0
	normal := NewVector(fortyFive, fortyFive, 0) // straight up
	reflectV := Reflect(v, normal)
	assert.True(t, TupleEquals(NewVector(1, 0, 0), reflectV))
}
