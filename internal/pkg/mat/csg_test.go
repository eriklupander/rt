package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateCSG(t *testing.T) {
	s1 := NewSphere()
	c1 := NewCube()
	csg := NewCSG("union", s1, c1)
	assert.Equal(t, "union", csg.Operation)
	assert.Equal(t, csg.Left, s1)
	assert.Equal(t, csg.Right, c1)
	assert.Equal(t, csg, s1.GetParent())
	assert.Equal(t, csg, c1.GetParent())
}
