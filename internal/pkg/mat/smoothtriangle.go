package mat

import "math"

func DefaultSmoothTriangle() *SmoothTriangle {
	return NewSmoothTriangle(
		NewPoint(0, 1, 0),
		NewPoint(-1, 0, 0),
		NewPoint(1, 0, 0),
		NewVector(0, 1, 0),
		NewVector(-1, 0, 0),
		NewVector(1, 0, 0))
}

func NewSmoothTriangle(p1 Tuple4, p2 Tuple4, p3 Tuple4, n1 Tuple4, n2 Tuple4, n3 Tuple4) *SmoothTriangle {

	e1 := Sub(p2, p1)
	e2 := Sub(p3, p1)
	n := Normalize(Cross(e2, e1))
	return &SmoothTriangle{P1: p1, P2: p2, P3: p3, E1: e1, E2: e2, N: n, N1: n1, N2: n2, N3: n3, Material: NewDefaultMaterial()}
}

type SmoothTriangle struct {
	P1       Tuple4
	P2       Tuple4
	P3       Tuple4
	E1       Tuple4
	E2       Tuple4
	N        Tuple4
	N1       Tuple4
	N2       Tuple4
	N3       Tuple4
	Material Material
	Shadow bool
}

func (s *SmoothTriangle) ID() int64 {
	return -1
}

func (s *SmoothTriangle) GetTransform() Mat4x4 {
	return IdentityMatrix
}

func (s *SmoothTriangle) GetInverse() Mat4x4 {
	elems := make([]float64, 16)
	copy(elems, IdentityMatrix.Elems)
	return Mat4x4{Elems: elems}
}

func (s *SmoothTriangle) SetTransform(transform Mat4x4) {
	panic("implement me")
}

func (s *SmoothTriangle) GetMaterial() Material {
	return s.Material
}

func (s *SmoothTriangle) SetMaterial(material Material) {
	s.Material = material
}

func (s *SmoothTriangle) IntersectLocal(ray Ray) []Intersection {
	dirCrossE2 := Cross(ray.Direction, s.E2)
	determinant := Dot(s.E1, dirCrossE2)
	if math.Abs(determinant) < Epsilon {
		return []Intersection{}
	}

	// Triangle misses over P1-P3 edge
	f := 1.0 / determinant
	p1ToOrigin := Sub(ray.Origin, s.P1)
	u := f * Dot(p1ToOrigin, dirCrossE2)
	if u < 0 || u > 1 {
		return []Intersection{}
	}

	originCrossE1 := Cross(p1ToOrigin, s.E1)
	v := f * Dot(ray.Direction, originCrossE1)
	if v < 0 || (u+v) > 1 {
		return []Intersection{}
	}
	tdist := f * Dot(s.E2, originCrossE1)
	return []Intersection{NewIntersectionUV(tdist, s, u, v)}
}

func (s *SmoothTriangle) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	return Add(Add(MultiplyByScalar(s.N2, intersection.U),
		MultiplyByScalar(s.N3, intersection.V)),
		MultiplyByScalar(s.N1, 1-intersection.U-intersection.V))
}

func (s *SmoothTriangle) GetLocalRay() Ray {
	panic("implement me")
}

func (s *SmoothTriangle) GetParent() Shape {
	return nil
}

func (s *SmoothTriangle) SetParent(shape Shape) {

}
func (s *SmoothTriangle) CastShadow()  bool {
	return s.Shadow
}