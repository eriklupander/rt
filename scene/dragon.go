package scene

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/eriklupander/rt/internal/pkg/obj"
	"io/ioutil"
	"math"
)

func Dragon() *Scene {
	camera := mat.NewCamera(1920, 1080, math.Pi/3.5)
	viewTransform := mat.ViewTransform(mat.NewPoint(-3, 3.5, -10), mat.NewPoint(-0.5, 2, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	// Model
	bytes, _ := ioutil.ReadFile("assets/models/dragon.obj")
	model := obj.ParseObj(string(bytes)).ToGroup()
	dm := mat.NewDefaultMaterial()
	dm.Color = mat.NewColor(0.77, 0.62, 0.24)
	dm.Ambient = 0.25
	dm.Diffuse = 0.7
	dm.Specular = 0.6
	dm.Shininess = 51.2
	model.SetMaterial(dm)

	mat.Divide(model, 100)
	model.Bounds()

	floor := mat.NewPlane()
	pm := mat.NewMaterial(mat.NewColor(1, 1, 1), 0.025, 0.67, 0, 200)
	pm.Reflectivity = 0.1
	floor.SetMaterial(pm)

	wm := mat.NewMaterial(mat.NewColor(1, 1, 1), 0.025, 0.67, 0, 200)
	wm.Pattern = mat.NewStripePattern(mat.NewColor(0, 0, 1), mat.NewColor(1, 1, 1))
	northWall := mat.NewPlane()
	northWall.SetTransform(mat.Translate(0, 0, 15))
	northWall.SetTransform(mat.RotateX(1.5708))
	northWall.SetMaterial(wm)

	wm2 := mat.NewMaterial(mat.NewColor(1, 1, 1), 0.025, 0.67, 0, 200)
	ptrn := mat.NewStripePattern(mat.NewColor(0, 0, 1), mat.NewColor(1, 1, 1))
	ptrn.SetPatternTransform(mat.RotateY(math.Pi / 2))
	wm2.Pattern = ptrn
	eastWall := mat.NewPlane()
	eastWall.SetTransform(mat.Translate(15, 0, 0))
	eastWall.SetTransform(mat.RotateZ(1.5708))
	eastWall.SetMaterial(wm2)

	return &Scene{
		Camera: camera,
		//Lights: []mat.Light{mat.NewLight(mat.NewPoint(-30, 15, -15), mat.NewColor(1.4, 1.4, 1.4))},
		AreaLights: []mat.AreaLight{mat.NewAreaLight(
			mat.NewPoint(-30, 15, -15),
			mat.NewVector(1, 0, 0),
			4,
			mat.NewVector(0, 0, 1),
			4,
			mat.NewColor(1.4, 1.4, 1.4))},
		Objects: []mat.Shape{
			floor, northWall, eastWall, model,
		},
	}
}
