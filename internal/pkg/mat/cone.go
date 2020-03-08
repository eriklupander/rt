package mat

import (
	"math"
	"math/rand"
)

func NewCone() *Cone {
	m1 := New4x4()  //NewMat4x4(make([]float64, 16))
	inv := New4x4() //NewMat4x4(make([]float64, 16))
	return &Cone{
		Id:         rand.Int63(),
		Transform:  m1,
		Inverse:    inv,
		Material:   NewDefaultMaterial(),
		MinY:       math.Inf(-1),
		MaxY:       math.Inf(1),
		CastShadow: true,
	}
}

func NewConeMMC(min, max float64, closed bool) *Cone {
	c := NewCone()
	c.MinY = min
	c.MaxY = max
	c.Closed = closed
	return c
}

type Cone struct {
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
	Closed           bool
	CastShadow       bool
}

func (c *Cone) CastsShadow() bool {
	return c.CastShadow
}

func (c *Cone) ID() int64 {
	return c.Id
}

func (c *Cone) GetTransform() Mat4x4 {
	return c.Transform
}
func (c *Cone) GetInverse() Mat4x4 {
	return c.Inverse
}
func (c *Cone) GetInverseTranspose() Mat4x4 {
	return c.InverseTranspose
}

func (c *Cone) SetTransform(transform Mat4x4) {
	c.Transform = Multiply(c.Transform, transform)
	c.Inverse = Inverse(c.Transform)
	c.InverseTranspose = Transpose(c.Inverse)
}

func (c *Cone) GetMaterial() Material {
	return c.Material
}

func (c *Cone) SetMaterial(material Material) {
	c.Material = material
}

func (c *Cone) IntersectLocal(ray Ray) []Intersection {
	var xs []Intersection
	rdx2 := ray.Direction.Get(0) * ray.Direction.Get(0)
	rdy2 := ray.Direction.Get(1) * ray.Direction.Get(1)
	rdz2 := ray.Direction.Get(2) * ray.Direction.Get(2)

	a := rdx2 - rdy2 + rdz2

	b := 2*ray.Origin.Get(0)*ray.Direction.Get(0) -
		2*ray.Origin.Get(1)*ray.Direction.Get(1) +
		2*ray.Origin.Get(2)*ray.Direction.Get(2)

	absA := math.Abs(a)
	absB := math.Abs(b)
	if absA < Epsilon && absB < Epsilon {
		return xs
	}

	rox2 := ray.Origin.Get(0) * ray.Origin.Get(0)
	roy2 := ray.Origin.Get(1) * ray.Origin.Get(1)
	roz2 := ray.Origin.Get(2) * ray.Origin.Get(2)

	c1 := rox2 - roy2 + roz2

	//if math.Abs(a) < Epsilon {
	//	return c.intercectCaps(ray, xs)
	//}

	disc := b*b - 4*a*c1

	// ray does not intersect the cone
	if disc < 0 {
		return xs
	}
	var t0, t1 float64
	if absA < Epsilon && absB > Epsilon {
		t0 = -c1 / (2.0 * b)
		y0 := ray.Origin.Get(1) + t0*ray.Direction.Get(1)
		if y0 > c.MinY && y0 < c.MaxY {
			xs = append(xs, NewIntersection(t0, c))
		}
		//t1 = -c1 / (2.0 * b)
	} else {
		t0 = (-b - math.Sqrt(disc)) / (2 * a)
		t1 = (-b + math.Sqrt(disc)) / (2 * a)

		// Capping check
		y0 := ray.Origin.Get(1) + t0*ray.Direction.Get(1)
		if y0 > c.MinY && y0 < c.MaxY {
			xs = append(xs, NewIntersection(t0, c))
		}

		y1 := ray.Origin.Get(1) + t1*ray.Direction.Get(1)
		if y1 > c.MinY && y1 < c.MaxY {
			xs = append(xs, NewIntersection(t1, c))
		}
	}

	// Lids on top and bottom
	return c.intercectCaps(ray, xs)
}

func (c *Cone) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {

	// compute the square of the distance from the y axis
	dist := math.Pow(point.Get(0), 2) + math.Pow(point.Get(2), 2)
	if dist < 1 && point.Get(1) >= c.MaxY-Epsilon {
		return NewVector(0, 1, 0)
	} else if dist < 1 && point.Get(1) <= c.MinY+Epsilon {
		return NewVector(0, -1, 0)
	} else {
		y := math.Sqrt(math.Pow(point.Get(0), 2) + math.Pow(point.Get(2), 2))
		if point.Get(1) > 0.0 {
			y = -y
		}
		return NewVector(point.Get(0), y, point.Get(2))
	}
}

func (c *Cone) GetLocalRay() Ray {
	return c.savedRay
}

// checkCap for cones changes so the MinY / MaxY is used instead of 1.0 since the cone narrows down.
// (remember, we're in unit space)
func (c *Cone) checkCap(ray Ray, t float64, minMaxY float64) bool {
	x := ray.Origin.Get(0) + t*ray.Direction.Get(0)
	z := ray.Origin.Get(2) + t*ray.Direction.Get(2)
	return math.Pow(x, 2)+math.Pow(z, 2) <= math.Abs(minMaxY)
}

func (c *Cone) intercectCaps(ray Ray, xs []Intersection) []Intersection {
	if !c.Closed || math.Abs(ray.Direction.Get(1)) < Epsilon {
		return xs
	}

	// check for an intersection with the lower end cap by intersecting
	// the ray with the plane at y=cyl.minimum
	t := (c.MinY - ray.Origin.Get(1)) / ray.Direction.Get(1)
	if c.checkCap(ray, t, c.MinY) {
		xs = append(xs, NewIntersection(t, c))
	}

	// check for an intersection with the upper end cap by intersecting
	// the ray with the plane at y=cyl.maximum
	t = (c.MaxY - ray.Origin.Get(1)) / ray.Direction.Get(1)
	if c.checkCap(ray, t, c.MaxY) {
		xs = append(xs, NewIntersection(t, c))
	}
	return xs
}

func (c *Cone) GetParent() Shape {
	return c.parent
}
func (c *Cone) SetParent(shape Shape) {
	c.parent = shape
}
func (c *Cone) Name() string {
	return c.Label
}
