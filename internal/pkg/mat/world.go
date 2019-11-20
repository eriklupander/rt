package mat

type World struct {
	Light   Light
	Objects []Sphere
}

func NewDefaultWorld() World {
	light := NewLight(NewPoint(-10, 10, -10), NewColor(1, 1, 1))
	material := NewDefaultMaterial()
	s1 := NewSphere()
	s1.Material = material

	s2 := NewSphere()
	SetTransform(&s2, Scale(0.5, 0.5, 0.5))
	return World{
		Light:   light,
		Objects: []Sphere{s1, s2},
	}
}

func NewWorld() World {
	return World{}
}
