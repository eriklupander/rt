package mat

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultView(t *testing.T) {
	from := NewPoint(0, 0, 0)
	to := NewPoint(0, 0, -1) // look away
	up := NewVector(0, 1, 0)
	view := ViewTransform(from, to, up)
	assert.True(t, Equals(view, IdentityMatrix))
}

func TestViewInOppositeDirection(t *testing.T) {
	from := NewPoint(0, 0, 0)
	to := NewPoint(0, 0, 1) // look away
	up := NewVector(0, 1, 0)
	view := ViewTransform(from, to, up)
	assert.True(t, Equals(view, Scale(-1, 1, -1)))
}

func TestViewTransformMovesTheWorld(t *testing.T) {
	from := NewPoint(0, 0, 8)
	to := NewPoint(0, 0, 0)
	up := NewVector(0, 1, 0)
	view := ViewTransform(from, to, up)
	tr := Translate(0, 0, -8)

	assert.True(t, Equals(view, tr))
}
