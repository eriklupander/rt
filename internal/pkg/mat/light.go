package mat

import (
	"math"
)

type Light struct {
	Position  Tuple4
	Intensity Tuple4

	cEffectiveColor Tuple4
	cAmbient        Tuple4
	cDiffuse        Tuple4
	cSpecular       Tuple4
	cColor          Tuple4
	cTempColor      Tuple4
	cLightVec       Tuple4
	cReflectVec     Tuple4
}

func NewLight(position Tuple4, intensity Tuple4) Light {

	return Light{
		Position:        position,
		Intensity:       intensity,
		cEffectiveColor: NewColor(0, 0, 0),
		cAmbient:        NewColor(0, 0, 0),
		cDiffuse:        NewColor(0, 0, 0),
		cSpecular:       NewColor(0, 0, 0),
		cColor:          NewColor(0, 0, 0),
		cLightVec:       NewVector(0, 0, 0),
		cReflectVec:     NewVector(0, 0, 0),
		cTempColor:      NewColor(0, 0, 0),
	}
}

func (l *Light) reset() {
	l.cEffectiveColor = l.cEffectiveColor.ResetVector()
	l.cAmbient = l.cAmbient.ResetVector()
	l.cDiffuse = l.cDiffuse.ResetVector()
	l.cSpecular = l.cSpecular.ResetVector()
	l.cColor = l.cColor.ResetVector()
	l.cTempColor = l.cTempColor.ResetVector()
	l.cLightVec = l.cLightVec.ResetVector()
	l.cReflectVec = l.cReflectVec.ResetVector()
}

// Lighting computes the color of a given pixel given phong shading
//func (l *Light) Lighting(material Material, object Shape, light Light, position, eyeVec, normalVec Tuple4, inShadow bool) Tuple4 {
//	// reset all stored stuff
//	//l.reset()
//
//	if material.HasPattern() {
//		l.cColor = PatternAtShape(material.Pattern, object, position)
//	} else {
//		l.cColor = material.Color
//	}
//	if inShadow {
//		return MultiplyByScalar(l.cColor, material.Ambient)
//	}
//	HadamardPtr(l.cColor, light.Intensity, &l.cEffectiveColor)
//
//	// get vector from point on sphere to light source by subtracting, normalized into unit space.
//	NormalizePtr2(Sub(light.Position, position), &l.cLightVec)
//
//	// Add the ambient portion
//	MultiplyByScalarPtr(l.cEffectiveColor, material.Ambient, &l.cAmbient)
//
//	// get dot product (angle) between the light and normal  vectors. If negative, it means the light source is
//	// on the backside
//	lightDotNormal := Dot(l.cLightVec, normalVec)
//	//specular := NewColor(0, 0, 0)
//	//diffuse := NewColor(0, 0, 0)
//	if lightDotNormal < 0 {
//		l.cDiffuse = black
//		l.cSpecular = black
//	} else {
//		// Diffuse contribution Precedense here??
//		l.cDiffuse = MultiplyByScalar(l.cEffectiveColor, material.Diffuse*lightDotNormal)
//
//		// reflect_dot_eye represents the cosine of the angle between the
//		// reflection vector and the eye vector. A negative number means the
//		// light reflects away from the eye.
//		// Note that we negate the light vector since we want to angle of the bounce
//		// of the light rather than the incoming light vector.
//		NegatePtr(l.cLightVec, &l.cLightVec)
//		ReflectPtr(l.cLightVec, normalVec, &l.cReflectVec)
//		reflectDotEye := Dot(l.cReflectVec, eyeVec)
//
//		if reflectDotEye <= 0.0 {
//			l.cSpecular = black
//		} else {
//			// compute the specular contribution
//			// (check what POW does mathematically, it combines two floating point numbers...)
//			factor := math.Pow(reflectDotEye, material.Shininess)
//
//			// again, check precedense here
//			MultiplyByScalarPtr(light.Intensity, material.Specular*factor, &l.cSpecular)
//		}
//	}
//	// Add the three contributions together to get the final shading
//	// Uses standard Tuple addition
//	Add3(l.cAmbient, l.cDiffuse, l.cSpecular, &l.cTempColor)
//
//	// I don't like this allocation, look at using pre-allocated storage for final color
//	return NewColor(l.cTempColor.Get(0), l.cTempColor.Get(1), l.cTempColor.Get(2))
//}

// Lighting computes the color of a given pixel given phong shading
func (l *Light) Lighting(material Material, object Shape, position, eyeVec, normalVec Tuple4, inShadow bool) Tuple4 {

	var color Tuple4
	if material.HasPattern() {
		color = PatternAtShape(material.Pattern, object, position)
	} else {
		color = material.Color
	}
	if inShadow {
		return MultiplyByScalar(color, material.Ambient)
	}
	effectiveColor := Hadamard(color, l.Intensity)

	// get vector from point on sphere to light source by subtracting, normalized into unit space.
	lightVec := Normalize(Sub(l.Position, position))

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
			// (check what POW does mathematically, it combines two floating point numbers...)
			factor := math.Pow(reflectDotEye, material.Shininess)

			// again, check precedense here
			specular = MultiplyByScalar(l.Intensity, material.Specular*factor)
		}
	}
	// Add the three contributions together to get the final shading
	// Uses standard Tuple addition
	return Add(Add(ambient, diffuse), specular)
}