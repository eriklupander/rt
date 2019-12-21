package mat

import "math"

type Triangle struct {
	P1       Tuple4
	P2       Tuple4
	P3       Tuple4
	E1       Tuple4
	E2       Tuple4
	N        Tuple4
	Material Material
}

func (t *Triangle) ID() int64 {
	panic("implement me")
}

func (t *Triangle) GetTransform() Mat4x4 {
	return IdentityMatrix
}

func (t *Triangle) GetInverse() Mat4x4 {
	panic("implement me")
}

func (t *Triangle) SetTransform(transform Mat4x4) {
	panic("implement me")
}

func (t *Triangle) GetMaterial() Material {
	return t.Material
}

func (t *Triangle) SetMaterial(material Material) {
	panic("implement me")
}

func (t *Triangle) IntersectLocal(ray Ray) []Intersection {
	dirCrossE2 := Cross(ray.Direction, t.E2)
	determinant := Dot(t.E1, dirCrossE2)
	if math.Abs(determinant) < Epsilon {
		return []Intersection{}
	}

	// Triangle misses over P1-P3 edge
	f := 1.0 / determinant
	p1ToOrigin := Sub(ray.Origin, t.P1)
	u := f * Dot(p1ToOrigin, dirCrossE2)
	if u < 0 || u > 1 {
		return []Intersection{}
	}

	originCrossE1 := Cross(p1ToOrigin, t.E1)
	v := f * Dot(ray.Direction, originCrossE1)
	if v < 0 || (u+v) > 1 {
		return []Intersection{}
	}
	tdist := f * Dot(t.E2, originCrossE1)
	return []Intersection{NewIntersection(tdist, t)}
}

func (t *Triangle) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	return t.N
}

func (t *Triangle) GetLocalRay() Ray {
	panic("implement me")
}

func (t *Triangle) GetParent() Shape {
	return nil
}

func (t *Triangle) SetParent(shape Shape) {

}

func NewTriangle(p1 Tuple4, p2 Tuple4, p3 Tuple4) *Triangle {

	e1 := Sub(p2, p1)
	e2 := Sub(p3, p1)
	n := Normalize(Cross(e2, e1))
	return &Triangle{P1: p1, P2: p2, P3: p3, E1: e1, E2: e2, N: n, Material: NewDefaultMaterial()}
}
