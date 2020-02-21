package mat

import "math"

type BoundingBox struct {
	Min Tuple4
	Max Tuple4
}
func NewEmptyBoundingBox() BoundingBox {
	return BoundingBox{
		Min: Tuple4{math.Inf(1),math.Inf(1),math.Inf(1)},
		Max: Tuple4{math.Inf(-1),math.Inf(-1),math.Inf(-1)},
	}
}
func NewBoundingBox(pointA Tuple4, pointB Tuple4) *BoundingBox {
	return &BoundingBox{Min: pointA, Max: pointB}
}
func NewBoundingBoxF(x1, y1, z1, x2, y2, z2 float64) *BoundingBox {
	return &BoundingBox{NewPoint(x1, y1, z1), NewPoint(x2, y2, z2)}
}

func FindBounds(shape Shape) *BoundingBox {
	return nil
}

func FindGroupBounds(group Group) *BoundingBox {
	return nil
}
