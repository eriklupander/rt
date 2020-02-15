package main

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/eriklupander/rt/internal/pkg/obj"
	"github.com/eriklupander/rt/internal/pkg/render"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"

	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	rand.Seed(time.Now().Unix())
	// we need a webserver to get the pprof webserver
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	//parse()
	//csg()
	//withModel()
	//groups()
	//refraction()
	worldWithPlane() // REFERENCE IMAGE!!
	//renderworld()
	//shadedSphereDemo()
	//circleDemo()
	//clockDemo()
	//projectileDemo()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan // Blocks here!!
	fmt.Println("shutting down!")
}

var white = mat.NewColor(1, 1, 1)
var black = mat.NewColor(0, 0, 0)

//
//func parse() {
//	scene := parser.ParseYAML("scenes/simple.yaml")
//	fmt.Printf("%v", scene)
//
//	w := scene.World
//	w.Light = scene.Lights
//	camera := scene.Camera
//	canvas := mat.RenderThreaded(*camera, *w)
//	// writec
//	data := canvas.ToPPM()
//	err := ioutil.WriteFile("fromyaml.ppm", []byte(data), os.FileMode(0755))
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
//
//func csg() {
//	w := mat.NewWorld()
//	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(0, 2, -2), mat.NewColor(1, 1, 1)))
//	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(0, 3, 0), mat.NewColor(1, 1, 1)))
//
//	camera := mat.NewCamera(640, 480, math.Pi/3)
//	viewTransform := mat.ViewTransform(mat.NewPoint(-4, 2, -5), mat.NewPoint(0, 0, 0), mat.NewVector(0, 1, 0))
//	camera.Transform = viewTransform
//
//	s1 := mat.NewSphere()
//	m1 := mat.NewDefaultReflectiveMaterial(0.5)
//	m1.Color = mat.NewColor(1, 0.1, 0.1)
//	s1.SetMaterial(m1)
//	c1 := mat.NewCube()
//	m2 := mat.NewDefaultReflectiveMaterial(0.5)
//	m2.Color = mat.NewColor(0.1, 0.1, 1.0)
//	c1.SetMaterial(m1)
//	c1.SetTransform(mat.Translate(-0.5, 0, 0))
//	c1.SetTransform(mat.Scale(0.75, 0.5, 0.5))
//	csg := mat.NewCSG("difference", s1, c1)
//	csg.SetTransform(mat.Translate(0, 1, 0))
//	csg.SetTransform(mat.RotateY(-math.Pi / 2))
//	w.Objects = append(w.Objects, csg)
//
//	floor := mat.NewPlane()
//	floor.SetMaterial(mat.NewMaterialWithReflectivity(mat.NewColor(0.2, 0.2, 1.0), 0.1, 0.9, 0.7, 200, 0.0))
//	//floor.Material.Pattern = mat.NewCheckerPattern(white, black)
//	w.Objects = append(w.Objects, floor)
//
//	canvas := mat.RenderThreaded(camera, w)
//	// writec
//	data := canvas.ToPPM()
//	err := ioutil.WriteFile("csg1.ppm", []byte(data), os.FileMode(0755))
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
//
func withModel() {

	bytes, err := ioutil.ReadFile("Toilet.1.obj")
	if err != nil {
		panic(err.Error())
	}
	camera := mat.NewCamera(96, 72, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-4.3, 5, -8), mat.NewPoint(0, 2.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform

	worlds := make([]mat.World, 0)

	for i := 0; i < 8; i++ {
		parseObj := obj.ParseObj(string(bytes))

		w := mat.NewWorld()
		w.Objects = append(w.Objects, parseObj.ToGroup())
		//w.Objects[0].SetTransform(mat.Scale(0.6, 0.6, 0.6))
		m := mat.NewDefaultMaterial()
		m.Ambient = 0.3
		m.Reflectivity = 0.5
		w.Objects[0].SetMaterial(m)

		floor := mat.NewPlane()
		floor.SetMaterial(mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.5, 0.5), 0.1, 0.9, 0.7, 200, 0.1))
		floor.Material.Pattern = mat.NewCheckerPattern(white, black)
		w.Objects = append(w.Objects, floor)
		w.Light = append(w.Light, mat.NewLight(mat.NewPoint(-1.5, 2.5, -3), mat.NewColor(1, 1, 1)))
		worlds = append(worlds, w)
	}

	//canvas := mat.RenderThreaded(camera, w)
	canvas := render.Threaded(camera, worlds)
	//mat.RenderReferenceAxises(canvas, camera)

	// writec
	data := canvas.ToPPM()
	err = ioutil.WriteFile("toilet.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func refraction() {
	camera := mat.NewCamera(600, 600, 0.5)
	camera.Transform = mat.ViewTransform(mat.NewPoint(-4.5, 0.85, -4), mat.NewPoint(0, 0.85, 0), mat.NewVector(0, 1, 0))
	camera.Inverse = mat.Inverse(camera.Transform)

	worlds := make([]mat.World, 0)

	for i := 0; i < 8; i++ {
		wallMat := mat.NewDefaultMaterial()
		ptrn := mat.NewCheckerPattern(black, mat.NewColor(0.75, 0.75, 0.74))
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

		w := mat.NewWorld()
		w.Objects = append(w.Objects, ceil)
		w.Objects = append(w.Objects, floor)
		w.Objects = append(w.Objects, northWall)
		w.Objects = append(w.Objects, eastWall)
		w.Objects = append(w.Objects, southWall)
		w.Objects = append(w.Objects, westWall)

		w.Objects = append(w.Objects, backBall1)
		w.Objects = append(w.Objects, backBall2)
		w.Objects = append(w.Objects, backBall3)
		w.Objects = append(w.Objects, glassBall)

		light := mat.NewLight(mat.NewPoint(-4.9, 4.9, 1), mat.NewColor(1, 1, 1))

		w.Light = append(w.Light, light)

		worlds = append(worlds, w)
	}

	canvas := render.Threaded(camera, worlds)

	// write
	data := canvas.ToPPM()
	err := ioutil.WriteFile("refractions.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func worldWithPlane() {
	camera := mat.NewCamera(640, 480, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-2, 1.0, -4), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	light := mat.NewLight(mat.NewPoint(-5, 2.5, -3), mat.NewColor(1, 1, 1))

	worlds := make([]mat.World, 8)
	for i := 0; i < 8; i++ {
		w := mat.NewWorld()
		w.Light = append(w.Light, light)

		floor := mat.NewPlane()
		floor.SetTransform(mat.Translate(0, 0.01, 0))
		floor.SetMaterial(mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.5, 0.5), 0.1, 0.9, 0.7, 240, 0.2))
		floor.Material.Pattern = mat.NewCheckerPattern(white, black)
		floor.Material.Pattern.SetPatternTransform(mat.Scale(2, 2, 2))
		w.Objects = append(w.Objects, floor)

		ceil := mat.NewPlane()
		ceil.SetTransform(mat.Translate(0, 5, 0))
		ceilMat := mat.NewDefaultMaterial()
		ceilPtrn := mat.NewCheckerPattern(mat.NewColor(0.85, 0.85, 0.85), mat.NewColor(1, 1, 1))
		ceilPtrn.Transform = mat.Scale(0.2, 0.2, 0.2)
		ceilMat.Pattern = ceilPtrn
		ceilMat.Ambient = 0.5
		ceilMat.Specular = 0
		ceil.SetMaterial(ceilMat)
		w.Objects = append(w.Objects, ceil)

		wall := mat.NewPlane()
		wall.SetMaterial(mat.NewMaterial(mat.NewColor(0.9, 0.9, 0.9), 0.1, 0.9, 0.7, 200))
		wall.Material.Pattern = ceilPtrn
		wall.SetTransform(mat.Translate(0, 0, 8))
		wall.SetTransform(mat.RotateX(math.Pi / 2))
		w.Objects = append(w.Objects, wall)

		// transparent sphere
		middle := mat.NewSphere()
		middle.SetTransform(mat.Translate(-0.5, 0.75, 0.5))
		middle.SetTransform(mat.Scale(0.75, 0.75, 0.75))
		glassMtrl := mat.NewMaterial(mat.NewColor(0.8, 0.8, 0.9), 0, 0.2, 0.9, 300)
		glassMtrl.Transparency = 1.0
		glassMtrl.RefractiveIndex = 1.57
		glassMtrl.Reflectivity = 0.3
		middle.SetMaterial(glassMtrl)
		w.Objects = append(w.Objects, middle)

		// back sphere
		right := mat.NewSphere()
		right.SetTransform(mat.Multiply(mat.Translate(-0.75, 0.5, 2), mat.Scale(0.5, 0.5, 0.5)))
		right.Material = mat.NewDefaultMaterial()
		right.Material.Color = mat.NewColor(1, 0, 0)
		right.Material.Diffuse = 0.7
		right.Material.Specular = 0.9
		right.Material.Reflectivity = 0.0
		right.Material.Ambient = 0.1
		w.Objects = append(w.Objects, right)

		// cube
		cube := mat.NewCube()
		cube.SetTransform(mat.Multiply(mat.Translate(1, 0.25, -1.25), mat.Scale(0.25, 0.25, 0.25)))
		cube.Material = mat.NewDefaultMaterial()
		cube.Material.Color = mat.NewColor(1, 0.88, 0.63)
		cube.Material.Transparency = 0.0
		cube.Material.Diffuse = 0.3
		cube.Material.Specular = 0.9
		cube.Material.Reflectivity = 0.9
		cube.Material.Shininess = 300
		cube.Material.Ambient = 0.4
		//w.Objects = append(w.Objects, cube)

		//  Cylinder
		cyl := mat.NewCylinderMMC(0.0, 3.0, true)
		cyl.SetTransform(mat.Translate(1, 0.0, -1)) //1, 0.25, -1
		cyl.SetTransform(mat.Scale(0.2, 0.4, 0.2))
		cyl.Material.Color = mat.NewColor(1, 0.88, 0.63)
		cyl.Material.Transparency = 0.0
		cyl.Material.Diffuse = 0.3
		cyl.Material.Specular = 0.9
		cyl.Material.Reflectivity = 0.3
		cyl.Material.Shininess = 300
		w.Objects = append(w.Objects, cyl)

		gr := mat.NewGroup()

		s1 := mat.NewSphere()
		s1.SetTransform(mat.Multiply(mat.Translate(-2, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
		mat1 := mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.1, 0.1), 0.1, 0.5, 0.8, 220.0, 0.4)
		s1.SetMaterial(mat1)
		gr.AddChild(s1)

		s2 := mat.NewSphere()
		s2.SetTransform(mat.Multiply(mat.Translate(-1, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
		mat2 := mat.NewMaterialWithReflectivity(mat.NewColor(0.1, 1.0, 0.1), 0.1, 0.5, 0.8, 220.0, 0.4)
		s2.SetMaterial(mat2)
		gr.AddChild(s2)

		s3 := mat.NewSphere()
		s3.SetTransform(mat.Multiply(mat.Translate(0, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
		mat3 := mat.NewMaterialWithReflectivity(mat.NewColor(0.1, 0.1, 1), 0.1, 0.5, 0.8, 220.0, 0.4)
		s3.SetMaterial(mat3)
		gr.AddChild(s3)

		s4 := mat.NewSphere()
		s4.SetTransform(mat.Multiply(mat.Translate(1, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
		mat4 := mat.NewDefaultReflectiveMaterial(1.0)
		mat4.Color = mat.NewColor(0, 0, 0)
		s4.SetMaterial(mat4)
		//gr.AddChild(s4)

		w.Objects = append(w.Objects, gr)

		worlds[i] = w
	}

	canvas := render.Threaded(camera, worlds)
	// 300x139
	//color := render.RenderSinglePixel(camera, worlds, 300, 139)
	//fmt.Printf("%v\n", color)
	//canvas := mat.RenderThreaded(camera, w)
	//mat.RenderReferenceAxises(canvas, camera)

	// write
	myImage := image.NewRGBA(image.Rect(0, 0, canvas.W, canvas.H))

	for i := 0; i < len(canvas.Pixels); i++ {
		myImage.Pix[i*4] = clamp(canvas.Pixels[i].Elems[0])
		myImage.Pix[i*4+1] = clamp(canvas.Pixels[i].Elems[1])
		myImage.Pix[i*4+2] = clamp(canvas.Pixels[i].Elems[2])
		myImage.Pix[i*4+3] = 255
	}

	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test-0.5-shadow-50-samples.png")
	if err != nil {
		panic(err.Error())
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, myImage)

	// Don't forget to close files
	outputFile.Close()

	data := canvas.ToPPM()
	err = ioutil.WriteFile("world-transparency-new-threaded.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func clamp(clr float64) uint8 {
	c := clr * 255.0
	rounded := math.Round(c)
	if rounded > 255.0 {
		rounded = 255.0
	} else if rounded < 0.0 {
		rounded = 0.0
	}
	return uint8(rounded)
}

//
//func groups() {
//	w := mat.NewWorld()
//	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(-3, 2.5, -3), mat.NewColor(1, 1, 1)))
//
//	camera := mat.NewCamera(640, 480, math.Pi/3)
//	//camera := mat.NewCamera(320, 240, math.Pi/3)
//	viewTransform := mat.ViewTransform(mat.NewPoint(-1.3, 2, -5), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
//	camera.Transform = viewTransform
//	camera.Inverse = mat.Inverse(viewTransform)
//
//	gr := mat.NewGroup()
//
//	s1 := mat.NewSphere()
//	//s1.SetTransform(mat.Multiply(mat.Translate(-2, 0.25, -1), mat.Scale(1.25, 0.25, 0.25)))
//	s1.SetTransform(mat.Translate(-2, -1, 0))
//	mat1 := mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.1, 0.1), 0.1, 0.5, 0.8, 220.0, 0.4)
//	s1.SetMaterial(mat1)
//	gr.AddChild(s1)
//
//	s4 := mat.NewSphere()
//	//s4.SetTransform(mat.Multiply(mat.Translate(1, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
//	s4.SetTransform(mat.Translate(0, 0, 0))
//	mat4 := mat.NewMaterialWithReflectivity(mat.NewColor(0.1, 0.1, 1.0), 0.1, 0.5, 0.8, 220.0, 0.4)
//	s4.SetMaterial(mat4)
//	gr.AddChild(s4)
//
//	w.Objects = append(w.Objects, gr)
//
//	canvas := mat.Render(camera, w)
//	mat.RenderReferenceAxises(canvas, camera)
//
//	// write
//	data := canvas.ToPPM()
//	err := ioutil.WriteFile("world-group.ppm", []byte(data), os.FileMode(0755))
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
//
//func renderworld() {
//	w := mat.NewWorld()
//	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(-10, 1, -10), mat.NewColor(1, 1, 1)))
//	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(1, 13, 1), mat.NewColor(0.5, 0.5, 0.5)))
//
//	camera := mat.NewCamera(480, 320, math.Pi/3)
//	viewTransform := mat.ViewTransform(mat.NewPoint(0, 1.5, -5), mat.NewPoint(0, 1, 0), mat.NewVector(0, 1, 0))
//	camera.Transform = viewTransform
//
//	// Create floor
//	floor := mat.NewSphere()
//	floor.Transform = mat.Scale(10, 0.01, 10)
//	floor.Material = mat.NewDefaultMaterial()
//	floor.Material.Color = mat.NewColor(1, 0.9, 0.9)
//	floor.Material.Specular = 0.0
//	floor.Material.Reflectivity = 0.2
//	w.Objects = append(w.Objects, floor)
//
//	// create left wall
//	leftWall := mat.NewSphere()
//
//	scaleM := mat.Scale(10, 0.01, 10)
//	rotXM := mat.RotateX(math.Pi / 2)
//	rotYM := mat.RotateY(-math.Pi / 4)
//	transM := mat.Translate(0, 0, 5)
//
//	m1 := mat.Multiply(transM, rotYM)
//	m2 := mat.Multiply(m1, rotXM)
//	m3 := mat.Multiply(m2, scaleM)
//	leftWall.Transform = m3
//	leftWall.Material = floor.Material
//	w.Objects = append(w.Objects, leftWall)
//
//	// create right wall
//	rightWall := mat.NewSphere()
//
//	scaleM = mat.Scale(10, 0.01, 10)
//	rotXM = mat.RotateX(math.Pi / 2)
//	rotYM = mat.RotateY(math.Pi / 4)
//	transM = mat.Translate(0, 0, 5)
//
//	m1 = mat.Multiply(transM, rotYM)
//	m2 = mat.Multiply(m1, rotXM)
//	m3 = mat.Multiply(m2, scaleM)
//	rightWall.Transform = m3
//	rightWall.Material = floor.Material
//	w.Objects = append(w.Objects, rightWall)
//
//	// middle sphere
//	middle := mat.NewSphere()
//	middle.Transform = mat.Translate(-0.5, 1, 0.5)
//	middle.Material = mat.NewDefaultMaterial()
//	middle.Material.Color = mat.NewColor(0.1, 1, 0.5)
//	middle.Material.Diffuse = 0.7
//	middle.Material.Specular = 0.3
//	w.Objects = append(w.Objects, middle)
//
//	// right sphere
//	right := mat.NewSphere()
//	right.Transform = mat.Multiply(mat.Translate(1.5, 0.5, -0.5), mat.Scale(0.5, 0.5, 0.5))
//	right.Material = mat.NewDefaultMaterial()
//	right.Material.Color = mat.NewColor(0.5, 1, 0.1)
//	right.Material.Diffuse = 0.7
//	right.Material.Specular = 0.3
//	w.Objects = append(w.Objects, right)
//
//	// left sphere
//	left := mat.NewSphere()
//	left.Transform = mat.Multiply(mat.Translate(-1.5, 0.33, -0.75), mat.Scale(0.33, 0.33, 0.33))
//	left.Material = mat.NewDefaultMaterial()
//	left.Material.Color = mat.NewColor(1, 0.8, 0.1)
//	left.Material.Diffuse = 0.7
//	left.Material.Specular = 0.3
//	w.Objects = append(w.Objects, left)
//
//	// cube
//	cube := mat.NewCube()
//	cube.Transform = mat.Multiply(mat.Translate(-.6, 0.25, -1.5), mat.Scale(0.25, 0.25, 0.25))
//	cube.Material = mat.NewDefaultMaterial()
//	cube.Material.Color = mat.NewColor(1, 0.6, 0.2)
//	cube.Material.Transparency = 0.0
//	cube.Material.Diffuse = 0.7
//	cube.Material.Specular = 0.3
//	cube.Material.Reflectivity = 0.0
//	w.Objects = append(w.Objects, cube)
//
//	canvas := mat.RenderThreaded(camera, w)
//	// write
//	data := canvas.ToPPM()
//	err := ioutil.WriteFile("world1.ppm", []byte(data), os.FileMode(0755))
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
//
//func shadedSphereDemo() {
//	c := mat.NewCanvas(512, 512)
//
//	// this is our eye starting 15 units "in front" of origo.
//	rayOrigin := mat.NewPoint(0, 0, -15.0)
//
//	// Note!! If I understand this correctly, the "wall" in this case is actually an abstraction
//	// of the "far" wall of the view frustum, giving something to cast our ray against, forming a vector between
//	// the eye and a point in world space.
//	wallZ := 20.0
//	wallSize := 7.0
//	pixelSize := wallSize / float64(c.W)
//	half := wallSize / 2
//	sphere := mat.NewSphere()
//	//mat.SetTransform(&sphere, mat.Translate(1, 1, 1))
//	material := mat.NewDefaultMaterial()
//	material.Color = mat.NewColor(1, 0.2, 1)
//	sphere.SetMaterial(material)
//
//	lightPos := mat.NewPoint(-10, 10, -10)
//	lightColor := mat.NewColor(1, 1, 1)
//	light := mat.NewLight(lightPos, lightColor)
//
//	for row := 0; row < c.W; row++ {
//		worldY := half - pixelSize*float64(row)
//
//		for col := 0; col < c.H; col++ {
//			worldX := -half + pixelSize*float64(col)
//			posOnWall := mat.NewPoint(worldX, worldY, wallZ)
//
//			// Build a ray (origin + direction)
//			rayFromOriginToPosOnWall := mat.NewRay(rayOrigin, mat.Normalize(mat.Sub(posOnWall, rayOrigin)))
//
//			// check if our ray intersects the sphere
//			intersections := mat.IntersectRayWithShape(sphere, rayFromOriginToPosOnWall)
//			intersection, found := mat.Hit(intersections)
//
//			if found {
//				pointOfHit := mat.Position(rayFromOriginToPosOnWall, intersection.T)
//				normalAtHit := mat.NormalAt(sphere, pointOfHit, nil)
//				minusEyeRayVector := mat.Negate(rayFromOriginToPosOnWall.Direction)
//				color := mat.Lighting(sphere.Material, sphere, light, pointOfHit, minusEyeRayVector, normalAtHit, false)
//
//				c.WritePixel(col, c.H-row, color)
//			}
//		}
//	}
//	// write
//	data := c.ToPPM()
//	err := ioutil.WriteFile("shadedcircle.ppm", []byte(data), os.FileMode(0755))
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}
//
//func circleDemo() {
//	c := mat.NewCanvas(100, 100)
//
//	rayOrigin := mat.NewPoint(0, 0, -15.0)
//	wallZ := 20.0
//	wallSize := 7.0
//	pixelSize := wallSize / float64(c.W)
//	half := wallSize / 2
//	color := mat.NewColor(1, 0, 0)
//	sphere := mat.NewSphere()
//
//	//mat.SetTransform(sphere, mat.Scale(1, 0.5, 1))
//	//mat.SetTransform(sphere, mat.Multiply(mat.RotateZ(math.Pi/4), mat.Scale(0.5, 1, 1)))
//
//	for row := 0; row < c.W; row++ {
//		worldY := half - pixelSize*float64(row)
//
//		for col := 0; col < c.H; col++ {
//			worldX := -half + pixelSize*float64(col)
//			posOnWall := mat.NewPoint(worldX, worldY, wallZ)
//
//			rayFromOriginToPosOnWall := mat.NewRay(rayOrigin, mat.Normalize(mat.Sub(posOnWall, rayOrigin)))
//
//			// check if our ray intersects the sphere
//			intersections := mat.IntersectRayWithShape(sphere, rayFromOriginToPosOnWall)
//			_, found := mat.Hit(intersections)
//			if found {
//				c.WritePixel(col, c.H-row, color)
//			}
//		}
//	}
//	// write
//	data := c.ToPPM()
//	err := ioutil.WriteFile("circle.ppm", []byte(data), os.FileMode(0755))
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//}

func clockDemo() {
	c := mat.NewCanvas(80, 80)
	center := (c.W/2 + c.H/2) / 2
	white := mat.NewColor(1, 1, 1)

	point := mat.NewPoint(0, 1, 0)
	for i := 0; i < 12; i++ {
		rotation := float64(i) * (2 * math.Pi) / 12
		rotMat := mat.RotateZ(rotation)
		p2 := mat.MultiplyByTuple(rotMat, point)
		p2 = mat.MultiplyByScalar(p2, 30.0)
		c.WritePixel(center+int(p2.Get(0)), center-int(p2.Get(1)), white)
	}

	// write
	data := c.ToPPM()
	err := ioutil.WriteFile("clock.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func projectileDemo() {
	prj := NewProjectile(mat.NewPoint(0, 1, 0), mat.MultiplyByScalar(mat.Normalize(mat.NewVector(1, 1.8, 0)), 11.25))
	env := NewEnvironment(mat.NewVector(0, -0.1, 0), mat.NewVector(-0.01, 0, 0))
	c := mat.NewCanvas(900, 550)
	red := mat.NewColor(1, 1, 1)
	for prj.pos.Get(1) > 0.0 {
		tick(prj, env)
		//time.Sleep(time.Millisecond * 100)
		fmt.Printf("Projectile pos %v at height %v with velocity %v\n", mat.Magnitude(prj.pos), prj.pos.Get(1), prj.velocity)
		fmt.Printf("Drawing at: %d %d\n", int(prj.pos.Get(0)), c.H-int(prj.pos.Get(1)))
		c.WritePixel(int(prj.pos.Get(0)), c.H-int(prj.pos.Get(1)), red)
	}
	fmt.Printf("Projectile flew %v\n", mat.Magnitude(prj.pos))
	data := c.ToPPM()
	err := ioutil.WriteFile("pic.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func tick(prj *Projectile, env *Environment) {
	prj.pos = mat.Add(prj.pos, prj.velocity)
	prj.velocity = mat.Add(prj.velocity, env.gravity)
	prj.velocity = mat.Add(prj.velocity, env.wind)
}

type Environment struct {
	gravity mat.Tuple4
	wind    mat.Tuple4
}

func NewEnvironment(gravity mat.Tuple4, wind mat.Tuple4) *Environment {
	return &Environment{gravity: gravity, wind: wind}
}

type Projectile struct {
	pos      mat.Tuple4
	velocity mat.Tuple4
}

func NewProjectile(pos mat.Tuple4, velocity mat.Tuple4) *Projectile {
	return &Projectile{pos: pos, velocity: velocity}
}
