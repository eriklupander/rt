package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewLight(t *testing.T) {
	p := NewPoint(0, 0, 0)
	color := NewColor(1, 1, 1)
	light := NewLight(p, color)
	assert.True(t, TupleEquals(*light.Position, *p))
	assert.True(t, TupleEquals(*light.Intensity, *color))
}

func setupBase() (*Material, *Tuple4) {
	return NewDefaultMaterial(), NewPoint(0, 0, 0)
}

func TestLightEyeBetweenLightAndSphere(t *testing.T) {
	material, position := setupBase()
	eyev := NewVector(0, 0, -1)
	normalv := NewVector(0, 0, -1)
	light := NewLight(NewPoint(0, 0, -10), NewColor(1, 1, 1))
	result := Lighting(material, light, position, eyev, normalv)
	assert.Equal(t, 1.9, result.Get(0))
	assert.Equal(t, 1.9, result.Get(1))
	assert.Equal(t, 1.9, result.Get(2))
}
