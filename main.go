package main

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/canvas"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"io/ioutil"
	"os"
)

func main() {
	projectileDemo()
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
		fmt.Printf("Drawing at: %d %d\n", int(prj.pos.Get(0)), c.H - int(prj.pos.Get(1)))
		c.WritePixel(int(prj.pos.Get(0)), c.H - int(prj.pos.Get(1)), red)
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
	wind *mat.Tuple4
}

func NewEnvironment(gravity *mat.Tuple4, wind *mat.Tuple4) *Environment {
	return &Environment{gravity: gravity, wind: wind}
}

type Projectile struct {
	pos *mat.Tuple4
	velocity *mat.Tuple4
}

func NewProjectile(pos *mat.Tuple4, velocity *mat.Tuple4) *Projectile {
	return &Projectile{pos: pos, velocity: velocity}
}