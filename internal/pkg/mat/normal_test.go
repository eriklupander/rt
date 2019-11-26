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

func TestPrecomputingReflectionVector(t *testing.T) {
	pl := NewPlane()
	ray := NewRay(NewPoint(0, 1, -1), NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := NewIntersection(math.Sqrt(2), pl)
	comps := PrepareComputationForIntersection(xs, ray)
	assert.Equal(t, 0.0, comps.ReflectVec.Get(0))
	assert.InEpsilon(t, math.Sqrt(2)/2, comps.ReflectVec.Get(1), Epsilon)
	assert.InEpsilon(t, math.Sqrt(2)/2, comps.ReflectVec.Get(2), Epsilon)
}

func TestReflectedColorForNonreflectiveMaterial(t *testing.T) {
	w := NewDefaultWorld()
	r := NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1))
	shape := w.Objects[1]
	material := shape.GetMaterial()
	material.Ambient = 1.0
	shape.SetMaterial(material)
	xs := NewIntersection(1, shape)
	comps := PrepareComputationForIntersection(xs, r)
	color := ReflectedColor(w, comps, 1)
	assert.Equal(t, black, color)
}

func TestReflectedColorForReflectiveMaterial(t *testing.T) {
	w := NewDefaultWorld()
	plane := NewPlane()
	plane.SetTransform(Translate(0, -1, -0))
	material := plane.GetMaterial()
	material.Reflectivity = 0.5
	plane.SetMaterial(material)
	w.Objects = append(w.Objects, plane)

	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := NewIntersection(math.Sqrt(2), plane)
	comps := PrepareComputationForIntersection(xs, r)
	color := ReflectedColor(w, comps, 1)
	assert.InEpsilon(t, 0.19032, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.2379, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.14274, color.Get(2), Epsilon)
}

func TestShadeHitWithReflectiveMaterial(t *testing.T) {
	w := NewDefaultWorld()
	plane := NewPlane()
	plane.SetTransform(Translate(0, -1, -0))
	material := plane.GetMaterial()
	material.Reflectivity = 0.5
	plane.SetMaterial(material)
	w.Objects = append(w.Objects, plane)

	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := NewIntersection(math.Sqrt(2), plane)
	comps := PrepareComputationForIntersection(xs, r)
	color := ShadeHit(w, comps, 1)

	assert.InEpsilon(t, 0.87677, color.Get(0), Epsilon)
	assert.InEpsilon(t, 0.92436, color.Get(1), Epsilon)
	assert.InEpsilon(t, 0.82918, color.Get(2), Epsilon)
}

func TestColorAtWithMutuallyReflectiveSurfaces(t *testing.T) {
	w := NewWorld()
	w.Light = NewLight(NewPoint(0, 0, 0), NewColor(1, 1, 1))
	lowerPlane := NewPlane()
	lowerPlane.SetMaterial(NewDefaultReflectiveMaterial(1.0))
	lowerPlane.SetTransform(Translate(0, -1, 0))
	w.Objects = append(w.Objects, lowerPlane)

	upperPlane := NewPlane()
	upperPlane.SetMaterial(NewDefaultReflectiveMaterial(1.0))
	upperPlane.SetTransform(Translate(0, 1, 0))
	w.Objects = append(w.Objects, upperPlane)

	r := NewRay(NewPoint(0, 0, 0), NewVector(0, 1, 0))
	_ = ColorAt(w, r, 1)
}

func TestTheReflectedColorAtMaxRecursiveDepth(t *testing.T) {
	w := NewWorld()
	pl := NewPlane()
	pl.SetMaterial(NewDefaultReflectiveMaterial(0.5))
	pl.SetTransform(Translate(0, -1, 0))
	w.Objects = append(w.Objects, pl)
	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := NewIntersection(math.Sqrt(2), pl)
	comps := PrepareComputationForIntersection(xs, r)
	color := ReflectedColor(w, comps, 0)
	assert.Equal(t, black, color)
}
