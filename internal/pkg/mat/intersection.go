package mat

import "sort"

type Intersection struct {
	T float64
	S Sphere
}

func NewIntersection(t float64, s Sphere) Intersection {
	return Intersection{T: t, S: s}
}

func IntersectionEqual(i1, i2 Intersection) bool {
	return i1.T == i2.T && i1.S.Id == i2.S.Id
}

func IntersectWithWorld(w World, r Ray) []Intersection {
	xs := make([]Intersection, 0)
	for idx, _ := range w.Objects {
		intersections := IntersectRayWithSphere(w.Objects[idx], r)
		if len(intersections) > 0 {
			for innerIdx := range intersections {
				if intersections[innerIdx].T >= 0.0 {
					xs = append(xs, intersections[innerIdx])
				}
			}
		}
	}
	// Remember that we must sort away negative ones?
	sort.Slice(xs, func(i, j int) bool {
		return xs[i].T < xs[j].T
	})
	return xs
}
