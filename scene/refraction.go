package scene

import "github.com/eriklupander/rt/internal/pkg/mat"

type Scene struct {
	Camera  mat.Camera
	Lights  []mat.Light
	Objects []mat.Shape
}

func Refraction() *Scene {
	camera := mat.NewCamera(600, 600, 0.5)
	camera.Transform = mat.ViewTransform(mat.NewPoint(-4.5, 0.85, -4), mat.NewPoint(0, 0.85, 0), mat.NewVector(0, 1, 0))
	camera.Inverse = mat.Inverse(camera.Transform)

	wallMat := mat.NewDefaultMaterial()
	ptrn := mat.NewCheckerPattern(mat.NewColor(0, 0, 0), mat.NewColor(0.75, 0.75, 0.74))
	ptrn.Transform = mat.Scale(0.5, 0.5, 0.5)
	wallMat.Pattern = ptrn
	wallMat.Specular = 0.0

	floor := mat.NewPlane()
	floor.SetTransform(mat.RotateY(0.31415))
	floorMat := mat.NewDefaultMaterial()
	floorMat.Pattern = ptrn
	floorMat.Ambient = 0.5
	floorMat.Diffuse = 0.4
	floorMat.Specular = 0.8
	floorMat.Reflectivity = 0.1
	floor.SetMaterial(floorMat)

	ceil := mat.NewPlane()
	ceil.SetTransform(mat.Translate(0, 5, 0))
	ceilMat := mat.NewDefaultMaterial()
	ceilPtrn := mat.NewCheckerPattern(mat.NewColor(0.85, 0.85, 0.85), mat.NewColor(1, 1, 1))
	ceilPtrn.Transform = mat.Scale(0.2, 0.2, 0.2)
	ceilMat.Pattern = ceilPtrn
	ceilMat.Ambient = 0.5
	ceilMat.Specular = 0
	ceil.SetMaterial(ceilMat)

	westWall := mat.NewPlane()
	westWall.SetTransform(mat.Translate(-5, 0, 0))
	westWall.SetTransform(mat.RotateZ(1.5708))
	westWall.SetTransform(mat.RotateY(1.5708))
	westWall.SetMaterial(wallMat)

	eastWall := mat.NewPlane()
	eastWall.SetTransform(mat.Translate(5, 0, 0))
	eastWall.SetTransform(mat.RotateZ(1.5708))
	eastWall.SetTransform(mat.RotateY(1.5708))
	eastWall.SetMaterial(wallMat)

	northWall := mat.NewPlane()
	northWall.SetTransform(mat.Translate(0, 0, 5))
	northWall.SetTransform(mat.RotateX(1.5708))
	northWall.SetMaterial(wallMat)

	southWall := mat.NewPlane()
	southWall.SetTransform(mat.Translate(0, 0, -5))
	southWall.SetTransform(mat.RotateX(1.5708))
	southWall.SetMaterial(wallMat)

	backBall1 := mat.NewSphere()
	backBall1.SetTransform(mat.Translate(4, 1, 4))
	mat1 := mat.NewDefaultMaterial()
	mat1.Color = mat.NewColor(0.8, 0.1, 0.3)
	mat1.Specular = 0
	backBall1.SetMaterial(mat1)

	backBall2 := mat.NewSphere()
	backBall2.SetTransform(mat.Translate(4.6, 0.4, 2.9))
	backBall2.SetTransform(mat.Scale(0.4, 0.4, 0.4))
	mat2 := mat.NewDefaultMaterial()
	mat2.Color = mat.NewColor(0.1, 0.8, 0.2)
	mat2.Shininess = 200
	backBall2.SetMaterial(mat2)

	backBall3 := mat.NewSphere()
	backBall3.SetTransform(mat.Translate(2.6, 0.6, 4.4))
	backBall3.SetTransform(mat.Scale(0.6, 0.6, 0.6))
	mat3 := mat.NewDefaultMaterial()
	mat3.Color = mat.NewColor(0.2, 0.1, 0.8)
	mat3.Shininess = 10
	mat3.Specular = 0.4
	backBall3.SetMaterial(mat3)

	glassBall := mat.NewSphere()
	glassBall.SetTransform(mat.Translate(0.25, 1, 0))
	glassBall.SetTransform(mat.Scale(1, 1, 1))

	glassMtrl := mat.NewMaterial(mat.NewColor(0.8, 0.8, 0.9), 0, 0.2, 0.9, 300)
	glassMtrl.Transparency = 0.8
	glassMtrl.RefractiveIndex = 1.5
	glassBall.SetMaterial(glassMtrl)

	return &Scene{
		Camera: camera,
		Lights: []mat.Light{mat.NewLight(mat.NewPoint(-4.9, 4.9, 1), mat.NewColor(1, 1, 1))},
		Objects: []mat.Shape{
			ceil, floor, northWall, eastWall, southWall, westWall, backBall1, backBall2, backBall3, glassBall,
		},
	}
}
