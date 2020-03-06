package mat

type Mtl struct {
	Ambient         Tuple4
	Diffuse         Tuple4
	Specular        Tuple4
	Shininess       float64
	Reflectivity    float64
	Transparency    float64
	RefractiveIndex float64
	Name            string
}
