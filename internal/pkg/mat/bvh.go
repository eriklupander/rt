package mat

// SplitBounds splits the passed bounding box perpendicular of its longest axis. (Impl from bonus chapter)
func SplitBounds(b1 *BoundingBox) (*BoundingBox, *BoundingBox) {
	// find the box's largest dimension
	dx := b1.Max[0] - b1.Min[0]
	dy := b1.Max[1] - b1.Min[1]
	dz := b1.Max[2] - b1.Min[2]

	greatest := max(dx, dy, dz)

	// variables to help construct the points on
	// the dividing plane
	x0 := b1.Min[0]
	y0 := b1.Min[1]
	z0 := b1.Min[2]

	x1 := b1.Max[0]
	y1 := b1.Max[1]
	z1 := b1.Max[2]

	// adjust the points so that they lie on the
	// dividing plane
	if greatest == dx {
		x0 = x0 + dx/2.0
		x1 = x0
	} else if greatest == dy {
		y0 = y0 + dy/2.0
		y1 = y0
	} else {
		z0 = z0 + dz/2.0
		z1 = z0
	}

	midMin := NewPoint(x0, y0, z0)
	midMax := NewPoint(x1, y1, z1)

	// construct and return the two halves of
	// the bounding box
	left := NewBoundingBox(b1.Min, midMax)
	right := NewBoundingBox(midMin, b1.Max)

	return left, right
}

func PartitionChildren(g *Group) (*Group, *Group) {
	left := NewGroup()
	right := NewGroup()
	bbound := BoundsOf(g)
	leftBounds, rightBounds := SplitBounds(bbound)

	remain := make([]Shape, 0)
	for i := range g.Children {
		childBound := ParentSpaceBounds(g.Children[i]) //BoundsOf(g.Children[i])
		if leftBounds.ContainsBox(childBound) {
			left.AddChild(g.Children[i])
		} else if rightBounds.ContainsBox(childBound) {
			right.AddChild(g.Children[i])
		} else {
			remain = append(remain, g.Children[i])
		}
	}
	// copy over the remaining ones
	g.Children = g.Children[:0]
	g.Children = append(g.Children, remain...)

	// we should really automate bounds-recalc whenever a group is mutated...
	g.Bounds()
	left.Bounds()
	right.Bounds()
	return left, right
}

func MakeSubGroup(g *Group, shapes ...Shape) {
	subgroup := NewGroup()
	for i := range shapes {
		subgroup.AddChild(shapes[i])
	}
	g.AddChild(subgroup)
}

func Divide(s Shape, threshold int) {
	switch g := s.(type) {
	case *CSG:
		Divide(g.Left, threshold)
		Divide(g.Right, threshold)
	case *Group:
		if threshold <= len(g.Children) {
			// split members of group into left, right or remain
			left, right := PartitionChildren(g)

			// if left or right contains any shapes, create a new subgroup containing those shapes
			// and add that subgroup to the passed group
			if len(left.Children) > 0 {
				MakeSubGroup(g, left.Children...)
			}
			if len(right.Children) > 0 {
				MakeSubGroup(g, right.Children...)
			}
		}

		// Now, iterate over all children and recursivley call divide on them.
		for i := range g.Children {
			Divide(g.Children[i], threshold)
		}
	default:
		// Do nothing
	}
}

func remove(a []Shape, i int) []Shape {
	// Remove the element at index i from a.
	copy(a[i:], a[i+1:]) // Shift a[i+1:] left one index.
	//a[len(a)-1] =      // Erase last element (write zero value).
	a = a[:len(a)-1] // Truncate slice.
	return a
}
