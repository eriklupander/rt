package scene

import (
	"github.com/eriklupander/rt/internal/pkg/constant"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/eriklupander/rt/internal/pkg/obj"
	"io/ioutil"
)

func Softshadows() *Scene {
	camera := mat.NewCamera(400, 160, 0.7854)
	viewTransform := mat.ViewTransform(mat.NewPoint(-3, 1.0, 2.5), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	cube := mat.NewCube()
	cm := mat.NewMaterial(mat.NewColor(1.5, 1.5, 1.5), 1, 0, 0, 100)
	cube.SetMaterial(cm)
	cube.SetTransform(mat.Translate(0, 3, 4))
	cube.SetTransform(mat.Scale(1, 1, 0.1))
	cube.CastShadow = false

	plane := mat.NewPlane()
	pm := mat.NewMaterial(mat.NewColor(1, 1, 1), 0.025, 0.67, 0, 200)
	plane.SetMaterial(pm)

	sphere1 := mat.NewSphere()
	sm1 := mat.NewMaterial(mat.NewColor(1, 0, 0), 0.1, 0.6, 0, 200)
	sm1.Reflectivity = 0.3
	sphere1.SetMaterial(sm1)
	sphere1.SetTransform(mat.Translate(0.5, 0.5, 0))
	sphere1.SetTransform(mat.Scale(0.5, 0.5, 0.5))

	sphere2 := mat.NewSphere()
	sm2 := mat.NewMaterial(mat.NewColor(0.5, 0.5, 1), 0.1, 0.6, 0, 200)
	sm2.Reflectivity = 0.3
	sphere2.SetMaterial(sm2)
	sphere2.SetTransform(mat.Translate(-0.25, 0.33, 0))
	sphere2.SetTransform(mat.Scale(0.33, 0.33, 0.33))

	bytes, _ := ioutil.ReadFile("assets/models/dragon.obj")

	// Model
	object := obj.ParseObj(string(bytes))
	model := object.ToGroup()
	model.SetTransform(mat.Translate(0.2, 0, 0.9))
	model.SetTransform(mat.Scale(0.2, 0.2, 0.2))
	m := mat.NewMaterial(mat.NewColor(1, 0, 0), 0.1, 0.6, 0, 200)
	m.Reflectivity = 0.2
	model.SetMaterial(m)
	mat.Divide(model, 100)
	model.Bounds()

	return &Scene{
		Camera: camera,
		//Lights: []mat.Light{mat.NewLight(mat.NewPoint(-4.9, 4.9, 1), mat.NewColor(1, 1, 1))},
		AreaLights: []mat.AreaLight{mat.NewAreaLight(
			mat.NewPoint(-1, 2, 4),
			mat.NewVector(2, 0, 0), constant.ShadowRays,
			mat.NewVector(0, 2, 0), constant.ShadowRays,
			mat.NewColor(1.5, 1.5, 1.5))},
		Objects: []mat.Shape{
			cube, plane, sphere1, sphere2, model,
		},
	}

}
