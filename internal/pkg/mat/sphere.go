package mat

import (
	"math/rand"
)

func NewSphere() *Sphere {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	return &Sphere{Id: rand.Int63(), Transform: m1}
}

type Sphere struct {
	Id        int64
	Transform *Mat4x4
	Material  *Material
}

func SetTransform(s *Sphere, translation *Mat4x4) {
	s.Transform = Multiply(s.Transform, translation)
}
func SetMaterial(s *Sphere, m *Material) {
	s.Material = m
}
