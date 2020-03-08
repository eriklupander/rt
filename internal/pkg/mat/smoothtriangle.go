package mat

import (
	"math"
)

const TriThreshold = 0.00000000001

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

	// for barycentric
	//d00 := Dot(e1, e1)
	//d01 := Dot(e1, e2)
	//d11 := Dot(e2, e2)
	//denom := d00*d11 - d01*d01

	return &SmoothTriangle{P1: p1, P2: p2, P3: p3, E1: e1, E2: e2, N: n, N1: n1, N2: n2, N3: n3,
		Material:   NewDefaultMaterial(),
		CastShadow: true,
		Label:      "SmoothTriangle",
		//D00:d00,
		//D01:d01,
		//D11:d11,
		//Denom: denom,
	}
}

type SmoothTriangle struct {
	P1         Tuple4
	P2         Tuple4
	P3         Tuple4
	E1         Tuple4
	E2         Tuple4
	N          Tuple4
	N1         Tuple4
	N2         Tuple4
	N3         Tuple4
	Material   Material
	CastShadow bool
	Label      string

	D00   float64
	D01   float64
	D11   float64
	Denom float64

	parent Shape
}

// Barycentric computes barycentric coordinates (u, v, w) for point p with respect to triangle defined by pre-computed
// vectors E1 and E2, which was derived into points d00, d01, d11 and denominator in constructor func.
func (s *SmoothTriangle) Barycentric(p Tuple4, u *float64, v *float64, w *float64) {

	v2 := NewTuple()
	SubPtr(p, s.P1, &v2)

	d20 := Dot(v2, s.E1)
	d21 := Dot(v2, s.E2)

	*v = (s.D11*d20 - s.D01*d21) / s.Denom
	*w = (s.D00*d21 - s.D01*d20) / s.Denom
	*u = 1.0 - *v - *w
}

func (s *SmoothTriangle) CastsShadow() bool {
	return s.CastShadow
}

func (s *SmoothTriangle) ID() int64 {
	return -1
}

func (s *SmoothTriangle) GetTransform() Mat4x4 {
	return IdentityMatrix
}

func (s *SmoothTriangle) GetInverse() Mat4x4 {
	return IdentityMatrix
}
func (s *SmoothTriangle) GetInverseTranspose() Mat4x4 {
	return IdentityMatrix
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
	//fmt.Printf("intersecting triangle with material: " + s.Material.Name + " having color: %v\n", s.Material.Color)

	dirCrossE2 := Cross(ray.Direction, s.E2)
	determinant := Dot(s.E1, dirCrossE2)
	if math.Abs(determinant) < TriThreshold {
		return nil
	}

	// Triangle misses over P1-P3 edge
	f := 1.0 / determinant
	p1ToOrigin := Sub(ray.Origin, s.P1)
	u := f * Dot(p1ToOrigin, dirCrossE2)
	if u < 0 || u > 1 {
		return nil
	}

	originCrossE1 := Cross(p1ToOrigin, s.E1)
	v := f * Dot(ray.Direction, originCrossE1)
	if v < 0 || (u+v) > 1 {
		return nil
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
	return s.parent
}

func (s *SmoothTriangle) SetParent(shape Shape) {
	s.parent = shape
}
func (s *SmoothTriangle) Name() string {
	return s.Label
}
