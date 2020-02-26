package main

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/constant"
	"github.com/eriklupander/rt/internal/pkg/helper"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/eriklupander/rt/internal/pkg/obj"
	"github.com/eriklupander/rt/internal/pkg/render"
	"github.com/eriklupander/rt/scene"
	"github.com/jinzhu/copier"
	"image"
	"image/png"
	"io/ioutil"
	"math"
	//_ "net/http/pprof"
	"os"
)

// main contains a load of old junk I've added while I completed chapters in the Ray Tracer Challenge book.
func main() {
	//runtime.SetBlockProfileRate(1)
	//runtime.SetMutexProfileFraction(1)
	// we need a webserver to get the pprof going
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	//parse()
	//csg()
	withModel()
	//groups()

	//refraction()
	//worldWithPlane() // REFERENCE IMAGE!!

	//termChan := make(chan os.Signal)
	//signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	//<-termChan // Blocks here!!
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

func csg() {

	camera := mat.NewCamera(640, 480, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-4, 2, -5), mat.NewPoint(0, 0, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(camera.Transform)
	worlds := make([]mat.World, constant.RenderThreads)
	for i := 0; i < constant.RenderThreads; i++ {
		w := mat.NewWorld()
		w.Light = append(w.Light, mat.NewLight(mat.NewPoint(0, 2, -2), mat.NewColor(1, 1, 1)))
		w.Light = append(w.Light, mat.NewLight(mat.NewPoint(0, 3, 0), mat.NewColor(1, 1, 1)))

		s1 := mat.NewSphere()
		m1 := mat.NewDefaultReflectiveMaterial(0.5)
		m1.Color = mat.NewColor(1, 0.1, 0.1)
		s1.SetMaterial(m1)
		c1 := mat.NewCube()
		m2 := mat.NewDefaultReflectiveMaterial(0.5)
		m2.Color = mat.NewColor(0.1, 0.1, 1.0)
		c1.SetMaterial(m1)
		c1.SetTransform(mat.Translate(-0.5, 0, 0))
		c1.SetTransform(mat.Scale(0.75, 0.5, 0.5))
		csg := mat.NewCSG("difference", s1, c1)
		csg.SetTransform(mat.Translate(0, 1, 0))
		csg.SetTransform(mat.RotateY(-math.Pi / 2))
		w.Objects = append(w.Objects, csg)

		floor := mat.NewPlane()
		floor.SetMaterial(mat.NewMaterialWithReflectivity(mat.NewColor(0.2, 0.2, 1.0), 0.1, 0.9, 0.7, 200, 0.0))
		w.Objects = append(w.Objects, floor)

		worlds[i] = w
	}

	writeImagePNG(render.Threaded(camera, worlds), "csg.png")
}

func withModel() {

	bytes, err := ioutil.ReadFile("assets/models/dragon.obj")
	parseObj := obj.ParseObj(string(bytes))

	// Model
	model := parseObj.ToGroup()
	model.SetTransform(mat.Translate(1, 0, 0))
	m := mat.NewDefaultMaterial()
	m.Color = mat.NewColor(0.92, 0.32, 0.3)
	m.Ambient = 0.3
	m.Diffuse = 0.6
	m.Specular = 0.3
	m.Shininess = 15
	model.SetMaterial(m)
	mat.Divide(model, 100)
	model.Bounds()

	if err != nil {
		panic(err.Error())
	}
	camera := mat.NewCamera(1280, 720, math.Pi/3)
	camera.Transform = mat.ViewTransform(mat.NewPoint(-8, 5.1, -8.5), mat.NewPoint(0, 2.5, 0), mat.NewVector(0, 1, 0))
	camera.Inverse = mat.Inverse(camera.Transform)

	worlds := setupModelScene(model, constant.RenderThreads)
	canvas := render.Threaded(camera, worlds)

	// write
	writeImagePNG(canvas, "dragon05-1.png")
}

func setupModelScene(model *mat.Group, instances int) []mat.World {
	worlds := make([]mat.World, instances, instances)
	for i := 0; i < instances; i++ {

		w := mat.NewWorld()
		w.Light = append(w.Light, mat.NewLight(mat.NewPoint(-4.5, 10, -6), mat.NewColor(1, 1, 1)))
		w.Light = append(w.Light, mat.NewLight(mat.NewPoint(-10, 10, 0), mat.NewColor(0.3, 0.3, 0.3)))
		w.Light = append(w.Light, mat.NewLight(mat.NewPoint(10, 10, 0), mat.NewColor(0.3, 0.3, 0.3)))

		dragon := mat.Group{}
		copier.Copy(&dragon, model)
		w.Objects = append(w.Objects, &dragon)

		floor := mat.NewPlane()
		floor.SetMaterial(mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.5, 0.5), 0.1, 0.9, 0.7, 200, 0))
		floor.Material.Pattern = mat.NewCheckerPattern(white, black)
		w.Objects = append(w.Objects, floor)

		northWall := mat.NewPlane()
		northWall.SetTransform(mat.Translate(0, 0, 10))
		northWall.SetTransform(mat.RotateX(1.5708))
		northWall.SetMaterial(mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.5, 0.5), 0.1, 0.9, 0.7, 200, 0))
		northWall.Material.Pattern = mat.NewCheckerPattern(white, black)
		w.Objects = append(w.Objects, northWall)
		worlds[i] = w
	}
	return worlds
}

func writeDataToPNG(canvas *mat.Canvas, myImage *image.RGBA) {
	for i := 0; i < len(canvas.Pixels); i++ {
		myImage.Pix[i*4] = clamp(canvas.Pixels[i][0])
		myImage.Pix[i*4+1] = clamp(canvas.Pixels[i][1])
		myImage.Pix[i*4+2] = clamp(canvas.Pixels[i][2])
		myImage.Pix[i*4+3] = 255
	}
}

// shows alternate way to load a scene and render it
func refraction() {
	sc := scene.Refraction()
	worlds := make([]mat.World, 8, 8)
	for i := 0; i < 8; i++ {
		w := mat.NewWorld()
		w.Light = sc.Lights
		w.Objects = sc.Objects
		worlds[i] = w
	}
	canvas := render.Threaded(sc.Camera, worlds)
	writeImagePNG(canvas, "refraction.png")
}

// This is my "reference image", used to benchmark the impl. in either 640x480 or 1920x1080
func worldWithPlane() {
	camera := mat.NewCamera(640, 480, math.Pi/3) // -4 är referens!
	viewTransform := mat.ViewTransform(mat.NewPoint(-2, 2.0, -4), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	light := mat.NewLight(mat.NewPoint(-5, 2.5, -3), mat.NewColor(1, 1, 1))

	worlds := make([]mat.World, constant.RenderThreads)
	for i := 0; i < constant.RenderThreads; i++ {
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
		s2.CastShadow = false
		s2.SetTransform(mat.Multiply(mat.Translate(-1, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
		mat2 := mat.NewMaterialWithReflectivity(mat.NewColor(0.1, 1.0, 0.1), 0.1, 0.5, 0.8, 220.0, 0.4)
		s2.SetMaterial(mat2)
		gr.AddChild(s2)

		s3 := mat.NewSphere()
		s3.SetTransform(mat.Multiply(mat.Translate(0, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
		mat3 := mat.NewMaterialWithReflectivity(mat.NewColor(0.1, 0.1, 1), 0.1, 0.5, 0.8, 220.0, 0.4)
		s3.SetMaterial(mat3)
		gr.AddChild(s3)

		gr.SetTransform(mat.RotateY(0.67))
		gr.Bounds() // For now, important to always call Bounds on Group once set up.
		mat.Divide(gr, 1)

		w.Objects = append(w.Objects, gr)
		w.Objects = append(w.Objects, gr.BoundsToCube())

		//cb := mat.NewCube()
		//w.Objects = append(w.Objects, cb)

		worlds[i] = w
	}

	canvas := render.Threaded(camera, worlds)

	// This is a hack, useful for debugging the renderering of a single pixel
	//color := render.RenderSinglePixel(camera, worlds, 300, 139)
	//fmt.Printf("%v\n", color)
	//canvas := mat.RenderThreaded(camera, w)

	// One can use this to render a unit-length XYZ axises superimposed on the image
	helper.RenderReferenceAxises(canvas, camera)

	// write
	myImage := image.NewRGBA(image.Rect(0, 0, canvas.W, canvas.H))

	writeDataToPNG(canvas, myImage)

	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test-sort-slice.png")
	if err != nil {
		panic(err.Error())
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, myImage)

	// Don't forget to close files
	outputFile.Close()
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

func writeImagePNG(canvas *mat.Canvas, filename string) {
	myImage := image.NewRGBA(image.Rect(0, 0, canvas.W, canvas.H))
	writeDataToPNG(canvas, myImage)
	outputFile, _ := os.Create(filename)
	_ = png.Encode(outputFile, myImage)
}
