package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultMaterial(t *testing.T) {
	m := NewDefaultMaterial()
	assert.True(t, TupleEquals(*NewColor(1, 1, 1), *m.Color))
	assert.Equal(t, 0.1, m.Ambient)
	assert.Equal(t, 0.9, m.Diffuse)
	assert.Equal(t, 0.9, m.Specular)
	assert.Equal(t, 200.0, m.Shininess)
}

func TestAssignMaterialToSphere(t *testing.T) {
	s := NewSphere()
	m := NewDefaultMaterial()
	m.Ambient = 1.0
	SetMaterial(s, m)
	assert.Equal(t, 1.0, s.Material.Ambient)

}