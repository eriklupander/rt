package mat

import (
	"math"
	"math/rand"
)

func NewCone() *Cone {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	return &Cone{Id: rand.Int63(), Transform: m1, Material: NewDefaultMaterial(), minY: math.Inf(-1), maxY: math.Inf(1)}
}

func NewConeMM(min, max float64) *Cone {
	c := NewCone()
	c.minY = min
	c.maxY = max
	return c
}

func NewConeMMC(min, max float64, closed bool) *Cone {
	c := NewCone()
	c.minY = min
	c.maxY = max
	c.closed = closed
	return c
}

type Cone struct {
	Id        int64
	Transform Mat4x4
	Material  Material
	Label     string
	Parent    Shape
	savedRay  Ray
	minY      float64
	maxY      float64
	closed    bool
}

func (c *Cone) ID() int64 {
	return c.Id
}

func (c *Cone) GetTransform() Mat4x4 {
	return c.Transform
}

func (c *Cone) SetTransform(transform Mat4x4) {
	c.Transform = Multiply(c.Transform, transform)
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
		if y0 > c.minY && y0 < c.maxY {
			xs = append(xs, NewIntersection(t0, c))
		}
		//t1 = -c1 / (2.0 * b)
	} else {
		t0 = (-b - math.Sqrt(disc)) / (2 * a)
		t1 = (-b + math.Sqrt(disc)) / (2 * a)

		// Capping check
		y0 := ray.Origin.Get(1) + t0*ray.Direction.Get(1)
		if y0 > c.minY && y0 < c.maxY {
			xs = append(xs, NewIntersection(t0, c))
		}

		y1 := ray.Origin.Get(1) + t1*ray.Direction.Get(1)
		if y1 > c.minY && y1 < c.maxY {
			xs = append(xs, NewIntersection(t1, c))
		}
	}

	// Lids on top and bottom
	return c.intercectCaps(ray, xs)
}

func (c *Cone) NormalAtLocal(point Tuple4) Tuple4 {

	// compute the square of the distance from the y axis
	dist := math.Pow(point.Get(0), 2) + math.Pow(point.Get(2), 2)
	if dist < 1 && point.Get(1) >= c.maxY-Epsilon {
		return NewVector(0, 1, 0)
	} else if dist < 1 && point.Get(1) <= c.minY+Epsilon {
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

// checkCap for cones changes so the minY / maxY is used instead of 1.0 since the cone narrows down.
// (remember, we're in unit space)
func (c *Cone) checkCap(ray Ray, t float64, minMaxY float64) bool {
	x := ray.Origin.Get(0) + t*ray.Direction.Get(0)
	z := ray.Origin.Get(2) + t*ray.Direction.Get(2)
	return math.Pow(x, 2)+math.Pow(z, 2) <= math.Abs(minMaxY)
}

func (c *Cone) intercectCaps(ray Ray, xs []Intersection) []Intersection {
	if !c.closed || math.Abs(ray.Direction.Get(1)) < Epsilon {
		return xs
	}

	// check for an intersection with the lower end cap by intersecting
	// the ray with the plane at y=cyl.minimum
	t := (c.minY - ray.Origin.Get(1)) / ray.Direction.Get(1)
	if c.checkCap(ray, t, c.minY) {
		xs = append(xs, NewIntersection(t, c))
	}

	// check for an intersection with the upper end cap by intersecting
	// the ray with the plane at y=cyl.maximum
	t = (c.maxY - ray.Origin.Get(1)) / ray.Direction.Get(1)
	if c.checkCap(ray, t, c.maxY) {
		xs = append(xs, NewIntersection(t, c))
	}
	return xs
}

func (c *Cone) GetParent() Shape {
	return c.Parent
}
func (c *Cone) SetParent(shape Shape) {
	c.Parent = shape
}
