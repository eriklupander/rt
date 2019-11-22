package mat

type World struct {
	Light   Light
	Objects []Shape
}

func NewDefaultWorld() World {
	light := NewLight(NewPoint(-10, 10, -10), NewColor(1, 1, 1))
	material := NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200)
	//material := NewMaterial(NewColor(1,1,1), 0.1, 0.9, 0.9, 200)
	s1 := NewSphere()
	s1.Material = material
	s1.Label = "OUTER SPHERE"

	s2 := NewSphere()
	s2.Label = "INNER SPHERE"
	s2.SetTransform(Scale(0.5, 0.5, 0.5))
	return World{
		Light:   light,
		Objects: []Shape{s1, s2},
	}
}

func NewWorld() World {
	return World{}
}
