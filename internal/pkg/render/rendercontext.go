package render

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/calcstats"
	"github.com/eriklupander/rt/internal/pkg/config"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/inhies/go-bytesize"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// constants
var originPoint = mat.NewPoint(0, 0, 0)
var black = mat.NewColor(0, 0, 0)

// New creates a new render context to be used exclusively by a single Render worker
func New(world mat.World) Context {

	// allocate a "stack" of 256 ShadeData instances, e.g. meaning that the render of a single pixel may recurse and
	// spawn up to 256 additonal rays/intersection tests without having to allocate new memory.
	cStack := make([]ShadeData, 1024)
	for i := 0; i < 1024; i++ {
		cStack[i] = NewShadeData()
	}

	samples := make([]mat.Tuple4, 16)
	for i := 0; i < len(samples); i++ {
		samples[i] = mat.NewColor(0, 0, 0)
	}

	return Context{
		world: world,
		total: 0,

		// allocate memory
		pointInView:     mat.NewPoint(0, 0, -1.0),
		pixel:           mat.NewColor(0, 0, 0),
		origin:          mat.NewPoint(0, 0, 0),
		direction:       mat.NewVector(0, 0, 0),
		subVec:          mat.NewVector(0, 0, 0),
		shadowDirection: mat.NewVector(0, 0, 0),

		// allocate ray
		firstRay: mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0)),

		// stack for shading
		cStack: cStack,

		samples: samples,
	}
}

// NewContext uses the passed parameters after creating a render context.
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
	pointInView     mat.Tuple4
	pixel           mat.Tuple4
	origin          mat.Tuple4
	direction       mat.Tuple4
	subVec          mat.Tuple4
	shadowDirection mat.Tuple4

	// ray cache
	firstRay mat.Ray

	// each renderContext needs to pre-allocate shade-data for sufficient number of recursions
	cStack []ShadeData

	// experiment, alloc memory for each sample of a given pixel
	samples []mat.Tuple4
}

// Threaded sets up workers and producer for rendering the passed camera + slice of worlds.
// Note that the length of worlds must be equal or greater than the number of config.Cfg.Threads
func Threaded(c mat.Camera, worlds []mat.World) *mat.Canvas {
	if len(worlds) < config.Cfg.Threads {
		panic("Number of world instances must be equal or greater than the configured renderThreads")
	}
	st := time.Now()
	canvas := mat.NewCanvas(c.Width, c.Height)
	jobs := make(chan *job)

	wg := sync.WaitGroup{}
	wg.Add(canvas.H)

	// Create the render contexts, one per worker
	renderContexts := make([]Context, config.Cfg.Threads)
	for i := 0; i < config.Cfg.Threads; i++ {
		renderContexts[i] = NewContext(i, worlds[i], c, canvas, jobs, &wg)
	}

	// start workers
	for i := 0; i < config.Cfg.Threads; i++ {
		go renderContexts[i].workerFuncPerLine()
	}

	// start passing work to the workers, one line at a time
	for row := 0; row < c.Height; row++ {
		jobs <- &job{row: row, col: 0}
		fmt.Printf("%d/%d\n", row, c.Height)
	}

	wg.Wait()
	fmt.Println("All done")
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)
	fmt.Printf("Memory: %v ", bytesize.New(float64(stats.Alloc)).String())
	fmt.Printf("Mallocs: %v ", stats.Mallocs)
	fmt.Printf("Total alloc: %v\n", bytesize.New(float64(stats.TotalAlloc)).String())
	fmt.Printf("%v\n", time.Now().Sub(st))
	fmt.Printf("XS skipped in group: %v\n", calcstats.Cnt)
	fmt.Printf("Transpose calls: %v\n", calcstats.Tpose)
	fmt.Printf("Dot calls: %v\n", calcstats.Dots)
	fmt.Printf("Cross calls: %v\n", calcstats.Crosses)
	fmt.Printf("Normalize calls: %v\n", calcstats.Ns)

	fmt.Println()
	fmt.Printf("|%v|%v|%v|%v|%v|%v|\n",
		bytesize.New(float64(stats.Alloc)).String(),
		stats.Mallocs,
		bytesize.New(float64(stats.TotalAlloc)).String(),
		time.Now().Sub(st),
		calcstats.Cnt,
		calcstats.Tpose)

	return canvas
}

func (rc *Context) workerFuncPerPixel() {
	for job := range rc.jobs {
		rc.renderPixelPinhole(job)
	}
}
func (rc *Context) workerFuncPerLine() {
	for job := range rc.jobs {
		for i := 0; i < rc.camera.Width; i++ {
			job.col = i
			if rc.camera.Aperture == 0.0 {
				rc.renderPixelPinhole(job)
			} else {
				rc.renderPixelWithAperture(job)
			}
		}
		rc.wg.Done()
	}
}

func (rc *Context) renderSinglePixel(col, row int) mat.Tuple4 {
	for i := 0; i < 1024; i++ {
		rc.cStack[i].WorldXS = rc.cStack[i].WorldXS[:0]
		rc.cStack[i].ShadowXS = rc.cStack[i].ShadowXS[:0]
	}
	rc.total = 0
	rc.depth = 0
	rc.rayForPixel(col, row, &rc.firstRay)
	color := rc.colorAt(rc.firstRay, 5, 5)
	return color
}

func (rc *Context) renderPixelPinhole(job *job) {
	for i := 0; i < 256; i++ {
		rc.cStack[i].WorldXS = rc.cStack[i].WorldXS[:0]
		rc.cStack[i].ShadowXS = rc.cStack[i].ShadowXS[:0]
	}
	rc.samples = rc.samples[:0]
	for i := 0; i < config.Cfg.Samples; i++ {
		rc.total = 0
		rc.depth = 0

		// TODO optimize so we take four samples in each corner and then see how much they differ. If below threshold,
		// just take a center one as well and return the interpolated result. Otherwise, pick N random samples.
		rc.rayForPixel(job.col, job.row, &rc.firstRay)
		rc.samples = append(rc.samples, rc.colorAt(rc.firstRay, 5, 5))
	}

	// calc avg color
	rc.canvas.WritePixelMutex(job.col, job.row, rc.sumColors())
}

func (rc *Context) renderPixelWithAperture(job *job) {
	tempPos := mat.Tuple4{}
	newVec := mat.Tuple4{}
	var pos = mat.Tuple4{}
	// experiment: run rayForPixel + colorAt N times, with random offset within the pixel
	// Then compute the average color of all
	rc.samples = rc.samples[:0]
	for i := 0; i < config.Cfg.Samples; i++ {
		for i := 0; i < 1024; i++ {
			rc.cStack[i].WorldXS = rc.cStack[i].WorldXS[:0]
			rc.cStack[i].ShadowXS = rc.cStack[i].ShadowXS[:0]
		}
		rc.total = 0
		rc.depth = 0

		rc.rayForPixel(job.col, job.row, &rc.firstRay)
		// DoF experiment starts here!
		// get focal point on this ray, which is distance on ray

		mat.PositionPtr(rc.firstRay, rc.camera.FocalLength, &pos)

		// now, move camera origin and cast a ray from new origin through pos

		for j := 0; j < 128; j++ { // -1, 0.5, -3
			tempPos[0] = rc.firstRay.Origin[0] + (-rc.camera.Aperture + rand.Float64()*rc.camera.Aperture*2)
			tempPos[1] = rc.firstRay.Origin[1] + (-rc.camera.Aperture + rand.Float64()*rc.camera.Aperture*2)
			tempPos[2] = rc.firstRay.Origin[2]
			tempPos[3] = 1
			mat.SubPtr(pos, tempPos, &newVec)
			tempRay := mat.NewRay(tempPos, newVec)
			rc.samples = append(rc.samples, rc.colorAt(tempRay, 5, 5))
		}
	}

	// calc avg color
	rc.canvas.WritePixelMutex(job.col, job.row, rc.sumColors())
}

func (rc *Context) rayForPixel(x, y int, out *mat.Ray) {

	xOffset := rc.camera.PixelSize * (float64(x) + 0.5) // 0.5
	yOffset := rc.camera.PixelSize * (float64(y) + 0.5) // 0.5

	// this feels a little hacky but actually works.
	worldX := rc.camera.HalfWidth - xOffset
	worldY := rc.camera.HalfHeight - yOffset

	// mat.NewPoint(worldX, worldY, -1.0)
	rc.pointInView[0] = worldX
	rc.pointInView[1] = worldY

	mat.MultiplyByTuplePtr(&rc.camera.Inverse, &rc.pointInView, &rc.pixel)
	mat.MultiplyByTuplePtr(&rc.camera.Inverse, &originPoint, &rc.origin)
	mat.SubPtr(rc.pixel, rc.origin, &rc.subVec)
	mat.NormalizePtr(&rc.subVec, &rc.direction)

	out.Direction = rc.direction
	out.Origin = rc.origin
}
func (rc *Context) rayForPixelRand(x, y int, out *mat.Ray) {

	xOffset := rc.camera.PixelSize * (float64(x) + rand.Float64()) // 0.5
	yOffset := rc.camera.PixelSize * (float64(y) + rand.Float64()) // 0.5

	// this feels a little hacky but actually works.
	worldX := rc.camera.HalfWidth - xOffset
	worldY := rc.camera.HalfHeight - yOffset

	// mat.NewPoint(worldX, worldY, -1.0)
	rc.pointInView[0] = worldX
	rc.pointInView[1] = worldY

	mat.MultiplyByTuplePtr(&rc.camera.Inverse, &rc.pointInView, &rc.pixel)
	mat.MultiplyByTuplePtr(&rc.camera.Inverse, &originPoint, &rc.origin)
	mat.SubPtr(rc.pixel, rc.origin, &rc.subVec)
	mat.NormalizePtr(&rc.subVec, &rc.direction)

	out.Direction = rc.direction
	out.Origin = rc.origin
}

func (rc *Context) rayForPixelDoF(x, y int, out *mat.Ray) {

	xOffset := rc.camera.PixelSize * (float64(x) + rand.Float64()) // 0.5
	yOffset := rc.camera.PixelSize * (float64(y) + rand.Float64()) // 0.5

	// this feels a little hacky but actually works.
	worldX := rc.camera.HalfWidth - xOffset
	worldY := rc.camera.HalfHeight - yOffset

	// mat.NewPoint(worldX, worldY, -1.0)
	rc.pointInView[0] = worldX
	rc.pointInView[1] = worldY

	mat.MultiplyByTuplePtr(&rc.camera.Inverse, &rc.pointInView, &rc.pixel)
	mat.MultiplyByTuplePtr(&rc.camera.Inverse, &originPoint, &rc.origin)
	mat.SubPtr(rc.pixel, rc.origin, &rc.subVec)
	mat.NormalizePtr(&rc.subVec, &rc.direction)

	out.Direction = rc.direction
	out.Origin = rc.origin
}

func (rc *Context) intensityAt(light mat.AreaLight, point mat.Tuple4) float64 {
	total := 0.0

	for v := 0; v < light.VSteps; v++ {
		for u := 0; u < light.USteps; u++ {
			lightPos := mat.PointOnLight(light, float64(u), float64(v))
			if !rc.isShadowed(lightPos, point) {
				total = total + 1.0
			}
		}
	}
	return total / light.Samples
}
func (rc *Context) isShadowed(lightPosition mat.Tuple4, p mat.Tuple4) bool {
	vecToLight := mat.Sub(lightPosition, p)
	distance := mat.Magnitude(vecToLight)

	mat.NormalizePtr(&vecToLight, &rc.shadowDirection)
	ray := mat.NewRay(p, rc.shadowDirection)

	return mat.ShadowIntersect(rc.world, ray, distance, &rc.cStack[rc.total].InRay) //mat.IntersectWithWorldPtrForShadow(rc.world, ray, rc.cStack[rc.total].ShadowXS, &rc.cStack[rc.total].InRay)
	// use stack...
	//rc.cStack[rc.total].ShadowXS = mat.IntersectWithWorldPtrForShadow(rc.world, ray, rc.cStack[rc.total].ShadowXS, distance, &rc.cStack[rc.total].InRay)
	//if len(rc.cStack[rc.total].ShadowXS) > 0 {
	//	for _, x := range rc.cStack[rc.total].ShadowXS {
	//		if x.T > 0.0 && x.T < distance {
	//			return true
	//		}
	//	}
	//}
	//return false
}
func (rc *Context) sumColors() mat.Tuple4 {
	var r, g, b float64
	for i := range rc.samples {
		r += rc.samples[i][0] * rc.samples[i][0]
		g += rc.samples[i][1] * rc.samples[i][1]
		b += rc.samples[i][2] * rc.samples[i][2]
	}
	n := float64(len(rc.samples))
	return mat.NewColor(math.Sqrt(r/n), math.Sqrt(g/n), math.Sqrt(b/n))
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
			return clr
		}
		return black
	}
	return black
}

func (rc *Context) reflectedColor(comps mat.Computation, remainingReflections, remainingRefractions int) mat.Tuple4 {
	if remainingReflections <= 0 || comps.Object.GetMaterial().Reflectivity == 0.0 {
		return black
	}
	reflectRay := mat.NewRay(comps.OverPoint, comps.ReflectVec)
	remainingReflections--
	reflectedColor := rc.colorAt(reflectRay, remainingReflections, remainingRefractions)
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

	return nextColor.Multiply(transparency)
}

func (rc *Context) shadeHit(comps mat.Computation, remainingReflections, remainingRefractions int) mat.Tuple4 {
	var surfaceColor = &mat.Tuple4{0, 0, 0, 1.0} //mat.NewColor(0,0,0)

	// Light for point lights
	for _, light := range rc.world.Light {
		isShadowed := rc.isShadowed(light.Position, comps.OverPoint)
		intensity := 1.0
		if isShadowed {
			intensity = 0.0
		}
		color := mat.LightingPointLight(comps.Object.GetMaterial(), comps.Object, light, comps.OverPoint, comps.EyeVec, comps.NormalVec, intensity == 0.0, rc.cStack[rc.total].LightData)
		surfaceColor.AddNoRet(color)
	}

	// Light for area lights
	for i := range rc.world.AreaLight {
		intensity := rc.intensityAt(rc.world.AreaLight[i], comps.OverPoint)
		color := mat.Lighting(comps.Object.GetMaterial(), comps.Object, rc.world.AreaLight[i], comps.OverPoint, comps.EyeVec, comps.NormalVec, intensity, rc.cStack[rc.total].LightData)
		surfaceColor.AddNoRet(color)
	}
	reflectedColor := rc.reflectedColor(comps, remainingReflections, remainingRefractions)
	refractedColor := rc.refractedColor(comps, remainingReflections, remainingRefractions)

	material := comps.Object.GetMaterial()
	if material.Reflectivity > 0.0 && material.Transparency > 0.0 {
		reflectance := mat.Schlick(comps)
		surfaceColor.AddNoRet(reflectedColor.Multiply(reflectance))
		surfaceColor.AddNoRet(refractedColor.Multiply(1 - reflectance))
		return *surfaceColor
	} else {
		surfaceColor.AddNoRetPtr(&reflectedColor)
		surfaceColor.AddNoRetPtr(&refractedColor)
		return *surfaceColor
	}
}

func (rc *Context) pointInShadow(light mat.Light, p mat.Tuple4) bool {

	vecToLight := mat.Sub(light.Position, p)
	distance := mat.Magnitude(vecToLight)

	mat.NormalizePtr(&vecToLight, &rc.shadowDirection)
	ray := mat.NewRay(p, rc.shadowDirection)

	// use stack...
	rc.cStack[rc.total].ShadowXS = mat.IntersectWithWorldPtrForShadow(rc.world, ray, rc.cStack[rc.total].ShadowXS, &rc.cStack[rc.total].InRay)
	if len(rc.cStack[rc.total].ShadowXS) > 0 {
		for _, x := range rc.cStack[rc.total].ShadowXS {
			if x.T > 0.0 && x.T < distance {
				return true
			}
		}
	}
	return false
}

//
//// lighting computes the color of a given pixel given phong shading
//func (rc *Context) lighting(material mat.Material, object mat.Shape, light mat.Light, position, eyeVec, normalVec mat.Tuple4, inShadow bool) mat.Tuple4 {
//	var color mat.Tuple4
//	if material.HasPattern() {
//		color = mat.PatternAtShape(material.Pattern, object, position)
//	} else {
//		color = material.Color
//	}
//	if inShadow {
//		return mat.MultiplyByScalar(color, material.Ambient)
//	}
//	effectiveColor := mat.Hadamard(color, light.Intensity)
//
//	// get vector from point on sphere to light source by subtracting, normalized into unit space.
//	lightVec := mat.Normalize(mat.Sub(light.Position, position))
//
//	// Add the ambient portion
//	ambient := mat.MultiplyByScalar(effectiveColor, material.Ambient)
//
//	// get dot product (angle) between the light and normal  vectors. If negative, it means the light source is
//	// on the backside
//	lightDotNormal := mat.Dot(lightVec, normalVec)
//	specular := mat.NewColor(0, 0, 0)
//	diffuse := mat.NewColor(0, 0, 0)
//	if lightDotNormal < 0 {
//		diffuse = black
//		specular = black
//	} else {
//		// Diffuse contribution Precedense here??
//		diffuse = mat.MultiplyByScalar(effectiveColor, material.Diffuse*lightDotNormal)
//
//		// reflect_dot_eye represents the cosine of the angle between the
//		// reflection vector and the eye vector. A negative number means the
//		// light reflects away from the eye.
//		// Note that we negate the light vector since we want to angle of the bounce
//		// of the light rather than the incoming light vector.
//		reflectVec := mat.Reflect(mat.Negate(lightVec), normalVec)
//		reflectDotEye := mat.Dot(reflectVec, eyeVec)
//
//		if reflectDotEye <= 0.0 {
//			specular = black
//		} else {
//			// compute the specular contribution
//			factor := math.Pow(reflectDotEye, material.Shininess)
//
//			// again, check precedense here
//			specular = mat.MultiplyByScalar(light.Intensity, material.Specular*factor)
//		}
//	}
//	// Add the three contributions together to get the final shading
//	// Uses standard Tuple addition
//	return ambient.Add(diffuse.Add(specular))
//}

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

// NewShadeData pre-allocates intersection lists etc
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
