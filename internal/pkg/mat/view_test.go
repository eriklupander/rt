package mat

import (
	"github.com/stretchr/testify/assert"
	"math"
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

func TestConstructCamera(t *testing.T) {
	cam := NewCamera(160, 120, math.Pi/2.0)
	assert.Equal(t, 160, cam.Width)
	assert.Equal(t, 120, cam.Height)
	assert.Equal(t, math.Pi/2, cam.Fov)
	assert.Equal(t, Mat4x4(IdentityMatrix), cam.Transform)
}

func TestCalcLandscapeCanvasPixelSize(t *testing.T) {
	cam := NewCamera(200, 125, math.Pi/2.0)
	assert.Equal(t, 0.01, cam.PixelSize)

}
func TestCalculatePortraitCanvasPixelSize(t *testing.T) {
	cam := NewCamera(125, 200, math.Pi/2.0)
	assert.Equal(t, 0.01, cam.PixelSize)
}
