package mat

type BoundingBox struct {
	PointA Tuple4
	PointB Tuple4
}

func NewBounds(pointA Tuple4, pointB Tuple4) *BoundingBox {
	return &BoundingBox{PointA: pointA, PointB: pointB}
}
func NewBoundsF(x1, y1, z1, x2, y2, z2 float64) *BoundingBox {
	return &BoundingBox{NewPoint(x1, y1, z1), NewPoint(x2, y2, z2)}
}

func FindBounds(shape Shape) *BoundingBox {
	return nil
}

func FindGroupBounds(group Group) *BoundingBox {
	return nil
}
