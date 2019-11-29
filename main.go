package main

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"io/ioutil"
	"math"
	"os"
)

func main() {
	worldWithPlane()
	//renderworld()
	//shadedSphereDemo()
	//circleDemo()
	//clockDemo()
	//projectileDemo()
}

var white = mat.NewColor(1, 1, 1)
var black = mat.NewColor(0, 0, 0)

func worldWithPlane() {
	w := mat.NewWorld()
	w.Light = mat.NewLight(mat.NewPoint(-10, 5, -10), mat.NewColor(1, 1, 1))

	camera := mat.NewCamera(2540, 1600, math.Pi/3)
	//camera := mat.NewCamera(320, 240, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(-2.1, 1, -4.5), mat.NewPoint(-1, 0.5, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform

	floor := mat.NewPlane()
	floor.SetMaterial(mat.NewMaterialWithReflectivity(mat.NewColor(1, 0.5, 0.5), 0.1, 0.9, 0.7, 200, 0.1))
	floor.Material.Pattern = mat.NewCheckerPattern(white, black)
	w.Objects = append(w.Objects, floor)

	wall := mat.NewPlane()
	wall.SetTransform(mat.Translate(0, 0, 10))
	wall.SetTransform(mat.RotateX(math.Pi / 2))
	material := mat.NewDefaultReflectiveMaterial(1.0)
	material.Transparency = 0.0
	material.Color = mat.NewColor(0.05, 0.05, 0.05)

	wall.SetMaterial(material)
	w.Objects = append(w.Objects, wall)

	rightWall := mat.NewPlane()
	rightWall.SetTransform(mat.Translate(5, 0, 0))
	rightWall.SetTransform(mat.RotateY(math.Pi / 2))
	rightWall.SetTransform(mat.RotateX(math.Pi / 2))
	w.Objects = append(w.Objects, rightWall)

	// middle sphere
	//pattern := mat.NewRingPattern(mat.NewColor(1, 0, 0), mat.NewColor(0, 0, 1))
	//pattern.SetPatternTransform(mat.Multiply(pattern.Transform, mat.RotateZ(math.Pi/2)))
	middle := mat.NewSphere()
	middle.Transform = mat.Translate(-0.5, 1, 0.5)
	middle.Material = mat.NewDefaultReflectiveMaterial(0.3)
	middle.Material.Color = mat.NewColor(0.1, 0.1, 0.1)
	middle.Material.Diffuse = 0.7
	middle.Material.Specular = 0.6
	middle.Material.Transparency = 0.95
	middle.Material.RefractiveIndex = 1.5
	//middle.Material.Pattern = pattern
	w.Objects = append(w.Objects, middle)

	// right sphere
	right := mat.NewSphere()
	right.Transform = mat.Multiply(mat.Translate(-0.75, 0.5, 2.5), mat.Scale(0.5, 0.5, 0.5))
	right.Material = mat.NewDefaultMaterial()
	right.Material.Color = mat.NewColor(1, 0, 0)
	right.Material.Diffuse = 0.7
	right.Material.Specular = 0.3
	right.Material.Reflectivity = 0.3
	w.Objects = append(w.Objects, right)

	// cube
	cube := mat.NewCube()
	cube.Transform = mat.Multiply(mat.Translate(-2.6, 0.25, -1.5), mat.Scale(0.25, 0.25, 0.25))
	cube.Material = mat.NewDefaultMaterial()
	cube.Material.Color = mat.NewColor(1, 0.6, 0.2)
	cube.Material.Transparency = 0.0
	cube.Material.Diffuse = 0.7
	cube.Material.Specular = 0.3
	cube.Material.Reflectivity = 0.0
	w.Objects = append(w.Objects, cube)

	// left sphere
	left := mat.NewSphere()
	left.Transform = mat.Multiply(mat.Translate(-2, 0.33, -1.0), mat.Scale(0.33, 0.33, 0.33))
	left.Material = mat.NewDefaultMaterial()
	left.Material.Color = mat.NewColor(1, 0.8, 0.1)
	left.Material.Diffuse = 0.7
	left.Material.Specular = 0.3
	left.Material.Reflectivity = 0.1
	w.Objects = append(w.Objects, left)

	canvas := mat.RenderThreaded(camera, w)

	//mat.RenderReferenceAxises(canvas, camera)

	// write
	data := canvas.ToPPM()
	err := ioutil.WriteFile("world-transparency.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func renderworld() {
	w := mat.NewWorld()
	w.Light = mat.NewLight(mat.NewPoint(-10, 1, -10), mat.NewColor(1, 1, 1))

	camera := mat.NewCamera(1376, 768, math.Pi/3)
	viewTransform := mat.ViewTransform(mat.NewPoint(0, 1.5, -5), mat.NewPoint(0, 1, 0), mat.NewVector(0, 1, 0))
	camera.Transform = viewTransform

	// Create floor
	floor := mat.NewSphere()
	floor.Transform = mat.Scale(10, 0.01, 10)
	floor.Material = mat.NewDefaultMaterial()
	floor.Material.Color = mat.NewColor(1, 0.9, 0.9)
	floor.Material.Specular = 0.0
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

	canvas := mat.Render(camera, w)
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
				normalAtHit := mat.NormalAt(sphere, pointOfHit)
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
