package main

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/canvas"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/eriklupander/rt/internal/pkg/ray"
	"io/ioutil"
	"math"
	"os"
)

func main() {
	circleDemo()
	//clockDemo()
	//projectileDemo()
}

func circleDemo() {
	c := canvas.NewCanvas(100, 100)

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

			rayFromOriginToPosOnWall := ray.New(rayOrigin, mat.Normalize(*mat.Sub(*posOnWall, *rayOrigin)))

			// check if our ray intersects the sphere
			intersections := ray.IntersectRayWithSphere(sphere, rayFromOriginToPosOnWall)
			intersection := ray.Hit(intersections)
			if intersection != nil {
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
	c := canvas.NewCanvas(80, 80)
	center := (c.W/2 + c.H/2) / 2
	white := mat.NewColor(1, 1, 1)

	point := mat.NewPoint(0, 1, 0)
	for i := 0; i < 12; i++ {
		rotation := float64(i) * (2 * math.Pi) / 12
		rotMat := mat.RotateZ(rotation)
		p2 := mat.MultiplyByTuple(*rotMat, *point)
		p2 = mat.MultiplyByScalar(*p2, 30.0)
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
	prj := NewProjectile(mat.NewPoint(0, 1, 0), mat.MultiplyByScalar(*mat.Normalize(*mat.NewVector(1, 1.8, 0)), 11.25))
	env := NewEnvironment(mat.NewVector(0, -0.1, 0), mat.NewVector(-0.01, 0, 0))
	c := canvas.NewCanvas(900, 550)
	red := mat.NewColor(1, 1, 1)
	for prj.pos.Get(1) > 0.0 {
		tick(prj, env)
		//time.Sleep(time.Millisecond * 100)
		fmt.Printf("Projectile pos %v at height %v with velocity %v\n", mat.Magnitude(*prj.pos), prj.pos.Get(1), *prj.velocity)
		fmt.Printf("Drawing at: %d %d\n", int(prj.pos.Get(0)), c.H-int(prj.pos.Get(1)))
		c.WritePixel(int(prj.pos.Get(0)), c.H-int(prj.pos.Get(1)), red)
	}
	fmt.Printf("Projectile flew %v\n", mat.Magnitude(*prj.pos))
	data := c.ToPPM()
	err := ioutil.WriteFile("pic.ppm", []byte(data), os.FileMode(0755))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func tick(prj *Projectile, env *Environment) {
	prj.pos = mat.Add(*prj.pos, *prj.velocity)
	prj.velocity = mat.Add(*prj.velocity, *env.gravity)
	prj.velocity = mat.Add(*prj.velocity, *env.wind)
}

type Environment struct {
	gravity *mat.Tuple4
	wind    *mat.Tuple4
}

func NewEnvironment(gravity *mat.Tuple4, wind *mat.Tuple4) *Environment {
	return &Environment{gravity: gravity, wind: wind}
}

type Projectile struct {
	pos      *mat.Tuple4
	velocity *mat.Tuple4
}

func NewProjectile(pos *mat.Tuple4, velocity *mat.Tuple4) *Projectile {
	return &Projectile{pos: pos, velocity: velocity}
}
