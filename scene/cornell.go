package scene

import (
	"github.com/eriklupander/rt/internal/pkg/config"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"math"
)

func Cornell() *Scene {
	camera := mat.NewCamera(config.Cfg.Width, config.Cfg.Height, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(0, 2, -10), mat.NewPoint(0, 2.0, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	s1 := mat.NewSphere()
	s1.SetTransform(mat.Translate(-1, 0, 0))
	s1.SetMaterial(mat.NewMaterial(mat.NewColor(.75, .25, .25), 0.9, 0.6, 0.7, 200))

	s2 := mat.NewSphere()
	s2.SetTransform(mat.Translate(1, 0, 0))
	s2.SetMaterial(mat.NewMaterial(mat.NewColor(.25, .25, .75), 0.9, 0.6, 0.7, 200))

	btm := mat.NewPlane()
	btm.SetTransform(mat.Translate(0, -1, 0))
	btm.SetMaterial(mat.NewMaterial(mat.NewColor(.75, .75, .75), 0.9, 0.6, 0.7, 200))

	back := mat.NewPlane()
	back.SetTransform(mat.Translate(0, 0, 4))
	back.SetTransform(mat.RotateX(math.Pi / 2))
	back.SetMaterial(mat.NewMaterial(mat.NewColor(.75, .75, .75), 0.9, 0.6, 0.7, 200))

	tp := mat.NewPlane()
	tp.SetTransform(mat.Translate(0, 4, 0))
	tp.SetMaterial(mat.NewMaterial(mat.NewColor(.75, .75, .75), 0.9, 0.6, 0.7, 200))

	left := mat.NewPlane()
	left.SetTransform(mat.Translate(-3, 0, 0))
	left.SetTransform(mat.RotateZ(math.Pi / 2))
	left.SetMaterial(mat.NewMaterial(mat.NewColor(.9, 0, 0), 0.9, 0.6, 0.7, 200))

	right := mat.NewPlane()
	right.SetTransform(mat.Translate(3, 0, 0))
	right.SetTransform(mat.RotateZ(math.Pi / -2))
	right.SetMaterial(mat.NewMaterial(mat.NewColor(0, 0.9, 0), 0.9, 0.6, 0.7, 200))

	light := mat.NewCube()
	light.SetTransform(mat.Translate(0, 3.9, 0))
	light.SetTransform(mat.Scale(0.4, 0.1, 0.4))
	light.SetMaterial(mat.NewMaterial(mat.NewColor(1, 1, 1), 1, 1, 1, 1))

	return &Scene{
		Camera: camera,
		Lights: []mat.Light{mat.NewLight(mat.NewPoint(0, 3.8, 0), mat.NewColor(1, 1, 1))},
		Objects: []mat.Shape{
			s1,
			s2,
			back,
			btm,
			tp,
			left,
			right,
			light,
		},
	}

}

/*

 */
