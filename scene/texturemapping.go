package scene

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/config"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"image"
	"os"
)

func TextureMapping() *Scene {
	camera := mat.NewCamera(config.Cfg.Width, config.Cfg.Height, 0.7854)
	viewTransform := mat.ViewTransform(mat.NewPoint(-3, 1.0, 2.5), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	cube := mat.NewCube()
	cm := mat.NewMaterial(mat.NewColor(1.5, 1.5, 1.5), 1, 0, 0, 100)
	cube.SetMaterial(cm)
	cube.SetTransform(mat.Translate(0, 3, 4))
	cube.SetTransform(mat.Scale(1, 1, 0.1))
	cube.CastShadow = false

	imgfile, err := os.Open("./assets/textures/marbletile.jpg")
	defer imgfile.Close()
	if err != nil {
		panic(err.Error())
	}
	img, _, err := image.Decode(imgfile)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Bounds: %v+\n", img.Bounds())

	plane := mat.NewPlane()
	pm := mat.NewMaterial(mat.NewColor(1, 1, 1), 0.025, 0.67, 0, 200)
	//pm.Texture = img
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

	return &Scene{
		Camera: camera,
		//Lights: []mat.Light{mat.NewLight(mat.NewPoint(-4.9, 4.9, 1), mat.NewColor(1, 1, 1))},
		AreaLights: []mat.AreaLight{mat.NewAreaLight(
			mat.NewPoint(-1, 2, 4),
			mat.NewVector(2, 0, 0), config.Cfg.SoftShadowSamples,
			mat.NewVector(0, 2, 0), config.Cfg.SoftShadowSamples,
			mat.NewColor(1.5, 1.5, 1.5))},
		Objects: []mat.Shape{
			cube, plane, sphere1, sphere2,
		},
	}

}
