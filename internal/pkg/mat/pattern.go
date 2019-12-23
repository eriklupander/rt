package mat

import (
	"math"
)

var black = NewColor(0, 0, 0)
var white = NewColor(1, 1, 1)

type Pattern interface {
	PatternAt(point Tuple4) Tuple4
	SetPatternTransform(transform Mat4x4)
	GetTransform() Mat4x4
	GetInverse() Mat4x4
}

func NewStripePattern(colorA Tuple4, colorB Tuple4) *StripePattern {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	inv := NewMat4x4(make([]float64, 16))
	copy(inv.Elems, IdentityMatrix.Elems)
	return &StripePattern{A: colorA, B: colorB, Transform: m1, Inverse: inv}
}

type TestPattern struct {
	Transform Mat4x4
	Inverse   Mat4x4
}

func NewTestPattern() *TestPattern {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	inv := NewMat4x4(make([]float64, 16))
	copy(inv.Elems, IdentityMatrix.Elems)
	return &TestPattern{Transform: m1, Inverse: inv}
}

func (t *TestPattern) PatternAt(point Tuple4) Tuple4 {
	return NewColor(point.Get(0), point.Get(1), point.Get(2))
}

func (t *TestPattern) SetPatternTransform(transform Mat4x4) {
	t.Transform = transform
	t.Inverse = Inverse(t.Transform)
}

func (t *TestPattern) GetTransform() Mat4x4 {
	return t.Transform
}
func (t *TestPattern) GetInverse() Mat4x4 {
	return t.Inverse
}

type StripePattern struct {
	A         Tuple4
	B         Tuple4
	Transform Mat4x4
	Inverse   Mat4x4
}

func (p *StripePattern) GetTransform() Mat4x4 {
	return p.Transform
}
func (p *StripePattern) GetInverse() Mat4x4 {
	return p.Inverse
}
func (p *StripePattern) SetPatternTransform(transform Mat4x4) {
	p.Transform = transform
	p.Inverse = Inverse(p.Transform)
}

func (p *StripePattern) PatternAt(point Tuple4) Tuple4 {
	if int(math.Floor(point.Get(0)))%2 == 0 {
		return p.A
	}
	return p.B
}

type GradientPattern struct {
	FromColor Tuple4
	ToColor   Tuple4
	Transform Mat4x4
	Inverse   Mat4x4
}

func NewGradientPattern(from, to Tuple4) *GradientPattern {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	inv := NewMat4x4(make([]float64, 16))
	copy(inv.Elems, IdentityMatrix.Elems)
	return &GradientPattern{FromColor: from, ToColor: to, Transform: m1, Inverse: inv}
}

func (g *GradientPattern) PatternAt(point Tuple4) Tuple4 {
	distance := Sub(g.ToColor, g.FromColor)
	fraction := point.Get(0) - math.Floor(point.Get(0))
	return Add(g.FromColor, MultiplyByScalar(distance, fraction))
}

func (g *GradientPattern) SetPatternTransform(transform Mat4x4) {
	g.Transform = transform
	g.Inverse = Inverse(g.Transform)
}

func (g *GradientPattern) GetTransform() Mat4x4 {
	return g.Transform
}

func (g *GradientPattern) GetInverse() Mat4x4 {
	return g.Inverse
}

type RingPattern struct {
	A         Tuple4
	B         Tuple4
	Transform Mat4x4
	Inverse   Mat4x4
}

func NewRingPattern(a Tuple4, b Tuple4) *RingPattern {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	inv := NewMat4x4(make([]float64, 16))
	copy(inv.Elems, IdentityMatrix.Elems)
	return &RingPattern{A: a, B: b, Transform: m1, Inverse: inv}
}

func (r *RingPattern) PatternAt(point Tuple4) Tuple4 {
	flooredDistance := math.Floor(math.Sqrt(point.Get(0)*point.Get(0) + point.Get(2)*point.Get(2)))
	if int(flooredDistance)%2 == 0 {
		return r.A
	}
	return r.B
}

func (r *RingPattern) SetPatternTransform(transform Mat4x4) {
	r.Transform = transform
	r.Inverse = Inverse(r.Transform)
}

func (r *RingPattern) GetTransform() Mat4x4 {
	return r.Transform
}

func (r *RingPattern) GetInverse() Mat4x4 {
	return r.Inverse
}

func PatternAtShape(pattern Pattern, s Shape, worldPoint Tuple4) Tuple4 {
	// Convert from world space to object space by inversing the shape transform and then multiply it by the point
	//objectPoint := MultiplyByTuple(Inverse(s.GetTransform()), worldPoint)
	objectPoint := WorldToObject(s, worldPoint)
	//patternPoint := MultiplyByTuple(Inverse(pattern.GetTransform()), objectPoint)
	patternPoint := MultiplyByTuple(pattern.GetInverse(), objectPoint)
	return pattern.PatternAt(patternPoint)
}

// use this new
//world_to_object() function when converting points from world space to object
//space.

func NewCheckerPattern(colorA Tuple4, colorB Tuple4) *CheckerPattern {
	m1 := NewMat4x4(make([]float64, 16))
	copy(m1.Elems, IdentityMatrix.Elems)
	inv := NewMat4x4(make([]float64, 16))
	copy(inv.Elems, IdentityMatrix.Elems)
	return &CheckerPattern{ColorA: colorA, ColorB: colorB, Transform: m1, Inverse: inv}
}

type CheckerPattern struct {
	ColorA    Tuple4
	ColorB    Tuple4
	Transform Mat4x4
	Inverse   Mat4x4
}

func (c *CheckerPattern) PatternAt(point Tuple4) Tuple4 {
	if (int(math.Floor(point.Get(0)))+int(math.Floor(point.Get(1)))+int(math.Floor(point.Get(2))))%2 == 0 {
		return c.ColorA
	}
	return c.ColorB
}

func (c *CheckerPattern) SetPatternTransform(transform Mat4x4) {
	c.Transform = transform
	c.Inverse = Inverse(c.Transform)
}

func (c *CheckerPattern) GetTransform() Mat4x4 {
	return c.Transform
}
func (c *CheckerPattern) GetInverse() Mat4x4 {
	return Inverse(c.Transform)//c.Inverse
}
