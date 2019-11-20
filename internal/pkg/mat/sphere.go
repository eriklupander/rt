package mat

import (
	"math/rand"
)

func NewSphere() Sphere {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	return Sphere{Id: rand.Int63(), Transform: m1, Material: NewDefaultMaterial()}
}

type Sphere struct {
	Id        int64
	Transform Mat4x4
	Material  Material
}

// SetTransform passes a pointer to the Sphere on which to apply the translation matrix
func SetTransform(s *Sphere, translation Mat4x4) {
	s.Transform = Multiply(s.Transform, translation)
}

// SetMaterial passes a pointer to the Sphere on which to set the material
func SetMaterial(s *Sphere, m Material) {
	s.Material = m
}
