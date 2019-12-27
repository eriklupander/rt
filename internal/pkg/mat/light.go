package mat

import (
	"math"
)

type Light struct {
	Position  Tuple4
	Intensity Tuple4
}

func NewLight(position Tuple4, intensity Tuple4) Light {
	return Light{Position: position, Intensity: intensity}
}

// Lighting computes the color of a given pixel given phong shading
func Lighting(material Material, object Shape, light Light, position, eyeVec, normalVec Tuple4, inShadow bool, lightData LightData) Tuple4 {
	var color Tuple4
	if material.HasPattern() {
		color = PatternAtShape(material.Pattern, object, position)
	} else {
		color = material.Color
	}
	if inShadow {
		MultiplyByScalarPtr(color, material.Ambient, &lightData.EffectiveColor)
		return lightData.EffectiveColor
	}

	HadamardPtr(color, light.Intensity, &lightData.EffectiveColor)

	// get vector from point on sphere to light source by subtracting, normalized into unit space.
	SubPtr(light.Position, position, &lightData.LightVec)
	NormalizePtr(lightData.LightVec, &lightData.LightVec)

	// Add the ambient portion
	MultiplyByScalarPtr(lightData.EffectiveColor, material.Ambient, &lightData.Ambient)

	// get dot product (angle) between the light and normal  vectors. If negative, it means the light source is
	// on the backside
	lightDotNormal := Dot(lightData.LightVec, normalVec)
	specular := lightData.Specular
	diffuse := lightData.Diffuse
	if lightDotNormal < 0 {
		diffuse = black
		specular = black
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
			specular = black
		} else {
			// compute the specular contribution
			factor := math.Pow(reflectDotEye, material.Shininess)

			// again, check precedense here
			MultiplyByScalarPtr(light.Intensity, material.Specular*factor, &specular)
		}
	}
	// Add the three contributions together to get the final shading
	// Uses standard Tuple addition
	return lightData.Ambient.Add(diffuse.Add(specular))
	//return Add(Add(ambient, diffuse), specular)
}
