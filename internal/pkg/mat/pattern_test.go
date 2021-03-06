package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStripePatternAtPoint(t *testing.T) {
	pattern := NewStripePattern(white, black)
	assert.Equal(t, pattern.A, white)
	assert.Equal(t, pattern.B, black)
}

func TestStripeAtConstantY(t *testing.T) {
	pattern := NewStripePattern(white, black)
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 1, 0)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 2, 0)))
}
func TestStripeAtConstantZ(t *testing.T) {
	pattern := NewStripePattern(white, black)
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 1)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 2)))
}
func TestStripeAtAlternateX(t *testing.T) {
	pattern := NewStripePattern(white, black)
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0.9, 0, 0)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(1, 0, 0)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(-0.1, 0, 0)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(-1, 0, 0)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(-1.1, 0, 0)))
}

func TestLightingWithPattern(t *testing.T) {
	s := NewSphere()
	pattern := NewStripePattern(NewColor(1, 1, 1), NewColor(0, 0, 0))
	material := Material{
		Ambient:  1,
		Diffuse:  0,
		Specular: 0,
		Pattern:  pattern,
	}
	eyeVec := NewVector(0, 0, -1)
	normalVec := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 0, -10), NewColor(1, 1, 1))

	c1 := LightingPointLight(material, s, light, NewPoint(0.9, 0, 0), eyeVec, normalVec, false, NewLightData())
	c2 := LightingPointLight(material, s, light, NewPoint(1.1, 0, 0), eyeVec, normalVec, false, NewLightData())
	assert.True(t, TupleXYZEq(white, c1))
	assert.True(t, TupleXYZEq(black, c2))
}

func TestStripeWithObjectTransform(t *testing.T) {
	s := NewSphere()
	s.SetTransform(Scale(2, 2, 2))
	s.Material.Pattern = NewStripePattern(white, black)
	color := PatternAtShape(s.Material.Pattern, s, NewPoint(1.5, 0, 0))
	assert.Equal(t, white, color)
}

func TestStripeWithPatternTransform(t *testing.T) {
	s := NewSphere()
	s.Material.Pattern = NewStripePattern(white, black)
	s.Material.Pattern.SetPatternTransform(Scale(2, 2, 2))
	color := PatternAtShape(s.Material.Pattern, s, NewPoint(1.5, 0, 0))
	assert.Equal(t, white, color)
}

func TestStripeWithBothObjectAndPatternTransform(t *testing.T) {
	s := NewSphere()
	s.SetTransform(Scale(2, 2, 2))
	s.Material.Pattern = NewStripePattern(white, black)
	s.Material.Pattern.SetPatternTransform(Translate(0.5, 0, 0))
	color := s.Material.Pattern.PatternAt(NewPoint(2.5, 0, 0))
	assert.Equal(t, white, color)
}

func TestStripeWithBothObjectAndPatternTransform2(t *testing.T) {
	s := NewSphere()
	s.SetTransform(Scale(2, 2, 2))
	s.Material.Pattern = NewTestPattern()
	s.Material.Pattern.SetPatternTransform(Translate(0.5, 1, 1.5))
	color := PatternAtShape(s.Material.Pattern, s, NewPoint(2.5, 3, 3.5))
	assert.Equal(t, NewColor(0.75, 0.5, 0.25), color)
}

func TestGradientPattern(t *testing.T) {
	pattern := NewGradientPattern(white, black)
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0)))
	assert.True(t, TupleXYZEq(NewColor(0.25, 0.25, 0.25), pattern.PatternAt(NewPoint(0.75, 0, 0))))
	assert.True(t, TupleXYZEq(NewColor(0.50, 0.5, 0.5), pattern.PatternAt(NewPoint(0.50, 0, 0))))
	assert.True(t, TupleXYZEq(NewColor(0.75, 0.75, 0.75), pattern.PatternAt(NewPoint(0.25, 0, 0))))
}
func TestRingPattern(t *testing.T) {

	pattern := NewRingPattern(white, black)
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(1, 0, 0)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(0, 0, 1)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(0.708, 0, 0.708)))
}
func TestCheckerPatternRepeatInX(t *testing.T) {
	pattern := NewCheckerPattern(white, black)
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0.99, 0, 0)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(1.01, 0, 0)))
}
func TestCheckerPatternRepeatInY(t *testing.T) {
	pattern := NewCheckerPattern(white, black)
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0.99, 0)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(0, 1.01, 0)))
}
func TestCheckerPatternRepeatInZ(t *testing.T) {
	pattern := NewCheckerPattern(white, black)
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0)))
	assert.Equal(t, white, pattern.PatternAt(NewPoint(0, 0, 0.99)))
	assert.Equal(t, black, pattern.PatternAt(NewPoint(0, 0, 1.01)))
}
