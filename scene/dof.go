package scene

import (
	"github.com/eriklupander/rt/internal/pkg/constant"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"math"
)

func DoF() *Scene {
	camera := mat.NewCamera(640, 480, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-1, 0.5, -3), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)
	camera.Aperture = 0.05
	camera.FocalLength = 2.0

	cube := mat.NewCube()
	cm := mat.NewMaterial(mat.NewColor(1.5, 1.5, 1.5), 1, 0, 0, 100)
	cube.SetMaterial(cm)
	cube.SetTransform(mat.Translate(-5, 2, -6))
	cube.SetTransform(mat.Scale(1, 1, 0.1))
	cube.CastShadow = false

	plane := mat.NewPlane()
	pm := mat.NewMaterial(mat.NewColor(1, 1, 1), 0.025, 0.67, 0, 200)
	plane.SetMaterial(pm)

	sphere1 := mat.NewSphere()
	sm1 := mat.NewMaterial(mat.NewColor(1, 0, 0), 0.1, 0.6, 0, 200)
	sm1.Reflectivity = 0.3
	sphere1.SetMaterial(sm1)
	sphere1.SetTransform(mat.Translate(0.25, 0.5, -0.5))
	sphere1.SetTransform(mat.Scale(0.5, 0.5, 0.5))

	sphere2 := mat.NewSphere()
	sm2 := mat.NewMaterial(mat.NewColor(0.5, 0.5, 1), 0.1, 0.6, 0, 200)
	sm2.Reflectivity = 0.3
	sphere2.SetMaterial(sm2)
	sphere2.SetTransform(mat.Translate(.25, 0.33, -1.5))
	sphere2.SetTransform(mat.Scale(0.33, 0.33, 0.33))

	sphere3 := mat.NewSphere()
	sm3 := mat.NewMaterial(mat.NewColor(0.25, 1, 0.25), 0.1, 0.6, 0, 200)
	sm3.Reflectivity = 0.3
	sphere3.SetMaterial(sm3)
	sphere3.SetTransform(mat.Translate(-0.5, 0.5, 1))
	sphere3.SetTransform(mat.Scale(0.5, 0.5, 0.5))

	sphere4 := mat.NewSphere()
	sm4 := mat.NewMaterial(mat.NewColor(0.9, 0.9, 0.96), 0.1, 0.6, 0, 200)
	sm4.Reflectivity = 0.3
	sphere4.SetMaterial(sm4)
	sphere4.SetTransform(mat.Translate(-1.5, 0.5, 4))
	sphere4.SetTransform(mat.Scale(0.5, 0.5, 0.5))

	return &Scene{
		Camera: camera,
		AreaLights: []mat.AreaLight{mat.NewAreaLight(
			mat.NewPoint(-5, 2, -6),
			mat.NewVector(2, 0, 0), constant.ShadowRays,
			mat.NewVector(0, 2, 0), constant.ShadowRays,
			mat.NewColor(1.5, 1.5, 1.5))},
		Objects: []mat.Shape{
			cube, plane, sphere1, sphere2, sphere3, sphere4,
		},
	}

}
