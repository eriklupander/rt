package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseShapeTransform(t *testing.T) {
	s := NewSphere()
	s.SetTransform(Multiply(s.Transform, Translate(2, 3, 4)))
	assert.Equal(t, Translate(2, 3, 4), s.Transform)
}
func TestBaseShapeMaterial(t *testing.T) {
	s := NewSphere()
	s.SetMaterial(NewMaterial(NewColor(1, 1, 1), 1.0, 0.1, 0.1, 0.1))
	assert.Equal(t, 1.0, s.Material.Ambient)
}
