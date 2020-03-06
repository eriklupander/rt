package mat

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewSphere() *Sphere {

	return &Sphere{
		Id:          rand.Int63(),
		Transform:   New4x4(),
		Inverse:     New4x4(),
		Material:    NewDefaultMaterial(),
		savedVec:    NewVector(0, 0, 0),
		savedNormal: NewVector(0, 0, 0),
		savedRay:    NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 0)),
		xsCache:     make([]Intersection, 2),
		xsEmpty:     make([]Intersection, 0),
		originPoint: NewPoint(0, 0, 0),
		CastShadow:  true,
		Label:       "Sphere",
	}
}

func NewGlassSphere() *Sphere {
	s := NewSphere()
	material := NewGlassMaterial(1.5)
	s.SetMaterial(material)
	return s
}

type Sphere struct {
	Id        int64
	Transform Mat4x4
	Inverse   Mat4x4
	Material  Material
	Label     string
	parent    Shape
	savedRay  Ray

	// cached stuff
	originPoint Tuple4
	savedVec    Tuple4
	xsCache     []Intersection
	xsEmpty     []Intersection

	savedNormal Tuple4

	CastShadow bool
}

func (s *Sphere) CastsShadow() bool {
	return s.CastShadow
}

func (s *Sphere) GetParent() Shape {
	return s.parent
}

func (s *Sphere) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	SubPtr(point, s.originPoint, &s.savedNormal)
	return s.savedNormal
}

func (s *Sphere) GetLocalRay() Ray {
	return s.savedRay
}

// IntersectLocal implements Sphere-ray intersection
func (s *Sphere) IntersectLocal(r Ray) []Intersection {
	s.savedRay = r
	//s.XsCache = s.XsCache[:0]
	// this is a vector from the origin of the ray to the center of the sphere at 0,0,0
	//SubPtr(r.Origin, s.originPoint, &s.savedVec)

	// Note that doing the Subtraction inlined was much faster than letting SubPtr do it.
	// Shouldn't the SubPtr be inlined by the compiler? Need to figure out what's going on here...
	for i := 0; i < 4; i++ {
		s.savedVec[i] = r.Origin[i] - s.originPoint[i]
	}

	// This dot product is
	a := Dot(r.Direction, r.Direction)

	// Take the dot of the direction and the vector from ray origin to sphere center times 2
	b := 2.0 * Dot(r.Direction, s.savedVec)

	// Take the dot of the two sphereToRay vectors and decrease by 1 (is that because the sphere is unit length 1?
	c := Dot(s.savedVec, s.savedVec) - 1.0

	// calculate the discriminant
	discriminant := (b * b) - 4*a*c
	if discriminant < 0.0 {
		return s.xsEmpty
	}

	// finally, find the intersection distances on our ray. Some values:
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	s.xsCache[0].T = t1
	s.xsCache[1].T = t2
	s.xsCache[0].S = s
	s.xsCache[1].S = s
	return s.xsCache
}

func (s *Sphere) ID() int64 {
	return s.Id
}
func (s *Sphere) GetTransform() Mat4x4 {
	return s.Transform
}
func (s *Sphere) GetInverse() Mat4x4 {
	return s.Inverse
}

func (s *Sphere) GetMaterial() Material {
	return s.Material
}

// SetTransform passes a pointer to the Sphere on which to apply the translation matrix
func (s *Sphere) SetTransform(translation Mat4x4) {
	s.Transform = Multiply(s.Transform, translation)
	s.Inverse = Inverse(s.Transform)
}

// SetMaterial passes a pointer to the Sphere on which to set the material
func (s *Sphere) SetMaterial(m Material) {
	s.Material = m
}

func (s *Sphere) SetParent(shape Shape) {
	s.parent = shape
}

func (s *Sphere) Name() string {
	return s.Label
}
