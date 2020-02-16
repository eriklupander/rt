package mat

import (
	"math"
	"math/rand"
)

type LightSource interface {
	Pos() Tuple4
	Intens() Tuple4
}

type Light struct {
	Position  Tuple4
	Intensity Tuple4
}

func (l Light) Pos() Tuple4 {
	return l.Position
}

func (l Light) Intens() Tuple4 {
	return l.Intensity
}

func NewLight(position Tuple4, intensity Tuple4) Light {

	return Light{
		Position:  position,
		Intensity: intensity,
	}
}

type AreaLight struct {
	Corner    Tuple4
	UVec      Tuple4
	USteps    int
	VVec      Tuple4
	VSteps    int
	Intensity Tuple4
	Samples   float64
	Position  Tuple4
}

func (al AreaLight) Pos() Tuple4 {
	return al.Position
}

func (al AreaLight) Intens() Tuple4 {
	return al.Intensity
}

func OrientAreaLight(light *AreaLight, source Tuple4, target Tuple4) {
	n := Normalize(Sub(target, source)) // Desired direction of the area light normal

	// Compute tangent and bitangent vectors
	a := NewVector(0, 1, 0)
	t := Normalize( Cross(a, n) )
	b := Cross( t, n )

	// Replace the uvec and vvec vectors, but preserve their length
	light.UVec = MultiplyByScalar(t, Magnitude(light.UVec))
	light.VVec = MultiplyByScalar(b, Magnitude(light.VVec))

	// Set the corner so that the source position falls in the center
	light.Corner = Sub(source, Sub(MultiplyByScalar(light.UVec, 0.5), MultiplyByScalar(light.VVec, 0.5)))
}

func NewAreaLight(corner Tuple4, uVec Tuple4, usteps int, vVec Tuple4, vsteps int, intensity Tuple4) AreaLight {
	return AreaLight{
		Corner:    corner,
		UVec:      DivideByScalar(uVec, float64(usteps)),
		USteps:    usteps,
		VVec:      DivideByScalar(vVec, float64(vsteps)),
		VSteps:    vsteps,
		Intensity: intensity,
		Samples:   float64(usteps * vsteps),
		Position: Add(corner, NewPoint(
			(uVec.Elems[0]+vVec.Elems[0])/2,
			(uVec.Elems[1]+vVec.Elems[1])/2,
			(uVec.Elems[2]+vVec.Elems[2])/2)),
	}
}

func PointOnLight(light AreaLight, u, v float64) Tuple4 {
	return Add(light.Corner,
		Add(
			MultiplyByScalar(light.UVec, u+rand.Float64()),  // used to be 0.5
			MultiplyByScalar(light.VVec, v+rand.Float64()))) // used to be 0.5
}

func Lighting(material Material, object Shape, light LightSource, position, eyeVec, normalVec Tuple4, intensity float64, lightData LightData) Tuple4 {
	var color Tuple4

	if material.HasPattern() {
		color = PatternAtShape(material.Pattern, object, position)
	} else {
		color = material.Color
	}
	if intensity == 0.0 {
		MultiplyByScalarPtr(color, material.Ambient, &lightData.EffectiveColor)
		return lightData.EffectiveColor
	}

	HadamardPtr(color, light.Intens(), &lightData.EffectiveColor)

	// sample each point on area light

	l := light.(AreaLight)
	//if !ok {
	//	panic("no support for point lights today!")
	//}
	sum := NewColor(0,0,0)
	for u := 0; u < l.USteps; u++ {
		for v := 0; v < l.VSteps; v++ {
			p := light.(AreaLight).Corner //PointOnLight(l, float64(u), float64(v))
			// get vector from point on sphere to light source by subtracting, normalized into unit space.
			SubPtr(p, position, &lightData.LightVec)
			NormalizePtr(lightData.LightVec, &lightData.LightVec)

			// Add the ambient portion
			MultiplyByScalarPtr(lightData.EffectiveColor, material.Ambient, &lightData.Ambient)

			lightDotNormal := Dot(lightData.LightVec, normalVec)

			// get dot product (angle) between the light and normal  vectors. If negative, it means the light source is
			// on the backside
			if lightDotNormal < 0 {
				lightData.Diffuse = black
				lightData.Specular = black
			} else {
				// Diffuse contribution Precedense here??

				MultiplyByScalarPtr(lightData.EffectiveColor, material.Diffuse*lightDotNormal, &lightData.Diffuse)

				// reflect_dot_eye represents the cosine of the angle between the
				// reflection vector and the eye vector. A negative number means the
				// light reflects away from the eye.
				// Note that we negate the light vector since we want to angle of the bounce
				// of the light rather than the incoming light vector.

				ReflectPtr(Negate(lightData.LightVec), normalVec, &lightData.ReflectVec)
				reflectDotEye := Dot(lightData.ReflectVec, eyeVec)

				if reflectDotEye <= 0.0 {
					lightData.Specular = black
				} else {
					// compute the specular contribution
					factor := math.Pow(reflectDotEye, material.Shininess)

					// again, check precedense here
					MultiplyByScalarPtr(light.Intens(), material.Specular*factor, &lightData.Specular)
				}
			}
			sum = Add(sum, Add(lightData.Diffuse, lightData.Specular))
		}
	}



	//intensity = 1.0
	// for soft shadows, multiply by intensity

	// Add the three contributions together to get the final shading
	// Uses standard Tuple addition
	return lightData.Ambient.Add(MultiplyByScalar(DivideByScalar(sum, l.Samples), intensity))
}

// Lighting computes the color of a given pixel given phong shading
func LightingPointLight(material Material, object Shape, light Light, position, eyeVec, normalVec Tuple4, intensity float64, lightData LightData) Tuple4 {
	var color Tuple4

	if material.HasPattern() {
		color = PatternAtShape(material.Pattern, object, position)
	} else {
		color = material.Color
	}
	if intensity == 0.0 {
		MultiplyByScalarPtr(color, material.Ambient, &lightData.EffectiveColor)
		return lightData.EffectiveColor
	}

	HadamardPtr(color, light.Intens(), &lightData.EffectiveColor)

	// get vector from point on sphere to light source by subtracting, normalized into unit space.
	SubPtr(light.Pos(), position, &lightData.LightVec)
	NormalizePtr(lightData.LightVec, &lightData.LightVec)

	// Add the ambient portion
	MultiplyByScalarPtr(lightData.EffectiveColor, material.Ambient, &lightData.Ambient)


	lightDotNormal := Dot(lightData.LightVec, normalVec)
	//diffuse := MultiplyByScalar(lightData.Diffuse, intensity)
	//specular := MultiplyByScalar(lightData.Specular, intensity)
	specular := lightData.Specular
	diffuse := lightData.Diffuse

	// get dot product (angle) between the light and normal  vectors. If negative, it means the light source is
	// on the backside
	if lightDotNormal < 0 {
		lightData.Diffuse = black
		lightData.Specular = black
	} else {
		// Diffuse contribution Precedense here??

		MultiplyByScalarPtr(lightData.EffectiveColor, material.Diffuse*lightDotNormal, &diffuse)

		// reflect_dot_eye represents the cosine of the angle between the
		// reflection vector and the eye vector. A negative number means the
		// light reflects away from the eye.
		// Note that we negate the light vector since we want to angle of the bounce
		// of the light rather than the incoming light vector.

		ReflectPtr(Negate(lightData.LightVec), normalVec, &lightData.ReflectVec)
		reflectDotEye := Dot(lightData.ReflectVec, eyeVec)

		if reflectDotEye <= 0.0 {
			lightData.Specular = black
		} else {
			// compute the specular contribution
			factor := math.Pow(reflectDotEye, material.Shininess)

			// again, check precedense here
			MultiplyByScalarPtr(light.Intens(), material.Specular*factor, &lightData.Specular)
		}
	}
	//intensity = 1.0
	// for soft shadows, multiply by intensity
	MultiplyByScalarPtr(diffuse, intensity, &diffuse)
	MultiplyByScalarPtr(specular, intensity, &specular)

	// Add the three contributions together to get the final shading
	// Uses standard Tuple addition
	return lightData.Ambient.Add(diffuse.Add(specular))
}
