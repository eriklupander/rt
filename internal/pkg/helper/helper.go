package helper

import "github.com/eriklupander/rt/internal/pkg/mat"

var white = mat.NewColor(1, 1, 1)

// RenderPointAt transforms the passed world coordinate point into view coords and projects it onto the 2D canvas.
func RenderPointAt(canvas *mat.Canvas, camera mat.Camera, worldPoint mat.Tuple4, color mat.Tuple4) {

	// TransformRay point into camera space
	mat.MultiplyByTuplePtr(camera.Transform, worldPoint, &worldPoint)

	// View is always 1 unit away, divide by translated pixel z to get 3D -> 2D translation factor.
	dDividedByZ := 1 / worldPoint.Get(2)

	// Multiply x and y by translation factor to get 2D coords on the projection surface
	x := worldPoint.Get(0) * dDividedByZ
	y := worldPoint.Get(1) * dDividedByZ

	// transform from projection surface space into actual X/Y pixels based on pixelsize.
	// note that this translation isn't quite correct since I think I've gotten the offset
	// calculation wrong. Possibly, the raw col float should be rounded rather than floored.
	col := int(((camera.HalfWidth + x) / camera.PixelSize) - 0.5)
	row := int(((camera.HalfHeight + y) / camera.PixelSize) - 0.5)
	canvas.WritePixel(col, row, color)
}

func RenderReferenceAxises(canvas *mat.Canvas, camera mat.Camera) {
	xVec := mat.NewVector(1, 0, 0)
	yVec := mat.NewVector(0, 1, 0)
	zVec := mat.NewVector(0, 0, 1)
	point := mat.NewPoint(0, 0, 0)

	rx := mat.NewRay(point, xVec)
	ry := mat.NewRay(point, yVec)
	rz := mat.NewRay(point, zVec)

	for i := 0.0; i < 1.0; i += 0.005 {
		RenderPointAt(canvas, camera, mat.Position(rx, i), mat.NewColor(1, 0, 0))
		RenderPointAt(canvas, camera, mat.Position(ry, i), mat.NewColor(0, 1, 0))
		RenderPointAt(canvas, camera, mat.Position(rz, i), mat.NewColor(0, 0, 1))
	}
}

func RenderLine(canvas *mat.Canvas, camera mat.Camera, from, to mat.Tuple4) {
	vec := mat.Sub(to, from)
	magn := mat.Magnitude(vec)
	ray := mat.NewRay(from, mat.Normalize(vec))
	for i := 0.0; i < magn; i += 0.01 {
		RenderPointAt(canvas, camera, mat.Position(ray, i), white)
	}
}

func RenderLineBetweenShapes(canvas *mat.Canvas, s1, s2 mat.Shape, camera mat.Camera) {
	origin := mat.NewPoint(s1.GetTransform().Get(0, 3), s1.GetTransform().Get(1, 3), s1.GetTransform().Get(2, 3))
	target := mat.NewPoint(s2.GetTransform().Get(0, 3), s2.GetTransform().Get(1, 3), s2.GetTransform().Get(2, 3))
	vec := mat.Sub(target, origin)
	magn := mat.Magnitude(vec)
	ray := mat.NewRay(origin, mat.Normalize(vec))
	for i := 0.0; i < magn; i += 0.01 {
		RenderPointAt(canvas, camera, mat.Position(ray, i), white)
	}
}
