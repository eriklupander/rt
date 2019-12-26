package render

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/inhies/go-bytesize"
	"math"
	"runtime"
	"sync"
	"time"
)

// constants
var originPoint = mat.NewPoint(0, 0, 0)
var black = mat.NewColor(0, 0, 0)

func New(world mat.World) Context {
	cStack := make([]ShadeData, 512)
	for i := 0; i < 512; i++ {
		cStack[i] = NewShadeData()
	}

	return Context{
		world: world,
		total: 0,

		// allocate memory
		pointInView: mat.NewPoint(0, 0, -1.0),
		pixel:       mat.NewColor(0, 0, 0),
		origin:      mat.NewPoint(0, 0, 0),
		direction:   mat.NewVector(0, 0, 0),
		subVec:      mat.NewVector(0, 0, 0),

		// allocate ray
		firstRay: mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0)),

		// stack for shading
		cStack: cStack,
	}
}

type Context struct {
	Id     int
	world  mat.World
	camera mat.Camera
	canvas *mat.Canvas
	jobs   chan *job
	wg     *sync.WaitGroup
	total  int

	// pixel cache
	pointInView mat.Tuple4
	pixel       mat.Tuple4
	origin      mat.Tuple4
	direction   mat.Tuple4
	subVec      mat.Tuple4

	// ray cache
	firstRay mat.Ray

	// each renderContext needs to pre-allocate 1+(MAX_REFLECTION*MAX_REFRACTION), e.g. 26. Use them linearly.
	cStack []ShadeData
}

func NewContext(id int, world mat.World, camera mat.Camera, canvas *mat.Canvas, jobs chan *job, wg *sync.WaitGroup) Context {
	cStack := make([]ShadeData, 512)
	for i := 0; i < 512; i++ {
		cStack[i] = NewShadeData()
	}

	return Context{
		Id:     id,
		world:  world,
		camera: camera,
		canvas: canvas,
		jobs:   jobs,
		wg:     wg,
		total:  0,

		// allocate memory
		pointInView: mat.NewPoint(0, 0, -1.0),
		pixel:       mat.NewColor(0, 0, 0),
		origin:      mat.NewPoint(0, 0, 0),
		direction:   mat.NewVector(0, 0, 0),
		subVec:      mat.NewVector(0, 0, 0),

		// allocate ray
		firstRay: mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0)),

		// stack for shading
		cStack: cStack,
	}
}

func Threaded(c mat.Camera, worlds []mat.World) *mat.Canvas {
	st := time.Now()
	canvas := mat.NewCanvas(c.Width, c.Height)
	jobs := make(chan *job)

	wg := sync.WaitGroup{}
	wg.Add(canvas.W * canvas.H)

	// allocate GOMAXPROCS render Contexts
	var GOMAXPROCS = 8
	renderContexts := make([]Context, GOMAXPROCS)
	for i := 0; i < GOMAXPROCS; i++ {
		renderContexts[i] = NewContext(i, worlds[i], c, canvas, jobs, &wg)
	}

	// start workers
	for i := 0; i < GOMAXPROCS; i++ {
		go renderContexts[i].workerFuncPerPixel()
	}

	// start passing work to the workers, one pixel at a time
	for row := 0; row < c.Height; row++ {
		for col := 0; col < c.Width; col++ {
			jobs <- &job{row: row, col: col}
		}
		fmt.Printf("%d/%d\n", row, c.Height)
	}

	wg.Wait()
	fmt.Println("All done")
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	fmt.Printf("Memory: %v ", bytesize.New(float64(stats.Alloc)).String())
	fmt.Printf("Mallocs: %v ", stats.Mallocs)
	fmt.Printf("Total alloc: %v\n", bytesize.New(float64(stats.TotalAlloc)).String())
	fmt.Printf("%v", time.Now().Sub(st))
	return canvas
}

func (rc *Context) workerFuncPerPixel() {
	for job := range rc.jobs {
		rc.renderPixel(job)
	}
}

func (rc *Context) renderPixel(job *job) {
	for i := 0; i < 64; i++ {
		rc.cStack[i].WorldXS = rc.cStack[i].WorldXS[:0]
		rc.cStack[i].ShadowXS = rc.cStack[i].ShadowXS[:0]
	}
	rc.total = 0
	rc.rayForPixel(job.col, job.row, &rc.firstRay)
	color := rc.colorAt(rc.firstRay, 5)
	//if rc.Id == 0 {
	//	fmt.Printf("finished color at %d %d, total: %d\n", job.col, job.row, rc.total)
	//}
	rc.canvas.WritePixelMutex(job.col, job.row, color)
	rc.wg.Done()
	//fmt.Printf("Thread %d remain: %d\n", rc.Id, rc.fakeremain)
}

func (rc *Context) rayForPixel(x, y int, out *mat.Ray) {

	xOffset := rc.camera.PixelSize * (float64(x) + 0.5)
	yOffset := rc.camera.PixelSize * (float64(y) + 0.5)

	// this feels a little hacky but actually works.
	worldX := rc.camera.HalfWidth - xOffset
	worldY := rc.camera.HalfHeight - yOffset

	// mat.NewPoint(worldX, worldY, -1.0)
	rc.pointInView.Elems[0] = worldX
	rc.pointInView.Elems[1] = worldY

	mat.MultiplyByTuplePtr(rc.camera.Inverse, rc.pointInView, &rc.pixel)
	mat.MultiplyByTuplePtr(rc.camera.Inverse, originPoint, &rc.origin)
	mat.SubPtr(rc.pixel, rc.origin, &rc.subVec)
	mat.NormalizePtr(rc.subVec, &rc.direction)

	out.Direction = rc.direction
	out.Origin = rc.origin
}

func (rc *Context) colorAt(r mat.Ray, remaining int) mat.Tuple4 {
	rc.total++

	rc.cStack[rc.total].WorldXS = mat.IntersectWithWorldPtr(rc.world, r, rc.cStack[rc.total].WorldXS, &rc.cStack[rc.total].ShadowInRay)
	if len(rc.cStack[rc.total].WorldXS) > 0 {
		mat.PrepareComputationForIntersectionPtr(rc.cStack[rc.total].WorldXS[0], r, &rc.cStack[rc.total].Comps, rc.cStack[rc.total].WorldXS...)
		return rc.shadeHit(rc.cStack[rc.total].Comps, remaining)
	} else {
		return black
	}
}

func (rc *Context) reflectedColor(comps mat.Computation, remaining int) mat.Tuple4 {
	if remaining <= 0 || comps.Object.GetMaterial().Reflectivity == 0.0 {
		return black
	}
	reflectRay := mat.NewRay(comps.OverPoint, comps.ReflectVec)
	remaining--
	reflectedColor := rc.colorAt(reflectRay, remaining)
	return mat.MultiplyByScalar(reflectedColor, comps.Object.GetMaterial().Reflectivity)
}

func (rc *Context) refractedColor(comps mat.Computation, remaining int) mat.Tuple4 {
	if remaining <= 0 || comps.Object.GetMaterial().Transparency == 0.0 {
		return black
	}

	// Find the ratio of first index of refraction to the second.
	nRatio := comps.N1 / comps.N2

	// cos(theta_i) is the same as the dot product of the two vectors
	cosI := mat.Dot(comps.EyeVec, comps.NormalVec)

	// Find sin(theta_t)^2 via trigonometric identity
	sin2Theta := (nRatio * nRatio) * (1.0 - (cosI * cosI))
	if sin2Theta > 1.0 {
		return black
	}

	// Find cos(theta_t) via trigonometric identity
	cosTheta := math.Sqrt(1.0 - sin2Theta)

	// Compute the direction of the refracted ray
	direction := mat.Sub(
		mat.MultiplyByScalar(comps.NormalVec, (nRatio*cosI)-cosTheta),
		mat.MultiplyByScalar(comps.EyeVec, nRatio))

	// Create the refracted ray
	refractRay := mat.NewRay(comps.UnderPoint, direction)

	// Find the color of the refracted ray, making sure to multiply
	// by the transparency value to account for any opacity
	remaining--
	color := mat.MultiplyByScalar(rc.colorAt(refractRay, remaining), comps.Object.GetMaterial().Transparency)

	return color
}

func (rc *Context) shadeHit(comps mat.Computation, remaining int) mat.Tuple4 {
	var surfaceColor = mat.NewColor(0, 0, 0)
	for _, light := range rc.world.Light {
		inShadow := rc.pointInShadow(light, comps.OverPoint)
		color := mat.Lighting(comps.Object.GetMaterial(), comps.Object, light, comps.OverPoint, comps.EyeVec, comps.NormalVec, inShadow)
		surfaceColor = mat.Add(surfaceColor, color)
	}
	reflectedColor := rc.reflectedColor(comps, remaining)
	refractedColor := rc.refractedColor(comps, remaining)

	material := comps.Object.GetMaterial()
	if material.Reflectivity > 0.0 && material.Transparency > 0.0 {
		reflectance := mat.Schlick(comps)
		return mat.Add(mat.Add(surfaceColor, mat.MultiplyByScalar(reflectedColor, reflectance)), mat.MultiplyByScalar(refractedColor, 1-reflectance))
	} else {
		return mat.Add(surfaceColor, mat.Add(reflectedColor, refractedColor))
	}
}

func (rc *Context) pointInShadow(light mat.Light, p mat.Tuple4) bool {

	vecToLight := mat.Sub(light.Position, p)
	distance := mat.Magnitude(vecToLight)

	ray := mat.NewRay(p, mat.Normalize(vecToLight))

	// use stack...
	rc.cStack[rc.total].ShadowXS = mat.IntersectWithWorldPtr(rc.world, ray, rc.cStack[rc.total].ShadowXS, &rc.cStack[rc.total].InRay)
	if len(rc.cStack[rc.total].ShadowXS) > 0 {
		for _, x := range rc.cStack[rc.total].ShadowXS {
			if x.T < distance {
				return true
			}
		}
	}
	return false
}

type job struct {
	row int
	col int
}

// ShadeData should contain pre-allocated memory for each "colorAt" recursion
type ShadeData struct {
	WorldXS     []mat.Intersection
	ShadowXS    []mat.Intersection
	InRay       mat.Ray
	ShadowInRay mat.Ray

	Comps mat.Computation
}

func NewShadeData() ShadeData {
	worldXS := make([]mat.Intersection, 8)
	shadowXS := make([]mat.Intersection, 8)

	worldXS = worldXS[:0]
	shadowXS = shadowXS[:0]

	return ShadeData{
		WorldXS:     worldXS,
		ShadowXS:    shadowXS,
		InRay:       mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0)),
		ShadowInRay: mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0)),
		Comps:       mat.NewComputation(),
	}
}
