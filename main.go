package main

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/eriklupander/rt/internal/pkg/obj"
	"github.com/eriklupander/rt/internal/pkg/parser"
	"io/ioutil"
	"math"
	_ "net/http/pprof"
	"os"
)

func main() {
	// we need a webserver to get the pprof webserver
	//go func() {
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	//parse()
	//csg()
	//withModel()
	//groups()
	worldWithPlane()
	//renderworld()
	//shadedSphereDemo()
	//circleDemo()
	//clockDemo()
	//projectileDemo()

	//termChan := make(chan os.Signal)
	//signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	//<-termChan // Blocks here!!
	//fmt.Println("shutting down!")
}

var white = mat.NewColor(1, 1, 1)
var black = mat.NewColor(0, 0, 0)

func parse() {
	scene := parser.ParseYAML("scenes/simple.yaml")
	fmt.Printf("%v", scene)

	w := scene.World
	w.Light = scene.Lights
	camera := scene.Camera
	canvas := mat.RenderThreaded(*camera, *w)
	// writec
	data := canvas.ToPPM()
	err := ioutil.WriteFile("fromyaml.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func csg() {
	w := mat.NewWorld()
	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(0, 2, -2), mat.NewColor(1, 1, 1)))
	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(0, 3, 0), mat.NewColor(1, 1, 1)))

	camera := mat.NewCamera(640, 480, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-4, 2, -5), mat.NewPoint(0, 0, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform

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
	//floor.Material.Pattern = mat.NewCheckerPattern(white, black)
	w.Objects = append(w.Objects, floor)

	canvas := mat.RenderThreaded(camera, w)
	// writec
	data := canvas.ToPPM()
	err := ioutil.WriteFile("csg1.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func withModel() {

	bytes, err := ioutil.ReadFile("Toilet.1.obj")
	if err != nil {
		panic(err.Error())
	}

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

	camera := mat.NewCamera(96, 72, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-4.3, 5, -8), mat.NewPoint(0, 2.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform

	canvas := mat.RenderThreaded(camera, w)

	mat.RenderReferenceAxises(canvas, camera)

	// writec
	data := canvas.ToPPM()
	err = ioutil.WriteFile("toilet.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func worldWithPlane() {
	w := mat.NewWorld()
	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(-3, 2.5, -3), mat.NewColor(1, 1, 1)))

	camera := mat.NewCamera(640, 480, math.Pi/3)
	//camera := mat.NewCamera(320, 240, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-2, 1.0, -4), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	floor := mat.NewPlane()
	floor.SetMaterial(mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.5, 0.5), 0.1, 0.9, 0.7, 240, 0.1))
	floor.Material.Pattern = mat.NewCheckerPattern(white, black)
	floor.Material.Pattern.SetPatternTransform(mat.Scale(0.5, 0.5, 0.5))
	w.Objects = append(w.Objects, floor)

	wall := mat.NewPlane()
	wall.SetMaterial(mat.NewMaterial(mat.NewColor(0.9, 0.9, 0.9), 0.1, 0.9, 0.7, 200))
	wall.SetTransform(mat.Translate(0, 0, 8))
	wall.SetTransform(mat.RotateX(math.Pi / 2))
	w.Objects = append(w.Objects, wall)

	// transparent sphere
	middle := mat.NewSphere()
	middle.SetTransform(mat.Translate(-0.5, 1, 0.5))
	middle.Material = mat.NewDefaultReflectiveMaterial(0.3)
	middle.Material.Color = mat.NewColor(0.1, 0.1, 0.1)
	middle.Material.Diffuse = 0.7
	middle.Material.Specular = 0.6
	middle.Material.Transparency = 0.95
	middle.Material.RefractiveIndex = 1.5
	w.Objects = append(w.Objects, middle)

	// back sphere
	right := mat.NewSphere()
	right.SetTransform(mat.Multiply(mat.Translate(-0.75, 0.5, 2.5), mat.Scale(0.5, 0.5, 0.5)))
	right.Material = mat.NewDefaultMaterial()
	right.Material.Color = mat.NewColor(1, 0, 0)
	right.Material.Diffuse = 0.7
	right.Material.Specular = 0.3
	right.Material.Reflectivity = 0.3
	w.Objects = append(w.Objects, right)

	// cube
	cube := mat.NewCube()
	cube.SetTransform(mat.Multiply(mat.Translate(-1.6, 0.25, 1.5), mat.Scale(0.25, 0.25, 0.25)))
	cube.Material = mat.NewDefaultMaterial()
	cube.Material.Color = mat.NewColor(1, 0.6, 0.2)
	cube.Material.Transparency = 0.0
	cube.Material.Diffuse = 0.7
	cube.Material.Specular = 0.3
	cube.Material.Reflectivity = 0.0
	w.Objects = append(w.Objects, cube)

	//  Cylinder
	cyl := mat.NewCylinderMMC(0.0, 3.0, true)
	cyl.SetTransform(mat.Translate(-1.6, 0.5, 1.5))
	cyl.SetTransform(mat.Scale(0.2, 0.2, 0.2))
	m := mat.NewDefaultReflectiveMaterial(0.6)
	m.Color = mat.NewColor(0.7, 0.5, 1.0)
	cyl.SetMaterial(m)
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
	mat4 := mat.NewMaterialWithReflectivity(mat.NewColor(0.6, 0.6, 0.6), 0.1, 0.5, 0.8, 220.0, 0.4)
	s4.SetMaterial(mat4)
	gr.AddChild(s4)

	w.Objects = append(w.Objects, gr)

	canvas := mat.Render(camera, w)
	mat.RenderReferenceAxises(canvas, camera)

	// write
	data := canvas.ToPPM()
	err := ioutil.WriteFile("world-transparency.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func groups() {
	w := mat.NewWorld()
	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(-3, 2.5, -3), mat.NewColor(1, 1, 1)))

	camera := mat.NewCamera(640, 480, math.Pi/3)
	//camera := mat.NewCamera(320, 240, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-1.3, 2, -5), mat.NewPoint(0, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform
	camera.Inverse = mat.Inverse(viewTransform)

	gr := mat.NewGroup()

	s1 := mat.NewSphere()
	//s1.SetTransform(mat.Multiply(mat.Translate(-2, 0.25, -1), mat.Scale(1.25, 0.25, 0.25)))
	s1.SetTransform(mat.Translate(-2, -1, 0))
	mat1 := mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.1, 0.1), 0.1, 0.5, 0.8, 220.0, 0.4)
	s1.SetMaterial(mat1)
	gr.AddChild(s1)

	s4 := mat.NewSphere()
	//s4.SetTransform(mat.Multiply(mat.Translate(1, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
	s4.SetTransform(mat.Translate(0, 0, 0))
	mat4 := mat.NewMaterialWithReflectivity(mat.NewColor(0.1, 0.1, 1.0), 0.1, 0.5, 0.8, 220.0, 0.4)
	s4.SetMaterial(mat4)
	gr.AddChild(s4)

	w.Objects = append(w.Objects, gr)

	canvas := mat.Render(camera, w)
	mat.RenderReferenceAxises(canvas, camera)

	// write
	data := canvas.ToPPM()
	err := ioutil.WriteFile("world-group.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func renderworld() {
	w := mat.NewWorld()
	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(-10, 1, -10), mat.NewColor(1, 1, 1)))
	w.Light = append(w.Light, mat.NewLight(mat.NewPoint(1, 13, 1), mat.NewColor(0.5, 0.5, 0.5)))

	camera := mat.NewCamera(480, 320, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(0, 1.5, -5), mat.NewPoint(0, 1, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform

	// Create floor
	floor := mat.NewSphere()
	floor.Transform = mat.Scale(10, 0.01, 10)
	floor.Material = mat.NewDefaultMaterial()
	floor.Material.Color = mat.NewColor(1, 0.9, 0.9)
	floor.Material.Specular = 0.0
	floor.Material.Reflectivity = 0.2
	w.Objects = append(w.Objects, floor)

	// create left wall
	leftWall := mat.NewSphere()

	scaleM := mat.Scale(10, 0.01, 10)
	rotXM := mat.RotateX(math.Pi / 2)
	rotYM := mat.RotateY(-math.Pi / 4)
	transM := mat.Translate(0, 0, 5)

	m1 := mat.Multiply(transM, rotYM)
	m2 := mat.Multiply(m1, rotXM)
	m3 := mat.Multiply(m2, scaleM)
	leftWall.Transform = m3
	leftWall.Material = floor.Material
	w.Objects = append(w.Objects, leftWall)

	// create right wall
	rightWall := mat.NewSphere()

	scaleM = mat.Scale(10, 0.01, 10)
	rotXM = mat.RotateX(math.Pi / 2)
	rotYM = mat.RotateY(math.Pi / 4)
	transM = mat.Translate(0, 0, 5)

	m1 = mat.Multiply(transM, rotYM)
	m2 = mat.Multiply(m1, rotXM)
	m3 = mat.Multiply(m2, scaleM)
	rightWall.Transform = m3
	rightWall.Material = floor.Material
	w.Objects = append(w.Objects, rightWall)

	// middle sphere
	middle := mat.NewSphere()
	middle.Transform = mat.Translate(-0.5, 1, 0.5)
	middle.Material = mat.NewDefaultMaterial()
	middle.Material.Color = mat.NewColor(0.1, 1, 0.5)
	middle.Material.Diffuse = 0.7
	middle.Material.Specular = 0.3
	w.Objects = append(w.Objects, middle)

	// right sphere
	right := mat.NewSphere()
	right.Transform = mat.Multiply(mat.Translate(1.5, 0.5, -0.5), mat.Scale(0.5, 0.5, 0.5))
	right.Material = mat.NewDefaultMaterial()
	right.Material.Color = mat.NewColor(0.5, 1, 0.1)
	right.Material.Diffuse = 0.7
	right.Material.Specular = 0.3
	w.Objects = append(w.Objects, right)

	// left sphere
	left := mat.NewSphere()
	left.Transform = mat.Multiply(mat.Translate(-1.5, 0.33, -0.75), mat.Scale(0.33, 0.33, 0.33))
	left.Material = mat.NewDefaultMaterial()
	left.Material.Color = mat.NewColor(1, 0.8, 0.1)
	left.Material.Diffuse = 0.7
	left.Material.Specular = 0.3
	w.Objects = append(w.Objects, left)

	// cube
	cube := mat.NewCube()
	cube.Transform = mat.Multiply(mat.Translate(-.6, 0.25, -1.5), mat.Scale(0.25, 0.25, 0.25))
	cube.Material = mat.NewDefaultMaterial()
	cube.Material.Color = mat.NewColor(1, 0.6, 0.2)
	cube.Material.Transparency = 0.0
	cube.Material.Diffuse = 0.7
	cube.Material.Specular = 0.3
	cube.Material.Reflectivity = 0.0
	w.Objects = append(w.Objects, cube)

	canvas := mat.RenderThreaded(camera, w)
	// write
	data := canvas.ToPPM()
	err := ioutil.WriteFile("world1.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func shadedSphereDemo() {
	c := mat.NewCanvas(512, 512)

	// this is our eye starting 15 units "in front" of origo.
	rayOrigin := mat.NewPoint(0, 0, -15.0)

	// Note!! If I understand this correctly, the "wall" in this case is actually an abstraction
	// of the "far" wall of the view frustum, giving something to cast our ray against, forming a vector between
	// the eye and a point in world space.
	wallZ := 20.0
	wallSize := 7.0
	pixelSize := wallSize / float64(c.W)
	half := wallSize / 2
	sphere := mat.NewSphere()
	//mat.SetTransform(&sphere, mat.Translate(1, 1, 1))
	material := mat.NewDefaultMaterial()
	material.Color = mat.NewColor(1, 0.2, 1)
	sphere.SetMaterial(material)

	lightPos := mat.NewPoint(-10, 10, -10)
	lightColor := mat.NewColor(1, 1, 1)
	light := mat.NewLight(lightPos, lightColor)

	for row := 0; row < c.W; row++ {
		worldY := half - pixelSize*float64(row)

		for col := 0; col < c.H; col++ {
			worldX := -half + pixelSize*float64(col)
			posOnWall := mat.NewPoint(worldX, worldY, wallZ)

			// Build a ray (origin + direction)
			rayFromOriginToPosOnWall := mat.NewRay(rayOrigin, mat.Normalize(mat.Sub(posOnWall, rayOrigin)))

			// check if our ray intersects the sphere
			intersections := mat.IntersectRayWithShape(sphere, rayFromOriginToPosOnWall)
			intersection, found := mat.Hit(intersections)

			if found {
				pointOfHit := mat.Position(rayFromOriginToPosOnWall, intersection.T)
				normalAtHit := mat.NormalAt(sphere, pointOfHit, nil)
				minusEyeRayVector := mat.Negate(rayFromOriginToPosOnWall.Direction)
				color := mat.Lighting(sphere.Material, sphere, light, pointOfHit, minusEyeRayVector, normalAtHit, false)

				c.WritePixel(col, c.H-row, color)
			}
		}
	}
	// write
	data := c.ToPPM()
	err := ioutil.WriteFile("shadedcircle.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func circleDemo() {
	c := mat.NewCanvas(100, 100)

	rayOrigin := mat.NewPoint(0, 0, -15.0)
	wallZ := 20.0
	wallSize := 7.0
	pixelSize := wallSize / float64(c.W)
	half := wallSize / 2
	color := mat.NewColor(1, 0, 0)
	sphere := mat.NewSphere()

	//mat.SetTransform(sphere, mat.Scale(1, 0.5, 1))
	//mat.SetTransform(sphere, mat.Multiply(mat.RotateZ(math.Pi/4), mat.Scale(0.5, 1, 1)))

	for row := 0; row < c.W; row++ {
		worldY := half - pixelSize*float64(row)

		for col := 0; col < c.H; col++ {
			worldX := -half + pixelSize*float64(col)
			posOnWall := mat.NewPoint(worldX, worldY, wallZ)

			rayFromOriginToPosOnWall := mat.NewRay(rayOrigin, mat.Normalize(mat.Sub(posOnWall, rayOrigin)))

			// check if our ray intersects the sphere
			intersections := mat.IntersectRayWithShape(sphere, rayFromOriginToPosOnWall)
			_, found := mat.Hit(intersections)
			if found {
				c.WritePixel(col, c.H-row, color)
			}
		}
	}
	// write
	data := c.ToPPM()
	err := ioutil.WriteFile("circle.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

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
