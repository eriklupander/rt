package mat

import "math"

type BoundingBox struct {
	Min Tuple4
	Max Tuple4
}

func NewEmptyBoundingBox() *BoundingBox {
	return &BoundingBox{
		Min: NewPoint(math.Inf(1), math.Inf(1), math.Inf(1)),
		Max: NewPoint(math.Inf(-1), math.Inf(-1), math.Inf(-1)),
	}
}
func NewBoundingBox(pointA Tuple4, pointB Tuple4) *BoundingBox {
	return &BoundingBox{Min: pointA, Max: pointB}
}
func NewBoundingBoxF(x1, y1, z1, x2, y2, z2 float64) *BoundingBox {
	return &BoundingBox{NewPoint(x1, y1, z1), NewPoint(x2, y2, z2)}
}

func (b *BoundingBox) ContainsPoint(p Tuple4) bool {
	return b.Min[0] <= p[0] && b.Min[1] <= p[1] && b.Min[2] <= p[2] &&
		b.Max[0] >= p[0] && b.Max[1] >= p[1] && b.Max[2] >= p[2]
}

func (b *BoundingBox) ContainsBox(b2 *BoundingBox) bool {
	//minOk := b.ContainsPoint(b2.Min)
	//maxOk := b.ContainsPoint(b2.Max)
	//return minOk && maxOk
	return b.ContainsPoint(b2.Min) && b.ContainsPoint(b2.Max)
}

func (b *BoundingBox) MergeWith(b2 *BoundingBox) {
	b.Add(b2.Min)
	b.Add(b2.Max)
}

func (b *BoundingBox) Add(p Tuple4) {
	if b.Min[0] > p[0] {
		b.Min[0] = p[0]
	}
	if b.Min[1] > p[1] {
		b.Min[1] = p[1]
	}
	if b.Min[2] > p[2] {
		b.Min[2] = p[2]
	}

	if b.Max[0] < p[0] {
		b.Max[0] = p[0]
	}
	if b.Max[1] < p[1] {
		b.Max[1] = p[1]
	}
	if b.Max[2] < p[2] {
		b.Max[2] = p[2]
	}
}

func ParentSpaceBounds(shape Shape) *BoundingBox {
	BoundingBox := BoundsOf(shape)
	return TransformBoundingBox(BoundingBox, shape.GetTransform())
}

func TransformBoundingBox(bbox *BoundingBox, m1 Mat4x4) *BoundingBox {
	p1 := bbox.Min
	p2 := NewPoint(bbox.Min[0], bbox.Min[1], bbox.Max[2])
	p3 := NewPoint(bbox.Min[0], bbox.Max[1], bbox.Min[2])
	p4 := NewPoint(bbox.Min[0], bbox.Max[1], bbox.Max[2])
	p5 := NewPoint(bbox.Max[0], bbox.Min[1], bbox.Min[2])
	p6 := NewPoint(bbox.Max[0], bbox.Min[1], bbox.Max[2])
	p7 := NewPoint(bbox.Max[0], bbox.Max[1], bbox.Min[2])
	p8 := bbox.Max

	out := NewEmptyBoundingBox()
	out.Add(MultiplyByTuple(m1, p1))
	out.Add(MultiplyByTuple(m1, p2))
	out.Add(MultiplyByTuple(m1, p3))
	out.Add(MultiplyByTuple(m1, p4))
	out.Add(MultiplyByTuple(m1, p5))
	out.Add(MultiplyByTuple(m1, p6))
	out.Add(MultiplyByTuple(m1, p7))
	out.Add(MultiplyByTuple(m1, p8))
	return out
}

func BoundsOf(shape Shape) *BoundingBox {
	switch val := shape.(type) {
	case *Group:
		box := NewEmptyBoundingBox()
		for i := 0; i < len(val.Children); i++ {
			cbox := ParentSpaceBounds(val.Children[i])
			box.MergeWith(cbox)
		}
		return box
	case *CSG:
		box := NewEmptyBoundingBox()
		box.MergeWith(ParentSpaceBounds(val.Left))
		box.MergeWith(ParentSpaceBounds(val.Right))
		return box
	case *Cube:
		return NewBoundingBoxF(-1, -1, -1, 1, 1, 1)
	case *Sphere:
		return NewBoundingBoxF(-1, -1, -1, 1, 1, 1)
	case *Plane:
		return NewBoundingBoxF(math.Inf(-1), 0, math.Inf(-1), math.Inf(1), 0, math.Inf(1))
	case *Cylinder:
		return NewBoundingBoxF(-1, val.MinY, -1, 1, val.MaxY, 1)
	case *Cone:
		xzMin := math.Abs(val.MinY)
		xzMax := math.Abs(val.MaxY)
		limit := xzMin
		if xzMax > limit {
			limit = xzMax
		}

		return NewBoundingBoxF(-limit, val.MinY, -limit, limit, val.MaxY, limit)
	case *Triangle:
		BoundingBox := NewEmptyBoundingBox()
		BoundingBox.Add(val.P1)
		BoundingBox.Add(val.P2)
		BoundingBox.Add(val.P3)
		return BoundingBox
	case *SmoothTriangle:
		BoundingBox := NewEmptyBoundingBox()
		BoundingBox.Add(val.P1)
		BoundingBox.Add(val.P2)
		BoundingBox.Add(val.P3)
		return BoundingBox

	default:
		return NewBoundingBoxF(-1, -1, -1, 1, 1, 1)
	}
}
