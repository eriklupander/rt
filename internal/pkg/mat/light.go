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
func Lighting(material Material, object Shape, light Light, position, eyeVec, normalVec Tuple4, inShadow bool) Tuple4 {
	var color Tuple4
	if material.HasPattern() {
		color = PatternAtShape(material.Pattern, object, position)
	} else {
		color = material.Color
	}
	if inShadow {
		return MultiplyByScalar(color, material.Ambient)
	}
	effectiveColor := Hadamard(color, light.Intensity)

	// get vector from point on sphere to light source by subtracting, normalized into unit space.
	lightVec := Normalize(Sub(light.Position, position))

	// Add the ambient portion
	ambient := MultiplyByScalar(effectiveColor, material.Ambient)

	// get dot product (angle) between the light and normal  vectors. If negative, it means the light source is
	// on the backside
	lightDotNormal := Dot(lightVec, normalVec)
	specular := NewColor(0, 0, 0)
	diffuse := NewColor(0, 0, 0)
	if lightDotNormal < 0 {
		diffuse = black
		specular = black
	} else {
		// Diffuse contribution Precedense here??
		diffuse = MultiplyByScalar(effectiveColor, material.Diffuse*lightDotNormal)

		// reflect_dot_eye represents the cosine of the angle between the
		// reflection vector and the eye vector. A negative number means the
		// light reflects away from the eye.
		// Note that we negate the light vector since we want to angle of the bounce
		// of the light rather than the incoming light vector.
		reflectVec := Reflect(Negate(lightVec), normalVec)
		reflectDotEye := Dot(reflectVec, eyeVec)

		if reflectDotEye <= 0.0 {
			specular = black
		} else {
			// compute the specular contribution
			factor := math.Pow(reflectDotEye, material.Shininess)

			// again, check precedense here
			specular = MultiplyByScalar(light.Intensity, material.Specular*factor)
		}
	}
	// Add the three contributions together to get the final shading
	// Uses standard Tuple addition
	return Add(Add(ambient, diffuse), specular)
}
