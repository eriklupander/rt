package mat

import (
	"math"
	"math/rand"
)

func NewCube() *Cube {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	return &Cube{Id: rand.Int63(), Transform: m1, Material: NewDefaultMaterial()}
}

type Cube struct {
	Id        int64
	Transform Mat4x4
	Material  Material
	Label     string
	Parent    Shape
	savedRay  Ray
}

func (c *Cube) ID() int64 {
	return c.Id
}

func (c *Cube) GetTransform() Mat4x4 {
	return c.Transform
}

func (c *Cube) SetTransform(transform Mat4x4) {
	c.Transform = Multiply(c.Transform, transform)
}

func (c *Cube) GetMaterial() Material {
	return c.Material
}

func (c *Cube) SetMaterial(material Material) {
	c.Material = material
}

func (c *Cube) IntersectLocal(ray Ray) []Intersection {
	// There is supposed  to be a way to optimize this for fewer checks by looking at early values.
	xtmin, xtmax := checkAxis(ray.Origin.Get(0), ray.Direction.Get(0))
	ytmin, ytmax := checkAxis(ray.Origin.Get(1), ray.Direction.Get(1))
	ztmin, ztmax := checkAxis(ray.Origin.Get(2), ray.Direction.Get(2))

	// Om det största av min-värdena är större än det minsta max-värdet.
	tmin := max(xtmin, ytmin, ztmin)
	tmax := min(xtmax, ytmax, ztmax)
	if tmin > tmax {
		return []Intersection{}
	}
	return []Intersection{NewIntersection(tmin, c), NewIntersection(tmax, c)}
}

// NormalAtLocal uses the fact that given a unit cube, the point of the surface axis X,Y or Z is always either
// 1.0 for positive XYZ and -1.0 for negative XYZ. I.e - if the point is 0.4, 1, -0.5, we know that the
// point is on the top Y surface and we can return a 0,1,0 normal
func (c *Cube) NormalAtLocal(point Tuple4, intersection *Intersection) Tuple4 {
	maxc := max(math.Abs(point.Get(0)), math.Abs(point.Get(1)), math.Abs(point.Get(2)))
	if maxc == math.Abs(point.Get(0)) {
		return NewVector(point.Get(0), 0, 0)
	} else if maxc == math.Abs(point.Get(1)) {
		return NewVector(0, point.Get(1), 0)
	}
	return NewVector(0, 0, point.Get(2))
}

func (c *Cube) GetLocalRay() Ray {
	return c.savedRay
}
func (c *Cube) GetParent() Shape {
	return c.Parent
}
func (c *Cube) SetParent(shape Shape) {
	c.Parent = shape
}

func checkAxis(origin float64, direction float64) (min float64, max float64) {
	tminNumerator := -1 - origin
	tmaxNumerator := 1 - origin
	var tmin, tmax float64
	if math.Abs(direction) >= Epsilon {
		tmin = tminNumerator / direction
		tmax = tmaxNumerator / direction
	} else {
		tmin = tminNumerator * math.Inf(1)
		tmax = tmaxNumerator * math.Inf(1)
	}
	if tmin > tmax {
		// swap
		temp := tmin
		tmin = tmax
		tmax = temp
	}
	return tmin, tmax
}

func max(values ...float64) float64 {
	c := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > c {
			c = values[i]
		}
	}
	return c
}

func min(values ...float64) float64 {
	c := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < c {
			c = values[i]
		}
	}
	return c
}
