package scene

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/eriklupander/rt/internal/pkg/obj"
	"io/ioutil"
	"math"
)

func SimpleGopher() *Scene {

	camera := mat.NewCamera(640, 480, math.Pi/3.5)
	camera.Transform = mat.ViewTransform(mat.NewPoint(-.1, 1.2, 6), mat.NewPoint(0.05, 1.1, 0.05), mat.NewVector(0, 1, 0))
	camera.Inverse = mat.Inverse(camera.Transform)

	// Model
	bytes, _ := ioutil.ReadFile("assets/models/gopher.obj")
	model := obj.ParseObj(string(bytes)).ToGroup()
	model.SetTransform(mat.Translate(0, 1.2, 0))
	model.SetTransform(mat.RotateX(math.Pi / 2))
	model.SetTransform(mat.RotateY(-math.Pi / 2))
	model.SetTransform(mat.RotateX(-math.Pi / 8))
	mat.Divide(model, 100)
	model.Bounds()

	w := mat.NewWorld()
	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(3.3, 4, 10.5), mat.NewColor(1, 1, 1)))

	floor := mat.NewPlane()
	pm := mat.NewMaterial(mat.NewColor(1, 1, 1), 0.025, 0.67, 0, 200)
	pm.Reflectivity = 0.2
	floor.SetMaterial(pm)

	return &Scene{
		Camera: camera,
		// Lights: []mat.Light{{mat.NewPoint(-1, 1, 2.5), mat.NewColor(0.3, 0.3, 0.3)}},
		AreaLights: []mat.AreaLight{mat.NewAreaLight(
			mat.NewPoint(1.5, 2.5, 4.5),
			mat.NewVector(-.5, 0, .5),
			4,
			mat.NewVector(0, 1, 0),
			4,
			mat.NewColor(0.9, 0.9, 0.9))},
		Objects: []mat.Shape{
			floor, model,
		},
	}
}
