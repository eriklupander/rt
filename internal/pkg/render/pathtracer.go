package render

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"math"
	"sort"
)

// trying to do a full path tracer in a single method...
func (rc *Context) renderPixelPathTracer(job *job) {

	// slice to store the final color of each sample in
	var samples = make([]mat.Tuple4, 0)

	// Sample N times per pixel
	for sample := 0; sample < 64; sample++ {
		// Part 1: Determine if the ray intersects something in our scene

		// cast the ray through a random spot in the pixel, store the result in "firstRay"
		rc.rayForPixelPathTracer(job.col, job.row, &rc.firstRay)
		color, hit := rc.castRay(rc.firstRay, 0, job.col, job.row)
		if hit {
			samples = append(samples, color)
		}
	}

	rc.canvas.WritePixelMutex(job.col, job.row, mat.MultiplyByScalar(sumColors(samples, 0), 2.2))
}

func (rc *Context) castRay(ray mat.Ray, depth, x, y int) (mat.Tuple4, bool) {
	if depth > 3 {
		return black, false
	}
	depth++
	// slice to store intersections for a single ray
	var xs []mat.Intersection

	// storage for a single ray transformed for doing an intersection test
	var transformedRay = mat.NewRay(mat.NewTuple(), mat.NewTuple())
	var normalizedVectorToLight = mat.NewTuple()

	// find all intersections
	for i := range rc.world.Objects {

		// transforming the ray into object space makes the intersection math much easier
		mat.TransformRayPtr(ray, rc.world.Objects[i].GetInverse(), &transformedRay)

		// Call the intersect function provided by the shape implementation (e.g. Sphere, Plane osv)
		// and append any results to the global intersection list
		xs = append(xs, rc.world.Objects[i].IntersectLocal(transformedRay)...)
	}

	// If there were no intersection
	if len(xs) == 0 {
		return black, false
	}

	// sort intersections
	intersections := mat.Intersections(xs)
	if len(xs) > 1 {
		sort.Sort(intersections)
	}

	// loop over all intersections and find the first positive one
	for i := 0; i < len(intersections); i++ {

		// Check is positive (in front of camera)
		if intersections[i].T > 0.0 {
			material := intersections[i].S.GetMaterial()

			// time to calculate normals and stuff for his particular intersection
			computations := mat.PrepareComputationForIntersection(intersections[i], ray, intersections...)

			// START SHADOW and LIGHT code
			// check if in shadow or not by getting a normalized vector from point to light

			// get vector from point on sphere to light source by subtracting, normalized into unit space.
			vecToLight := mat.Sub(rc.world.Light[0].Position, computations.OverPoint)
			distance := mat.Magnitude(vecToLight)
			mat.NormalizePtr(&vecToLight, &normalizedVectorToLight)

			// experiment. Use a random cone to simulate area light
			normalizedVectorToLight = rc.RandomConeInHemisphere(normalizedVectorToLight, 0.1)

			lightRay := mat.NewRay(computations.OverPoint, normalizedVectorToLight)

			// next for shadows, check if any object in the scene intersects the ray between point and light
			// source
			isShadowed := rc.checkShadowPathTracer(lightRay, distance)

			// END SHADOW CHECK

			// START COLORING OF PIXEL!!!!!!
			//if job.col == 158 && job.row == 288 {
			//	fmt.Printf("%v\n", "HEJ")
			//}
			var color mat.Tuple4
			if !isShadowed {

				// storage for light stuff
				lightData := mat.NewLightData()
				lightData.LightVec = normalizedVectorToLight
				light := rc.world.Light[0]

				// hadamard is a fancy name for multiplying each col of two vectors
				mat.HadamardPtr(&material.Color, light.Intensity, &lightData.EffectiveColor)

				// Add the ambient portion
				mat.MultiplyByScalarPtr(lightData.EffectiveColor, material.Ambient, &lightData.Ambient)

				// get dot product (angle) between the light and normal  vectors. If negative, it means the light source is
				// on the backside
				lightDotNormal := mat.Dot(lightData.LightVec, computations.NormalVec)
				specular := lightData.Specular
				diffuse := lightData.Diffuse

				if lightDotNormal < 0 {
					diffuse = black
					specular = black
				} else {
					// Diffuse contribution Precedense here??

					mat.MultiplyByScalarPtr(lightData.EffectiveColor, material.Diffuse*lightDotNormal, &diffuse)

					// reflect_dot_eye represents the cosine of the angle between the
					// reflection vector and the eye vector. A negative number means the
					// light reflects away from the eye.
					// Note that we negate the light vector since we want to angle of the bounce
					// of the light rather than the incoming light vector.

					mat.ReflectPtr(mat.Negate(lightData.LightVec), computations.NormalVec, &lightData.ReflectVec)
					reflectDotEye := mat.Dot(lightData.ReflectVec, computations.EyeVec)

					if reflectDotEye <= 0.0 {
						specular = black
					} else {
						// compute the specular contribution
						factor := math.Pow(reflectDotEye, material.Shininess)

						// again, check precedense here
						mat.MultiplyByScalarPtr(light.Intensity, material.Specular*factor, &specular)
					}
				}
				// Add the three contributions together to get the final shading
				// Uses standard Tuple addition
				color = lightData.Ambient.Add(diffuse.Add(specular))
			} else {
				// in shadow
				color = mat.MultiplyByScalar(material.Color, material.Ambient)
			}

			// we could say that once we're here, we are done with the DIRECT light.
			// It's then time to start the Monte Carlo integrating.
			// collect N indirect samples
			normalVec := computations.NormalVec

			var indirectSamples []mat.Tuple4
			for indirect := 0; indirect < 1; indirect++ {
				// get a random unit vector in the normal's hemisphere
				rndVec := rc.RandomConeInHemisphere(normalVec, 1)
				rndRay := mat.NewRay(computations.OverPoint, rndVec)
				// 158x288

				color2, hit := rc.castRay(rndRay, depth, x, y)
				if hit {
					indirectSamples = append(indirectSamples, color2)
				}
			}
			indirectSamples = append(indirectSamples, color)
			return sumColors(indirectSamples, depth), true
		}
	}

	return black, false
}

func (rc *Context) checkShadowPathTracer(lightRay mat.Ray, distance float64) bool {
	var transformedShadowRay = mat.NewRay(mat.NewTuple(), mat.NewTuple())

	for idx := range rc.world.Objects {
		if !rc.world.Objects[idx].CastsShadow() {
			continue
		}
		intersections := mat.IntersectRayWithShapePtr(rc.world.Objects[idx], lightRay, &transformedShadowRay)

		for innerIdx := range intersections {
			if intersections[innerIdx].T > 0.0 && intersections[innerIdx].T < distance {
				return true
			}
		}
	}
	return false
}

func sumColors(samples []mat.Tuple4, depth int) mat.Tuple4 {
	var r, g, b float64
	for i := range samples {
		r += samples[i][0] * samples[i][0]
		g += samples[i][1] * samples[i][1]
		b += samples[i][2] * samples[i][2]
	}
	n := float64(len(samples) * (depth+1)) // little test to decrease the contribution the deeper we've gone
	return mat.NewColor(math.Sqrt(r/n), math.Sqrt(g/n), math.Sqrt(b/n))
}

func (rc *Context) rayForPixelPathTracer(x, y int, out *mat.Ray) {

	xOffset := rc.camera.PixelSize * (float64(x) + rc.rnd.Float64()) // 0.5
	yOffset := rc.camera.PixelSize * (float64(y) + rc.rnd.Float64()) // 0.5

	// this feels a little hacky but actually works.
	worldX := rc.camera.HalfWidth - xOffset
	worldY := rc.camera.HalfHeight - yOffset

	rc.pointInView[0] = worldX
	rc.pointInView[1] = worldY

	mat.MultiplyByTuplePtr(&rc.camera.Inverse, &rc.pointInView, &rc.pixel)
	mat.MultiplyByTuplePtr(&rc.camera.Inverse, &originPoint, &rc.origin)
	mat.SubPtr(rc.pixel, rc.origin, &rc.subVec)
	mat.NormalizePtr(&rc.subVec, &rc.direction)

	out.Direction = rc.direction
	out.Origin = rc.origin
}


// From Hunter Loftis path tracer:
// https://github.com/hunterloftis/pbr/blob/1ce8b1c067eea7cf7298745d6976ba72ff12dd50/pkg/geom/dir.go
// In this program, think of "a" as the surface normal
//
// Cone returns a random vector within a cone about Direction a.
// size is 0-1, where 0 is the original vector and 1 is anything within the original hemisphere.
// https://github.com/fogleman/pt/blob/69e74a07b0af72f1601c64120a866d9a5f432e2f/pt/util.go#L24
func (rc *Context) RandomConeInHemisphere(startVec mat.Tuple4, size float64) mat.Tuple4 {
	u := rc.rnd.Float64()
	v := rc.rnd.Float64()
	theta := size * 0.5 * math.Pi * (1 - (2 * math.Acos(u) / math.Pi))
	m1 := math.Sin(theta)
	m2 := math.Cos(theta)
	a2 := v * 2 * math.Pi

	// q should be possible to store in Context?
	q := mat.Tuple4{}
	rc.RandDirection(&q)

	// should be possible to move s and t into Context?
	s := mat.Tuple4{}
	t := mat.Tuple4{}
	mat.Cross2(&startVec, &q, &s)
	mat.Cross2(&startVec, &s, &t)

	d := mat.Tuple4{}
	mat.MultiplyByScalarPtr(s, m1 * math.Cos(a2), &d)
	d = mat.Add(d, mat.MultiplyByScalar(t, m1 * math.Sin(a2)))
	d = mat.Add(d, mat.MultiplyByScalar(startVec, m2))
	return mat.Normalize(d)
}

// RandDirection returns a random unit vector (a point on the edge of a unit sphere).
func (rc *Context) RandDirection(out *mat.Tuple4) {
	AngleDirection(rc.rnd.Float64()*math.Pi*2, math.Asin(rc.rnd.Float64()*2-1), out)
}
func AngleDirection(theta, phi float64, out *mat.Tuple4) {
	out[0] = math.Cos(theta) * math.Cos(phi)
	out[1] = math.Sin(phi)
	out[2] = math.Sin(theta) * math.Cos(phi)
}