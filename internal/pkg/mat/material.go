package mat

type Material struct {
	Color        Tuple4
	Ambient      float64
	Diffuse      float64
	Specular     float64
	Shininess    float64
	Pattern      Pattern
	Reflectivity float64
}

func NewMaterial(color Tuple4, ambient float64, diffuse float64, specular float64, shininess float64) Material {
	return Material{Color: color, Ambient: ambient, Diffuse: diffuse, Specular: specular, Shininess: shininess}
}
func NewMaterialWithReflectivity(color Tuple4, ambient float64, diffuse float64, specular float64, shininess, reflectivity float64) Material {
	return Material{Color: color, Ambient: ambient, Diffuse: diffuse, Specular: specular, Shininess: shininess, Reflectivity: reflectivity}
}
func NewDefaultMaterial() Material {
	return Material{Color: NewColor(1, 1, 1), Ambient: 0.1, Diffuse: 0.9, Specular: 0.9, Shininess: 200.0, Reflectivity: 0.0}
}
func NewDefaultReflectiveMaterial(reflectivity float64) Material {
	m := NewDefaultMaterial()
	m.Reflectivity = reflectivity
	return m
}
func (m *Material) HasPattern() bool {
	return m.Pattern != nil
}
