package mat

// RenderPointAt transforms the passed world coordinate point into view coords and projects it onto the 2D canvas.
func RenderPointAt(canvas *Canvas, camera Camera, worldPoint Tuple4, color Tuple4) {

	// TransformRay point into camera space
	MultiplyByTuplePtr(camera.Transform, worldPoint, &worldPoint)

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

func RenderReferenceAxises(canvas *Canvas, camera Camera) {
	xVec := NewVector(1, 0, 0)
	yVec := NewVector(0, 1, 0)
	zVec := NewVector(0, 0, 1)
	point := NewPoint(0, 0, 0)

	rx := NewRay(point, xVec)
	ry := NewRay(point, yVec)
	rz := NewRay(point, zVec)

	for i := 0.0; i < 1.0; i += 0.005 {
		RenderPointAt(canvas, camera, Position(rx, i), NewColor(1, 0, 0))
		RenderPointAt(canvas, camera, Position(ry, i), NewColor(0, 1, 0))
		RenderPointAt(canvas, camera, Position(rz, i), NewColor(0, 0, 1))
	}
}

func RenderLine(canvas *Canvas, camera Camera, from, to Tuple4) {
	vec := Sub(to, from)
	magn := Magnitude(vec)
	ray := NewRay(from, Normalize(vec))
	for i := 0.0; i < magn; i += 0.01 {
		RenderPointAt(canvas, camera, Position(ray, i), white)
	}
}

func RenderLineBetweenShapes(canvas *Canvas, s1, s2 Shape, camera Camera) {
	origin := NewPoint(s1.GetTransform().Get(0, 3), s1.GetTransform().Get(1, 3), s1.GetTransform().Get(2, 3))
	target := NewPoint(s2.GetTransform().Get(0, 3), s2.GetTransform().Get(1, 3), s2.GetTransform().Get(2, 3))
	vec := Sub(target, origin)
	magn := Magnitude(vec)
	ray := NewRay(origin, Normalize(vec))
	for i := 0.0; i < magn; i += 0.01 {
		RenderPointAt(canvas, camera, Position(ray, i), white)
	}
}
