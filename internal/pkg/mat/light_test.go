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

func setupBase() (Material, Tuple4) {
	return NewDefaultMaterial(), NewPoint(0, 0, 0)
}

func TestLightEyeBetweenLightAndSphere(t *testing.T) {
	material, position := setupBase()
	eyev := NewVector(0, 0, -1)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 0, -10), NewColor(1, 1, 1))
	result := Lighting(material, light, position, eyev, normalv, false)
	assert.Equal(t, 1.9, result.Get(0))
	assert.Equal(t, 1.9, result.Get(1))
	assert.Equal(t, 1.9, result.Get(2))
}
func TestLight180ToSurfaceEye45(t *testing.T) {
	material, position := setupBase()
	eyev := NewVector(0, 0, -1)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 10, -10), NewColor(1, 1, 1))
	result := Lighting(material, light, position, eyev, normalv, false)
	assert.InEpsilon(t, 0.7364, result.Get(0), Epsilon)
	assert.InEpsilon(t, 0.7364, result.Get(1), Epsilon)
	assert.InEpsilon(t, 0.7364, result.Get(2), Epsilon)
}
func TestLight45ToSurfaceEye180(t *testing.T) {
	material, position := setupBase()
	eyev := NewVector(0, rt2, -rt2)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 0, -10), NewColor(1, 1, 1))
	result := Lighting(material, light, position, eyev, normalv, false)
	assert.Equal(t, 1.0, result.Get(0))
	assert.Equal(t, 1.0, result.Get(1))
	assert.Equal(t, 1.0, result.Get(2))
}
func TestLight45ToSurfaceEye45(t *testing.T) {
	material, position := setupBase()
	eyev := NewVector(0, -rt2, -rt2)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 10, -10), NewColor(1, 1, 1))
	result := Lighting(material, light, position, eyev, normalv, false)
	assert.InEpsilon(t, 1.6364, result.Get(0), Epsilon)
	assert.InEpsilon(t, 1.6364, result.Get(1), Epsilon)
	assert.InEpsilon(t, 1.6364, result.Get(2), Epsilon)
}
func TestLightBehind(t *testing.T) {
	material, position := setupBase()
	eyev := NewVector(0, 0, -1)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 0, 10), NewColor(1, 1, 1))
	result := Lighting(material, light, position, eyev, normalv, false)
	assert.Equal(t, 0.1, result.Get(0))
	assert.Equal(t, 0.1, result.Get(1))
	assert.Equal(t, 0.1, result.Get(2))
}
