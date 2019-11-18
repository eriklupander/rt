package mat

func NormalAt(s *Sphere, x, y, z float64) *Tuple4 {
	origin := NewPoint(0, 0, 0)
	point := NewPoint(x, y, z)
	return Normalize(*Sub(*point, *origin))
}
