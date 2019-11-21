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
	assert.Equal(t, cam.Width, 160)
	assert.Equal(t, cam.Height, 120)
	assert.Equal(t, cam.Fov, math.Pi/2)
	assert.Equal(t, cam.Transform, IdentityMatrix)
}

func TestCalcLandscapeCanvasPixelSize(t *testing.T) {
	cam := NewCamera(200, 125, math.Pi/2.0)
	assert.Equal(t, 0.01, cam.PixelSize)

}
func TestCalculatePortraitCanvasPixelSize(t *testing.T) {
	cam := NewCamera(125, 200, math.Pi/2.0)
	assert.Equal(t, 0.01, cam.PixelSize)
}

func TestRayForPixelThroughCenterOfCanvas(t *testing.T) {
	cam := NewCamera(201, 101, math.Pi/2.0)
	r := RayForPixel(cam, 100, 50)
	assert.Equal(t, NewPoint(0, 0, 0), r.Origin)
	assert.Equal(t, NewVector(0, 0, -1), r.Direction)
}

func TestRayForPixelThroughCornerOfCanvas(t *testing.T) {
	cam := NewCamera(201, 101, math.Pi/2.0)
	r := RayForPixel(cam, 0, 0)
	assert.Equal(t, NewPoint(0, 0, 0), r.Origin)
	assert.InEpsilon(t, 0.66519, r.Direction.Get(0), Epsilon)
	assert.InEpsilon(t, 0.33259, r.Direction.Get(1), Epsilon)
	assert.InEpsilon(t, -0.66851, r.Direction.Get(2), Epsilon)
}

// Page 103, third testx
func TestRayForPixelWhenCamIsTransformed(t *testing.T) {
	cam := NewCamera(201, 101, math.Pi/2.0)
	cam.Transform = Multiply(RotateY(math.Pi/4), Translate(0, -2, 5))
	r := RayForPixel(cam, 100, 50)
	assert.Equal(t, NewPoint(0, 2, -5), r.Origin)
	assert.True(t, TupleEquals(NewVector(math.Sqrt(2.0)/2.0, 0.0, -math.Sqrt(2.0)/2.0), r.Direction))
}

// Page 104
func TestRender(t *testing.T) {
	/*
		Scenario: Rendering a world with a camera
		Given w ← default_world()
		And c ← camera(11, 11, π/2)
		And from ← point(0, 0, -5)
		And to ← point(0, 0, 0)
		And up ← vector(0, 1, 0)
		And c.transform ← view_transform(from, to, up)
		When image ← render(c, w)
		Then pixel_at(image, 5, 5) = color(0.38066, 0.47583, 0.2855)
	*/
	w := NewDefaultWorld()
	c := NewCamera(11, 11, math.Pi/2)
	from := NewPoint(0, 0, -5)
	to := NewPoint(0, 0, 0)
	upVec := NewVector(0, 1, 0)
	c.Transform = ViewTransform(from, to, upVec)
	canvas := Render(c, w)

	assert.InEpsilon(t, 0.38066, canvas.ColorAt(5, 5).Get(0), Epsilon)
	assert.InEpsilon(t, 0.47583, canvas.ColorAt(5, 5).Get(1), Epsilon)
	assert.InEpsilon(t, 0.2855, canvas.ColorAt(5, 5).Get(2), Epsilon)
}
