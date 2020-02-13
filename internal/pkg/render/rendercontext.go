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
	cStack := make([]ShadeData, 256)
	for i := 0; i < 256; i++ {
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

func NewContext(id int, world mat.World, camera mat.Camera, canvas *mat.Canvas, jobs chan *job, wg *sync.WaitGroup) Context {
	ctx := New(world)
	ctx.Id = id
	ctx.camera = camera
	ctx.canvas = canvas
	ctx.jobs = jobs
	ctx.wg = wg
	return ctx
}

type Context struct {
	Id     int
	world  mat.World
	camera mat.Camera
	canvas *mat.Canvas
	jobs   chan *job
	wg     *sync.WaitGroup
	total  int
	depth  int

	// pixel cache
	pointInView mat.Tuple4
	pixel       mat.Tuple4
	origin      mat.Tuple4
	direction   mat.Tuple4
	subVec      mat.Tuple4

	// ray cache
	firstRay mat.Ray

	// each renderContext needs to pre-allocate shade-data for sufficient number of recursions
	cStack []ShadeData
}

func Threaded(c mat.Camera, worlds []mat.World) *mat.Canvas {
	st := time.Now()
	canvas := mat.NewCanvas(c.Width, c.Height)
	jobs := make(chan *job)

	wg := sync.WaitGroup{}
	wg.Add(canvas.W * canvas.H) // PIXEL RENDER
	//wg.Add(canvas.H)               // LINE RENDER

	// allocate GOMAXPROCS render Contexts
	var GOMAXPROCS = 8
	renderContexts := make([]Context, GOMAXPROCS)
	for i := 0; i < GOMAXPROCS; i++ {
		renderContexts[i] = NewContext(i, worlds[i], c, canvas, jobs, &wg)
	}

	// start workers
	// Per-pixel worker:
	for i := 0; i < GOMAXPROCS; i++ {
		go renderContexts[i].workerFuncPerPixel()
	}

	// Per line worker:
	//for i := 0; i < GOMAXPROCS; i++ {
	//	go renderContexts[i].workerFuncPerLine()
	//}

	// start passing work to the workers:
	// One pixel at a time
	for row := 0; row < c.Height; row++ {
		for col := 0; col < c.Width; col++ {
			jobs <- &job{row: row, col: col}
		}
		fmt.Printf("%d/%d\n", row, c.Height)
	}

	// Pass by line
	//for row := 0; row < c.Height; row++ {
	//	//for col := 0; col < c.Width; col++ {
	//	jobs <- &job{row: row, col: 0}
	//	//}
	//	fmt.Printf("%d/%d\n", row, c.Height)
	//}

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
func (rc *Context) workerFuncPerLine() {
	for job := range rc.jobs {
		for i := 0; i < 1920; i++ {
			job.col = i
			rc.renderPixel(job)
		}
		rc.wg.Done()
	}
}

func (rc *Context) renderSinglePixel(col, row int) mat.Tuple4 {
	for i := 0; i < 256; i++ {
		rc.cStack[i].WorldXS = rc.cStack[i].WorldXS[:0]
		rc.cStack[i].ShadowXS = rc.cStack[i].ShadowXS[:0]
	}
	rc.total = 0
	rc.depth = 0
	rc.rayForPixel(col, row, &rc.firstRay)
	color := rc.colorAt(rc.firstRay, 5, 5)
	return color
}

func (rc *Context) renderPixel(job *job) {
	for i := 0; i < 256; i++ {
		rc.cStack[i].WorldXS = rc.cStack[i].WorldXS[:0]
		rc.cStack[i].ShadowXS = rc.cStack[i].ShadowXS[:0]
	}
	rc.total = 0
	rc.depth = 0
	rc.rayForPixel(job.col, job.row, &rc.firstRay)
	color := rc.colorAt(rc.firstRay, 5, 5)
	//if rc.Id == 0 {
	//fmt.Printf("finished color at %d %d, total: %d depth: %d\n", job.col, job.row, rc.total, rc.depth)
	//}
	rc.canvas.WritePixelMutex(job.col, job.row, color)
	// MUST BE COMMENTED OUT WHEN RUNNING LINE-MODE
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

func (rc *Context) colorAt(r mat.Ray, remainingReflections int, remainingRefractions int) mat.Tuple4 {
	rc.total++
	rc.depth++

	rc.cStack[rc.total].WorldXS = mat.IntersectWithWorldPtr(rc.world, r, rc.cStack[rc.total].WorldXS, &rc.cStack[rc.total].ShadowInRay)
	if len(rc.cStack[rc.total].WorldXS) > 0 {
		var found = false
		hit, found := mat.Hit(rc.cStack[rc.total].WorldXS)
		if found {
			mat.PrepareComputationForIntersectionPtr(hit, r, &rc.cStack[rc.total].Comps, rc.cStack[rc.total].WorldXS...)
			clr := rc.shadeHit(rc.cStack[rc.total].Comps, remainingReflections, remainingRefractions)

			if clr.Elems[0] < 0 || clr.Elems[1] < 0 || clr.Elems[2] < 0 {
				panic("negative color!!")
			}
			return clr
		}
		return black
	} else {
		return black
	}
}

func (rc *Context) reflectedColor(comps mat.Computation, remainingReflections, remainingRefractions int) mat.Tuple4 {
	if remainingReflections <= 0 || comps.Object.GetMaterial().Reflectivity == 0.0 {
		return black
	}
	reflectRay := mat.NewRay(comps.OverPoint, comps.ReflectVec)
	remainingReflections--
	reflectedColor := rc.colorAt(reflectRay, remainingReflections, remainingRefractions)
	//return mat.MultiplyByScalar(reflectedColor, comps.Object.GetMaterial().Reflectivity)
	return reflectedColor.Multiply(comps.Object.GetMaterial().Reflectivity)
}

func (rc *Context) refractedColor(comps mat.Computation, remainingReflections, remainingRefractions int) mat.Tuple4 {
	if remainingRefractions <= 0 || comps.Object.GetMaterial().Transparency == 0.0 {
		return black
	}

	if comps.N2 == 0.0 {
		fmt.Println("Warn: About to divide by zero. Im not Chuck Norris.")
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
	remainingRefractions--
	nextColor := rc.colorAt(refractRay, remainingRefractions, remainingReflections)
	transparency := comps.Object.GetMaterial().Transparency
	//color := mat.MultiplyByScalar(nextColor, transparency)

	return nextColor.Multiply(transparency)
}

func (rc *Context) shadeHit(comps mat.Computation, remainingReflections, remainingRefractions int) mat.Tuple4 {
	var surfaceColor = mat.NewColor(0, 0, 0)
	for _, light := range rc.world.Light {
		inShadow := rc.pointInShadow(light, comps.OverPoint)
		color := mat.Lighting(comps.Object.GetMaterial(), comps.Object, light, comps.OverPoint, comps.EyeVec, comps.NormalVec, inShadow, rc.cStack[rc.total].LightData)
		surfaceColor = mat.Add(surfaceColor, color)
	}
	reflectedColor := rc.reflectedColor(comps, remainingReflections, remainingRefractions)
	refractedColor := rc.refractedColor(comps, remainingReflections, remainingRefractions)

	material := comps.Object.GetMaterial()
	if material.Reflectivity > 0.0 && material.Transparency > 0.0 {
		reflectance := mat.Schlick(comps)
		//return surfaceColor.Add(reflectedColor.Multiply(reflectance)).Add(refractedColor.Multiply(1 - reflectance))
		return mat.Add(mat.Add(surfaceColor, reflectedColor.Multiply(reflectance)), refractedColor.Multiply(1-reflectance))
		//mat.Add3(surfaceColor, reflectedColor.Multiply(reflectance), refractedColor.Multiply(1 - reflectance), &comps.ShadeColor)
		//return comps.ShadeColor
	} else {
		//return surfaceColor.Add(reflectedColor.Add(refractedColor))
		return mat.Add(surfaceColor, mat.Add(reflectedColor, refractedColor))
		//mat.Add3(surfaceColor, reflectedColor, refractedColor, &comps.ShadeColor)
		//return comps.ShadeColor
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
			if x.T > 0.0 && x.T < distance {
				return true
			}
		}
	}
	return false
}

// lighting computes the color of a given pixel given phong shading
func (rc *Context) lighting(material mat.Material, object mat.Shape, light mat.Light, position, eyeVec, normalVec mat.Tuple4, inShadow bool) mat.Tuple4 {
	var color mat.Tuple4
	if material.HasPattern() {
		color = mat.PatternAtShape(material.Pattern, object, position)
	} else {
		color = material.Color
	}
	if inShadow {
		return mat.MultiplyByScalar(color, material.Ambient)
	}
	effectiveColor := mat.Hadamard(color, light.Intensity)

	// get vector from point on sphere to light source by subtracting, normalized into unit space.
	lightVec := mat.Normalize(mat.Sub(light.Position, position))

	// Add the ambient portion
	ambient := mat.MultiplyByScalar(effectiveColor, material.Ambient)

	// get dot product (angle) between the light and normal  vectors. If negative, it means the light source is
	// on the backside
	lightDotNormal := mat.Dot(lightVec, normalVec)
	specular := mat.NewColor(0, 0, 0)
	diffuse := mat.NewColor(0, 0, 0)
	if lightDotNormal < 0 {
		diffuse = black
		specular = black
	} else {
		// Diffuse contribution Precedense here??
		diffuse = mat.MultiplyByScalar(effectiveColor, material.Diffuse*lightDotNormal)

		// reflect_dot_eye represents the cosine of the angle between the
		// reflection vector and the eye vector. A negative number means the
		// light reflects away from the eye.
		// Note that we negate the light vector since we want to angle of the bounce
		// of the light rather than the incoming light vector.
		reflectVec := mat.Reflect(mat.Negate(lightVec), normalVec)
		reflectDotEye := mat.Dot(reflectVec, eyeVec)

		if reflectDotEye <= 0.0 {
			specular = black
		} else {
			// compute the specular contribution
			factor := math.Pow(reflectDotEye, material.Shininess)

			// again, check precedense here
			specular = mat.MultiplyByScalar(light.Intensity, material.Specular*factor)
		}
	}
	// Add the three contributions together to get the final shading
	// Uses standard Tuple addition
	return ambient.Add(diffuse.Add(specular))
	//return Add(Add(ambient, diffuse), specular)
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

	LightData mat.LightData
}

func NewShadeData() ShadeData {
	worldXS := make([]mat.Intersection, 16)
	shadowXS := make([]mat.Intersection, 16)

	worldXS = worldXS[:0]
	shadowXS = shadowXS[:0]

	return ShadeData{
		WorldXS:     worldXS,
		ShadowXS:    shadowXS,
		InRay:       mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0)),
		ShadowInRay: mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0)),
		Comps:       mat.NewComputation(),
		LightData:   mat.NewLightData(),
	}
}
