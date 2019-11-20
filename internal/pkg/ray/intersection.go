package ray

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
)

type Intersection struct {
	T float64
	S mat.Sphere
}

func NewIntersection(t float64, s mat.Sphere) *Intersection {
	return &Intersection{T: t, S: s}
}

func IntersectionEqual(i1, i2 Intersection) bool {
	return i1.T == i2.T && i1.S.Id == i2.S.Id
}
