package mat

type Material struct {
	Color           Tuple4
	Ambient         float64
	Diffuse         float64
	Specular        float64
	Shininess       float64
	Pattern         Pattern
	Reflectivity    float64
	Transparency    float64
	RefractiveIndex float64
	Name            string
}

func NewMaterial(color Tuple4, ambient float64, diffuse float64, specular float64, shininess float64) Material {
	return Material{Color: color, Ambient: ambient, Diffuse: diffuse, Specular: specular, Shininess: shininess, RefractiveIndex: 1.0}
}
func NewMaterialWithReflectivity(color Tuple4, ambient float64, diffuse float64, specular float64, shininess, reflectivity float64) Material {
	return Material{Color: color, Ambient: ambient, Diffuse: diffuse, Specular: specular, Shininess: shininess, Reflectivity: reflectivity, RefractiveIndex: 1.0}
}
func NewDefaultMaterial() Material {
	return Material{Color: NewColor(1, 1, 1), Ambient: 0.1, Diffuse: 0.9, Specular: 0.9, Shininess: 200.0, Reflectivity: 0.0, Transparency: 0.0, RefractiveIndex: 1.0}
}
func NewDefaultReflectiveMaterial(reflectivity float64) Material {
	m := NewDefaultMaterial()
	m.Reflectivity = reflectivity
	return m
}
func NewGlassMaterial(refractiveIndex float64) Material {
	m := NewDefaultMaterial()
	m.RefractiveIndex = refractiveIndex
	m.Transparency = 1.0
	return m
}
func (m *Material) HasPattern() bool {
	return m.Pattern != nil
}
