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
	m1 := NewMat4x4(make([]float64, 16))
	inv := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	copy(inv.Elems, IdentityMatrix.Elems)
	return &Sphere{Id: rand.Int63(), Transform: m1, Inverse: inv, Material: NewDefaultMaterial()}
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
	Parent    Shape
	savedRay  Ray
}

func (s *Sphere) GetParent() Shape {
	return s.Parent
}

func (s *Sphere) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	return Sub(point, NewPoint(0, 0, 0))
}

func (s *Sphere) GetLocalRay() Ray {
	return s.savedRay
}

// IntersectLocal implements Sphere-ray intersection
func (s *Sphere) IntersectLocal(r Ray) []Intersection {
	s.savedRay = r
	// this is a vector from the origin of the ray to the center of the sphere at 0,0,0
	sphereToRay := Sub(r.Origin, NewPoint(0, 0, 0))

	// This dot product is
	a := Dot(r.Direction, r.Direction)

	// Take the dot of the direction and the vector from ray origin to sphere center times 2
	b := 2.0 * Dot(r.Direction, sphereToRay)

	// Take the dot of the two sphereToRay vectors and decrease by 1 (is that because the sphere is unit length 1?
	c := Dot(sphereToRay, sphereToRay) - 1.0

	// calculate the discriminant
	discriminant := (b * b) - 4*a*c
	if discriminant < 0.0 {
		return []Intersection{}
	}

	// finally, find the intersection distances on our ray. Some values:
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	return []Intersection{
		{T: t1, S: s},
		{T: t2, S: s},
	}
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
	s.Parent = shape
}
