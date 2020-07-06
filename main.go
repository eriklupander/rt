package main

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/config"
	"github.com/eriklupander/rt/internal/pkg/helper"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/eriklupander/rt/internal/pkg/render"
	"github.com/eriklupander/rt/scene"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"image"
	"image/png"
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"
)

// main contains a load of old junk I've added while I completed chapters in the Ray Tracer Challenge book.
func main() {

	var configFlags = pflag.NewFlagSet("config", pflag.ExitOnError)
	configFlags.Int("threads", runtime.NumCPU(), "Image width")
	configFlags.Int("width", 640, "Image width")
	configFlags.Int("height", 480, "Image height")
	configFlags.Int("samples", 1, "Number of samples per pixel")
	configFlags.Int("softshadowsamples", 0, "Number of shadow rays for soft shadows")
	configFlags.String("scene", "reference", "scene from /scenes")
	//configFlags.String("output",  time.Now().Format(time.RFC3339)+ ".png", "output filename")

	if err := configFlags.Parse(os.Args[1:]); err != nil {
		panic(err.Error())
	}
	if err := viper.BindPFlags(configFlags); err != nil {
		panic(err.Error())
	}
	viper.AutomaticEnv()

	config.FromConfig()
	fmt.Printf("Running with %d CPUs\n", viper.GetInt("threads"))

	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	// we need a webserver to get the pprof going
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	switch viper.GetString("scene") {
	case "threads":
		threaded()
	case "reference":
		worldWithPlane()
	case "dragon":
		withDragonModel()
	case "gopher":
		withSimpleGopherModel()
	case "csg":
		csg()
	case "softshadows":
		softshadows()
	case "refraction":
		refraction()
	case "dof":
		depthOfField()
	case "cornell":
		cornell()
	default:
		fmt.Println("no scene specified, rendering reference scene")
	}

	fmt.Printf("Rendered scene '%v'\n", viper.GetString("scene"))
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
func cornell() {

	worlds := make([]mat.World, config.Cfg.Threads, config.Cfg.Threads)
	sc := scene.Cornell()
	for i := 0; i < config.Cfg.Threads; i++ {
		w := mat.NewWorld()
		sc := scene.Cornell()
		w.Light = sc.Lights
		//w.AreaLight = sc.AreaLights
		w.Objects = sc.Objects
		worlds[i] = w
	}
	canvas := render.Threaded(sc.Camera, worlds)
	writeImagePNG(canvas, viper.GetString("scene")+".png")
}

func depthOfField() {

	worlds := make([]mat.World, config.Cfg.Threads, config.Cfg.Threads)
	sc := scene.DoF()
	for i := 0; i < config.Cfg.Threads; i++ {
		w := mat.NewWorld()
		sc := scene.DoF()
		w.Light = sc.Lights
		w.AreaLight = sc.AreaLights
		w.Objects = sc.Objects
		worlds[i] = w
	}
	canvas := render.Threaded(sc.Camera, worlds)
	writeImagePNG(canvas, viper.GetString("scene")+".png")
}

func softshadows() {

	worlds := make([]mat.World, config.Cfg.Threads, config.Cfg.Threads)
	sc := scene.Softshadows()
	for i := 0; i < config.Cfg.Threads; i++ {
		w := mat.NewWorld()
		sc := scene.Softshadows()
		w.Light = sc.Lights
		w.AreaLight = sc.AreaLights
		w.Objects = sc.Objects
		worlds[i] = w
	}
	canvas := render.Threaded(sc.Camera, worlds)
	writeImagePNG(canvas, viper.GetString("scene")+".png")
}

func csg() {

	camera := mat.NewCamera(640, 480, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-4, 2, -5), mat.NewPoint(0, 0, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(camera.Transform)
	worlds := make([]mat.World, config.Cfg.Threads)
	for i := 0; i < config.Cfg.Threads; i++ {
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

	writeImagePNG(render.Threaded(camera, worlds), viper.GetString("scene")+".png")
}

func withDragonModel() {
	worlds := make([]mat.World, config.Cfg.Threads, config.Cfg.Threads)
	sc := scene.Dragon()
	for i := 0; i < config.Cfg.Threads; i++ {
		w := mat.NewWorld()
		sc := scene.Dragon()
		w.Light = sc.Lights
		w.AreaLight = sc.AreaLights
		w.Objects = sc.Objects
		worlds[i] = w
	}
	canvas := render.Threaded(sc.Camera, worlds)
	writeImagePNG(canvas, viper.GetString("scene")+".png")
}

func withSimpleGopherModel() {

	worlds := make([]mat.World, config.Cfg.Threads, config.Cfg.Threads)
	sc := scene.SimpleGopher()
	for i := 0; i < config.Cfg.Threads; i++ {
		w := mat.NewWorld()
		sc := scene.SimpleGopher()
		w.Light = sc.Lights
		w.AreaLight = sc.AreaLights
		w.Objects = sc.Objects
		worlds[i] = w
	}
	canvas := render.Threaded(sc.Camera, worlds)
	writeImagePNG(canvas, viper.GetString("scene")+".png")
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
	writeImagePNG(canvas, viper.GetString("scene")+".png")
}

// This is my "reference image", used to benchmark the impl. in either 640x480 or 1920x1080
func worldWithPlane() {
	camera := mat.NewCamera(config.Cfg.Width, config.Cfg.Height, math.Pi/3) // -4 är referens!
	viewTransform := mat.ViewTransform(mat.NewPoint(-2, 2.0, -4), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	light := mat.NewLight(mat.NewPoint(-5, 2.5, -3), mat.NewColor(1, 1, 1))

	worlds := make([]mat.World, config.Cfg.Threads)
	for i := 0; i < config.Cfg.Threads; i++ {
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
		mat1 := mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.1, 0.1), 0.1, 0.75, 0.8, 220.0, 0.4)
		s1.SetMaterial(mat1)
		gr.AddChild(s1)

		s2 := mat.NewSphere()
		//s2.CastShadow = false
		s2.SetTransform(mat.Multiply(mat.Translate(-1, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
		mat2 := mat.NewMaterialWithReflectivity(mat.NewColor(0.1, 1.0, 0.1), 0.1, 0.75, 0.8, 220.0, 0.4)
		s2.SetMaterial(mat2)
		gr.AddChild(s2)

		s3 := mat.NewSphere()
		s3.SetTransform(mat.Multiply(mat.Translate(0, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
		mat3 := mat.NewMaterialWithReflectivity(mat.NewColor(0.1, 0.1, 1), 0.1, 0.75, 0.8, 220.0, 0.4)
		s3.SetMaterial(mat3)
		gr.AddChild(s3)

		gr.SetTransform(mat.RotateY(0.67))
		gr.Bounds() // For now, important to always call Bounds on Group once set up.
		mat.Divide(gr, 1)

		w.Objects = append(w.Objects, gr)
		//w.Objects = append(w.Objects, gr.BoundsToCube())

		//cb := mat.NewCube()
		//w.Objects = append(w.Objects, cb)

		worlds[i] = w
	}

	canvas := render.Threaded(camera, worlds)

	// One can use this to render a unit-length XYZ axises superimposed on the image
	//helper.RenderReferenceAxises(canvas, camera)

	// write
	writeImagePNG(canvas, viper.GetString("scene")+".png")
}

// This is my "reference image", used to benchmark the impl. in either 640x480 or 1920x1080
func threaded() {
	camera := mat.NewCamera(config.Cfg.Width, config.Cfg.Height, math.Pi/3) // -4 är referens!
	viewTransform := mat.ViewTransform(mat.NewPoint(-2, 2.0, -4), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
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
		//w.Objects = append(w.Objects, gr.BoundsToCube())

		//cb := mat.NewCube()
		//w.Objects = append(w.Objects, cb)

		worlds[i] = w
	}
	var canvas *mat.Canvas
	for i := 1; i < 9; i++ {
		config.Cfg.Threads = i
		st := time.Now()
		canvas = render.Threaded(camera, worlds)
		fmt.Printf("%v\n", time.Since(st).Milliseconds())
	}

	// One can use this to render a unit-length XYZ axises superimposed on the image
	helper.RenderReferenceAxises(canvas, camera)

	// write
	writeImagePNG(canvas, viper.GetString("scene")+".png")
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
	fmt.Printf("writing output to file %v\n", filename)
	myImage := image.NewRGBA(image.Rect(0, 0, canvas.W, canvas.H))
	writeDataToPNG(canvas, myImage)
	outputFile, _ := os.Create(filename)
	defer outputFile.Close()
	_ = png.Encode(outputFile, myImage)
}
