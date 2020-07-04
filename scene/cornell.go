package scene

import (
	"github.com/eriklupander/rt/internal/pkg/config"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"math"
)

func Cornell() *Scene {
	camera := mat.NewCamera(config.Cfg.Width, config.Cfg.Height, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(0, 0, -100), mat.NewPoint(0, 0.0, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	s1 := mat.NewSphere()
	s1.SetTransform(mat.Translate(1e5+1, 40.8, 81.6))
	s1.SetTransform(mat.Scale(1e5, 1e5, 1e5))
	s1.SetMaterial(mat.NewMaterial(mat.NewColor(.75, .25, .25), 0.9, 0.6, 0.7, 200))

	s2 := mat.NewSphere()
	s2.SetTransform(mat.Translate(-1e5+99, 40.8, 81.6))
	s2.SetTransform(mat.Scale(1e5, 1e5, 1e5))
	s2.SetMaterial(mat.NewMaterial(mat.NewColor(.25, .25, .75), 0.9, 0.6, 0.7, 200))

	s3 := mat.NewSphere()
	s3.SetTransform(mat.Translate(50, 40.8, 1e5))
	s3.SetTransform(mat.Scale(1e5, 1e5, 1e5))
	s3.SetMaterial(mat.NewMaterial(mat.NewColor(.75, .75, .75), 0.9, 0.6, 0.7, 200))

	s4 := mat.NewSphere()
	s4.SetTransform(mat.Translate(50, 40.8, -1e5+170))
	s4.SetTransform(mat.Scale(1e5, 1e5, 1e5))
	s4.SetMaterial(mat.NewMaterial(mat.NewColor(0, 0, 0), 0.9, 0.6, 0.7, 200))

	btm := mat.NewSphere()
	btm.SetTransform(mat.Translate(50, 1e5, 81.6))
	btm.SetTransform(mat.Scale(1e5, 1e5, 1e5))
	btm.SetMaterial(mat.NewMaterial(mat.NewColor(.75, .75, .75), 0.9, 0.6, 0.7, 200))

	tp := mat.NewSphere()
	tp.SetTransform(mat.Translate(50, -1e5+81.6, 81.6))
	tp.SetTransform(mat.Scale(1e5, 1e5, 1e5))
	tp.SetMaterial(mat.NewMaterial(mat.NewColor(.75, .75, .75), 0.9, 0.6, 0.7, 200))

	mirr := mat.NewSphere()
	mirr.SetTransform(mat.Translate(27, 16.5, 47))
	mirr.SetTransform(mat.Scale(16.5, 16.5, 16.5))
	mirr.SetMaterial(mat.NewMaterial(mat.NewColor(0.99, 0.99, 0.99), 0.9, 0.6, 0.7, 200))

	glas := mat.NewSphere()
	glas.SetTransform(mat.Translate(73, 16.5, 78))
	glas.SetTransform(mat.Scale(16.5, 16.5, 16.5))
	glas.SetMaterial(mat.NewMaterial(mat.NewColor(0.99, 0.99, 0.99), 0.9, 0.6, 0.7, 200))

	lght := mat.NewSphere()
	lght.SetTransform(mat.Translate(50, 681.6-.27, 81.6))
	lght.SetTransform(mat.Scale(600, 600, 600))
	lght.SetMaterial(mat.NewMaterial(mat.NewColor(12, 12, 12), 0.9, 0.6, 0.7, 200))

	//		Sphere(16.5,Vec(27,16.5,47),       Vec(),Vec(1,1,1)*.999, SPEC),//Mirr
	//		Sphere(16.5,Vec(73,16.5,78),       Vec(),Vec(1,1,1)*.999, REFR),//Glas
	//		Sphere(600, Vec(50,681.6-.27,81.6),Vec(12,12,12),  Vec(), DIFF) //Lite
	//};

	return &Scene{
		Camera: camera,
		Lights: []mat.Light{mat.NewLight(mat.NewPoint(-4.9, 4.9, -25), mat.NewColor(1, 1, 1))},
		Objects: []mat.Shape{
			s1, s2, s3, s4,
			//btm, tp,
			glas, mirr,
			//lght,
		},
	}

}

/*

 */
