package mat

func ViewTransform(from, to, up Tuple4) Mat4x4 {
	// Create a new matrix from the identity matrix.
	vt := Mat4x4{Elems: make([]float64, 16)}
	copy(vt.Elems, IdentityMatrix.Elems)

	// Sub creates the initial vector between the eye and what we're looking at.
	forward := Normalize(Sub(to, from))

	// Use the cross product to get the "third" axis (in this case, not the forward or up one)
	left := Cross(forward, Normalize(up))

	// Again, use cross product between the just computed left and forward to get the "true" up.
	trueUp := Cross(left, forward)

	// copy each axis into the matrix
	vt.Elems[0] = left.Get(0)
	vt.Elems[1] = left.Get(1)
	vt.Elems[2] = left.Get(3)

	vt.Elems[4] = trueUp.Get(0)
	vt.Elems[5] = trueUp.Get(1)
	vt.Elems[6] = trueUp.Get(2)

	vt.Elems[8] = -forward.Get(0)
	vt.Elems[9] = -forward.Get(1)
	vt.Elems[10] = -forward.Get(2)

	// finally, move the view matrix opposite the camera position to emulate that the camera has moved.
	return Multiply(vt, Translate(-from.Get(0), -from.Get(1), -from.Get(2)))
}
