package mat

import (
	"math"
	"math/rand"
)

func NewCylinder() *Cylinder {
	m1 := New4x4()  //NewMat4x4(make([]float64, 16))
	inv := New4x4() //NewMat4x4(make([]float64, 16))
	invTranspose := New4x4()

	savedXs := make([]Intersection, 4)

	return &Cylinder{
		Id:               rand.Int63(),
		Transform:        m1,
		Inverse:          inv,
		InverseTranspose: invTranspose,
		Material:         NewDefaultMaterial(),
		MinY:             math.Inf(-1),
		MaxY:             math.Inf(1),
		savedXs:          savedXs,
		CastShadow:       true,
	}
}

func NewCylinderMM(min, max float64) *Cylinder {
	c := NewCylinder()
	c.MinY = min
	c.MaxY = max
	return c
}

func NewCylinderMMC(min, max float64, closed bool) *Cylinder {
	c := NewCylinder()
	c.MinY = min
	c.MaxY = max
	c.closed = closed
	return c
}

type Cylinder struct {
	Id               int64
	Transform        Mat4x4
	Inverse          Mat4x4
	InverseTranspose Mat4x4
	Material         Material
	Label            string
	parent           Shape
	savedRay         Ray
	MinY             float64
	MaxY             float64
	closed           bool
	CastShadow       bool

	savedXs []Intersection
}

func (c *Cylinder) CastsShadow() bool {
	return c.CastShadow
}

func (c *Cylinder) ID() int64 {
	return c.Id
}

func (c *Cylinder) GetTransform() Mat4x4 {
	return c.Transform
}

func (c *Cylinder) GetInverse() Mat4x4 {
	return c.Inverse
}
func (c *Cylinder) GetInverseTranspose() Mat4x4 {
	return c.InverseTranspose
}

func (c *Cylinder) SetTransform(transform Mat4x4) {
	c.Transform = Multiply(c.Transform, transform)
	c.Inverse = Inverse(c.Transform)
	c.InverseTranspose = Transpose(c.Inverse)
}

func (c *Cylinder) GetMaterial() Material {
	return c.Material
}

func (c *Cylinder) SetMaterial(material Material) {
	c.Material = material
}

func (c *Cylinder) IntersectLocal(ray Ray) []Intersection {
	//var xs []Intersection
	rdx2 := ray.Direction.Get(0) * ray.Direction.Get(0)
	rdz2 := ray.Direction.Get(2) * ray.Direction.Get(2)

	a := rdx2 + rdz2
	c.savedXs = c.savedXs[:0]
	if math.Abs(a) < Epsilon {
		return c.intercectCaps(ray, c.savedXs)
	}

	b := 2*ray.Origin.Get(0)*ray.Direction.Get(0) +
		2*ray.Origin.Get(2)*ray.Direction.Get(2)

	rox2 := ray.Origin.Get(0) * ray.Origin.Get(0)
	roz2 := ray.Origin.Get(2) * ray.Origin.Get(2)

	c1 := rox2 + roz2 - 1

	disc := b*b - 4*a*c1

	// ray does not intersect the cylinder
	if disc < 0 {
		return c.savedXs
	}

	t0 := (-b - math.Sqrt(disc)) / (2 * a)
	t1 := (-b + math.Sqrt(disc)) / (2 * a)

	y0 := ray.Origin.Get(1) + t0*ray.Direction.Get(1)
	if y0 > c.MinY && y0 < c.MaxY {
		c.savedXs = append(c.savedXs, NewIntersection(t0, c))
	}

	y1 := ray.Origin.Get(1) + t1*ray.Direction.Get(1)
	if y1 > c.MinY && y1 < c.MaxY {
		c.savedXs = append(c.savedXs, NewIntersection(t1, c))
	}

	return c.intercectCaps(ray, c.savedXs)
}

func (c *Cylinder) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {

	// compute the square of the distance from the y axis
	dist := math.Pow(point.Get(0), 2) + math.Pow(point.Get(2), 2)
	if dist < 1 && point.Get(1) >= c.MaxY-Epsilon {
		return NewVector(0, 1, 0)
	} else if dist < 1 && point.Get(1) <= c.MinY+Epsilon {
		return NewVector(0, -1, 0)
	} else {
		return NewVector(point.Get(0), 0, point.Get(2))
	}
}

func (c *Cylinder) GetLocalRay() Ray {
	return c.savedRay
}
func (c *Cylinder) GetParent() Shape {
	return c.parent
}
func (c *Cylinder) SetParent(shape Shape) {
	c.parent = shape
}

func (c *Cylinder) Init() {
	c.savedXs = make([]Intersection, 4)
}

func checkCap(ray Ray, t float64) bool {
	x := ray.Origin.Get(0) + t*ray.Direction.Get(0)
	z := ray.Origin.Get(2) + t*ray.Direction.Get(2)
	return math.Pow(x, 2)+math.Pow(z, 2) <= 1.0
}

func (c *Cylinder) intercectCaps(ray Ray, xs []Intersection) []Intersection {
	if !c.closed || math.Abs(ray.Direction.Get(1)) < Epsilon {
		return xs
	}

	// check for an intersection with the lower end cap by intersecting
	// the ray with the plane at y=cyl.minimum
	t := (c.MinY - ray.Origin.Get(1)) / ray.Direction.Get(1)
	if checkCap(ray, t) {
		xs = append(xs, NewIntersection(t, c))
	}

	// check for an intersection with the upper end cap by intersecting
	// the ray with the plane at y=cyl.maximum
	t = (c.MaxY - ray.Origin.Get(1)) / ray.Direction.Get(1)
	if checkCap(ray, t) {
		xs = append(xs, NewIntersection(t, c))
	}
	return xs
}
func (c *Cylinder) Name() string {
	return c.Label
}
