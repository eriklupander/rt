package scene

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"math"
)

func Hello() *Scene {
	floor := mat.NewPlane() //.UnitCube(material.Plastic(1, 1, 1, 0.05))
	plastic := mat.NewMaterial(mat.NewColor(1, 0.48, 0.43), 0.1, 0.3, 0.9, 300)

	floor.SetMaterial(plastic)
	floor.SetTransform(mat.Translate(0, -0.1, 0))

	ball := mat.NewSphere()
	gold := mat.NewMaterial(mat.NewColor(1, 0.88, 0.63), 0.1, 0.3, 0.9, 300)
	ball.SetMaterial(gold)
	ball.SetTransform(mat.Scale(0.1, 0.1, 0.1))

	c := mat.NewCamera(898, 450, math.Pi/3.5)
	c.Transform = mat.ViewTransform(mat.NewPoint(0, 0, -.5), mat.NewPoint(0, 0, 0), mat.NewVector(0, 1, 0))
	c.Inverse = mat.Inverse(c.Transform)

	return &Scene{
		Camera: c,
		Lights: []mat.Light{
			{Position: [4]float64{-1, 1, -3, 0}, Intensity: [4]float64{1, 1, 1, 1}},
		},
		Objects: []mat.Shape{
			floor, ball,
		},
	}
}
