package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var rt2 = math.Sqrt(2.0) / 2.0

func TestNewLight(t *testing.T) {
	p := NewPoint(0, 0, 0)
	color := NewColor(1, 1, 1)
	light := NewLight(p, color)
	assert.True(t, TupleEquals(light.Position, p))
	assert.True(t, TupleEquals(light.Intensity, color))
}

func setupBase() (Shape, Material, Tuple4) {
	return NewSphere(), NewDefaultMaterial(), NewPoint(0, 0, 0)
}

func TestLightEyeBetweenLightAndSphere(t *testing.T) {
	s, material, position := setupBase()
	eyev := NewVector(0, 0, -1)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 0, -10), NewColor(1, 1, 1))

	result := LightingPointLight(material, s, light, position, eyev, normalv, false, NewLightData())
	assert.InEpsilon(t, 1.9, result.Get(0), Epsilon)
	assert.InEpsilon(t, 1.9, result.Get(1), Epsilon)
	assert.InEpsilon(t, 1.9, result.Get(2), Epsilon)
}
func TestLight180ToSurfaceEye45(t *testing.T) {
	s, material, position := setupBase()
	eyev := NewVector(0, 0, -1)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 10, -10), NewColor(1, 1, 1))

	result := LightingPointLight(material, s, light, position, eyev, normalv, false, NewLightData())
	assert.InEpsilon(t, 0.7364, result.Get(0), Epsilon)
	assert.InEpsilon(t, 0.7364, result.Get(1), Epsilon)
	assert.InEpsilon(t, 0.7364, result.Get(2), Epsilon)
}
func TestLight45ToSurfaceEye180(t *testing.T) {
	s, material, position := setupBase()
	eyev := NewVector(0, rt2, -rt2)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 0, -10), NewColor(1, 1, 1))

	result := LightingPointLight(material, s, light, position, eyev, normalv, false, NewLightData())
	assert.Equal(t, 1.0, result.Get(0))
	assert.Equal(t, 1.0, result.Get(1))
	assert.Equal(t, 1.0, result.Get(2))
}
func TestLight45ToSurfaceEye45(t *testing.T) {
	s, material, position := setupBase()
	eyev := NewVector(0, -rt2, -rt2)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 10, -10), NewColor(1, 1, 1))

	result := LightingPointLight(material, s, light, position, eyev, normalv, false, NewLightData())
	assert.InEpsilon(t, 1.6364, result.Get(0), Epsilon)
	assert.InEpsilon(t, 1.6364, result.Get(1), Epsilon)
	assert.InEpsilon(t, 1.6364, result.Get(2), Epsilon)
}
func TestLightBehind(t *testing.T) {
	s, material, position := setupBase()
	eyev := NewVector(0, 0, -1)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 0, 10), NewColor(1, 1, 1))

	result := LightingPointLight(material, s, light, position, eyev, normalv, false, NewLightData())
	assert.Equal(t, 0.1, result.Get(0))
	assert.Equal(t, 0.1, result.Get(1))
	assert.Equal(t, 0.1, result.Get(2))
}

func TestCreateAreaLight(t *testing.T) {

	corner := NewPoint(0, 0, 0)
	v1 := NewVector(2, 0, 0)
	v2 := NewVector(0, 0, 1)
	light := NewAreaLight(corner, v1, 4, v2, 2, NewColor(1, 1, 1))
	assert.Equal(t, corner, light.Corner)
	assert.Equal(t, NewVector(0.5, 0, 0), light.UVec)
	assert.Equal(t, 4, light.USteps)

	assert.Equal(t, NewVector(0, 0, 0.5), light.VVec)
	assert.Equal(t, 2, light.VSteps)

	assert.Equal(t, 8.0, light.Samples)
	assert.True(t, TupleXYZEq(NewPoint(1, 0, 0.5), light.Position))
}

func TestFindPointOnAreaLight(t *testing.T) {
	corner := NewPoint(0, 0, 0)
	v1 := NewVector(2, 0, 0)
	v2 := NewVector(0, 0, 1)
	light := NewAreaLight(corner, v1, 4, v2, 2, NewColor(1, 1, 1))

	testcases := []struct {
		u      float64
		v      float64
		result Tuple4
	}{
		{u: 0, v: 0, result: NewPoint(0.25, 0, 0.25)},
		{u: 1, v: 0, result: NewPoint(0.75, 0, 0.25)},
		{u: 0, v: 1, result: NewPoint(0.25, 0, 0.75)},
		{u: 2, v: 0, result: NewPoint(1.25, 0, 0.25)},
		{u: 3, v: 1, result: NewPoint(1.75, 0, 0.75)},
	}

	for _, tc := range testcases {
		pt := PointOnLightNoRandom(light, tc.u, tc.v)
		assert.True(t, TupleXYZEq(pt, tc.result))
	}
}

func TestLightingSamplesAreaLight(t *testing.T) {
	corner := NewPoint(-0.5, -0.5, -5)
	v1 := NewVector(1, 0, 0)
	v2 := NewVector(0, 1, 0)

	light := NewAreaLight(corner, v1, 2, v2, 2, NewColor(1, 1, 1))
	sh := NewSphere()
	sh.Material.Ambient = 0.1
	sh.Material.Diffuse = 0.9
	sh.Material.Specular = 0
	sh.Material.Color = NewColor(1, 1, 1)

	eye := NewPoint(0, 0, -5)

	testcases := []struct {
		point Tuple4
		color Tuple4
	}{
		{point: NewPoint(0, 0, -1), color: NewColor(0.9965, 0.9965, 0.9965)},
		{point: NewPoint(0, 0.7071, -0.7071), color: NewColor(0.6232, 0.6232, 0.6232)},
	}

	for _, tc := range testcases {
		eyev := Normalize(Sub(eye, tc.point))
		normalv := NewVector(tc.point[0], tc.point[1], tc.point[2])
		result := Lighting(sh.Material, sh, light, tc.point, eyev, normalv, 1.0, NewLightData())
		assert.InEpsilon(t, tc.color[0], result[0], Epsilon)
		assert.InEpsilon(t, tc.color[1], result[1], Epsilon)
		assert.InEpsilon(t, tc.color[2], result[2], Epsilon)
	}
	/*
		lighting() samples the area light
		  Given corner ← point(-0.5, -0.5, -5)
		    And v1 ← vector(1, 0, 0)
		    And v2 ← vector(0, 1, 0)
		    And light ← area_light(corner, v1, 2, v2, 2, color(1, 1, 1))
		    And shape ← sphere()
		    And shape.material.ambient ← 0.1
		    And shape.material.diffuse ← 0.9
		    And shape.material.specular ← 0
		    And shape.material.color ← color(1, 1, 1)
		    And eye ← point(0, 0, -5)
		    And pt ← <point>
		    And eyev ← normalize(eye - pt)
		    And normalv ← vector(pt.x, pt.y, pt.z)
		  When result ← lighting(shape.material, shape, light, pt, eyev, normalv, 1.0)
		  Then result = <result>
		  Examples:
		    | point                      | result                        |
		    | point(0, 0, -1)            | color(0.9965, 0.9965, 0.9965) |
		    | point(0, 0.7071, -0.7071)  | color(0.6232, 0.6232, 0.6232) |
	*/
}
