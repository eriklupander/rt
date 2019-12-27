package render

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestReflectedColorForNonreflectiveMaterial(t *testing.T) {

	w := mat.NewDefaultWorld()
	rc := Context{world: w}
	r := mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 1))
	shape := w.Objects[1]
	material := shape.GetMaterial()
	material.Ambient = 1.0
	shape.SetMaterial(material)
	xs := mat.NewIntersection(1, shape)

	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs, r, &comps)
	color := rc.reflectedColor(comps, 1, 1)
	assert.Equal(t, black, color)
}

func TestReflectedColorForReflectiveMaterial(t *testing.T) {
	w := mat.NewDefaultWorld()
	plane := mat.NewPlane()
	plane.SetTransform(mat.Translate(0, -1, -0))
	material := plane.GetMaterial()
	material.Reflectivity = 0.5
	plane.SetMaterial(material)
	w.Objects = append(w.Objects, plane)

	rc := New(w)

	r := mat.NewRay(mat.NewPoint(0, 0, -3), mat.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := mat.NewIntersection(math.Sqrt(2), plane)
	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs, r, &comps)
	color := rc.reflectedColor(comps, 1, 1)
	assert.InEpsilon(t, 0.19032, color.Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.2379, color.Get(1), mat.Epsilon)
	assert.InEpsilon(t, 0.14274, color.Get(2), mat.Epsilon)
}

func TestShadeHitWithReflectiveMaterial(t *testing.T) {
	w := mat.NewDefaultWorld()
	plane := mat.NewPlane()
	plane.SetTransform(mat.Translate(0, -1, -0))
	material := plane.GetMaterial()
	material.Reflectivity = 0.5
	plane.SetMaterial(material)
	w.Objects = append(w.Objects, plane)

	rc := New(w)

	r := mat.NewRay(mat.NewPoint(0, 0, -3), mat.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := mat.NewIntersection(math.Sqrt(2), plane)
	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs, r, &comps)
	color := rc.shadeHit(comps, 1, 1)

	assert.InEpsilon(t, 0.87677, color.Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.92436, color.Get(1), mat.Epsilon)
	assert.InEpsilon(t, 0.82918, color.Get(2), mat.Epsilon)
}

func TestColorAtWithMutuallyReflectiveSurfaces(t *testing.T) {
	w := mat.NewWorld()
	w.Light = []mat.Light{mat.NewLight(mat.NewPoint(0, 0, 0), mat.NewColor(1, 1, 1))}
	lowerPlane := mat.NewPlane()
	lowerPlane.SetMaterial(mat.NewDefaultReflectiveMaterial(1.0))
	lowerPlane.SetTransform(mat.Translate(0, -1, 0))
	w.Objects = append(w.Objects, lowerPlane)

	upperPlane := mat.NewPlane()
	upperPlane.SetMaterial(mat.NewDefaultReflectiveMaterial(1.0))
	upperPlane.SetTransform(mat.Translate(0, 1, 0))
	w.Objects = append(w.Objects, upperPlane)

	rc := New(w)

	r := mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 1, 0))
	_ = rc.colorAt(r, 1, 5)
}

func TestTheReflectedColorAtMaxRecursiveDepth(t *testing.T) {
	w := mat.NewWorld()
	pl := mat.NewPlane()
	pl.SetMaterial(mat.NewDefaultReflectiveMaterial(0.5))
	pl.SetTransform(mat.Translate(0, -1, 0))
	w.Objects = append(w.Objects, pl)

	rc := New(w)

	r := mat.NewRay(mat.NewPoint(0, 0, -3), mat.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := mat.NewIntersection(math.Sqrt(2), pl)
	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs, r, &comps)
	color := rc.reflectedColor(comps, 0, 0)
	assert.Equal(t, black, color)
}

func TestRayForPixelThroughCenterOfCanvas(t *testing.T) {
	cam := mat.NewCamera(201, 101, math.Pi/2.0)
	copy(cam.Inverse.Elems, mat.IdentityMatrix.Elems)
	rc := NewContext(0, mat.NewWorld(), cam, nil, nil, nil)
	r := mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0))
	rc.rayForPixel(100, 50, &r)
	assert.Equal(t, mat.NewPoint(0, 0, 0), r.Origin)
	assert.Equal(t, mat.NewVector(0, 0, -1), r.Direction)
}

func TestRayForPixelThroughCornerOfCanvas(t *testing.T) {
	cam := mat.NewCamera(201, 101, math.Pi/2.0)
	rc := NewContext(0, mat.NewWorld(), cam, nil, nil, nil)
	r := mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0))
	rc.rayForPixel(0, 0, &r)
	assert.Equal(t, mat.NewPoint(0, 0, 0), r.Origin)
	assert.InEpsilon(t, 0.66519, r.Direction.Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.33259, r.Direction.Get(1), mat.Epsilon)
	assert.InEpsilon(t, -0.66851, r.Direction.Get(2), mat.Epsilon)
}

// Page 103, third testx
func TestRayForPixelWhenCamIsTransformed(t *testing.T) {
	cam := mat.NewCamera(201, 101, math.Pi/2.0)
	cam.Transform = mat.Multiply(mat.RotateY(math.Pi/4), mat.Translate(0, -2, 5))
	cam.Inverse = mat.Inverse(cam.Transform)
	rc := NewContext(0, mat.NewWorld(), cam, nil, nil, nil)
	r := mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0))
	rc.rayForPixel(100, 50, &r)
	assert.Equal(t, mat.NewPoint(0, 2, -5), r.Origin)
	assert.True(t, mat.TupleEquals(mat.NewVector(math.Sqrt(2.0)/2.0, 0.0, -math.Sqrt(2.0)/2.0), r.Direction))
}

// Page 104
func TestRender(t *testing.T) {
	worlds := make([]mat.World, 0)
	for i := 0; i < 8; i++ {
		w := mat.NewDefaultWorld()
		worlds = append(worlds, w)
	}
	c := mat.NewCamera(11, 11, math.Pi/2)
	from := mat.NewPoint(0, 0, -5)
	to := mat.NewPoint(0, 0, 0)
	upVec := mat.NewVector(0, 1, 0)
	c.Transform = mat.ViewTransform(from, to, upVec)
	c.Inverse = mat.Inverse(c.Transform)

	canvas := Threaded(c, worlds)

	assert.InEpsilon(t, 0.38066, canvas.ColorAt(5, 5).Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.47583, canvas.ColorAt(5, 5).Get(1), mat.Epsilon)
	assert.InEpsilon(t, 0.2855, canvas.ColorAt(5, 5).Get(2), mat.Epsilon)
}

// Page 95
func TestShadeIntersection(t *testing.T) {
	w := mat.NewDefaultWorld()
	rc := New(w)
	r := mat.NewRay(mat.NewPoint(0, 0, -5), mat.NewVector(0, 0, 1))
	i := mat.Intersection{T: 4.0, S: w.Objects[0]}

	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(i, r, &comps)
	color := rc.shadeHit(comps, 1, 1)
	assert.InEpsilon(t, 0.38066, color.Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.47583, color.Get(1), mat.Epsilon)
	assert.InEpsilon(t, 0.2855, color.Get(2), mat.Epsilon)
}

func TestShadeIntersectionFromInside(t *testing.T) {
	w := mat.NewDefaultWorld()
	w.Light = []mat.Light{mat.NewLight(mat.NewPoint(0, 0.25, 0), mat.NewColor(1, 1, 1))}
	rc := New(w)
	r := mat.NewRay(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 1))

	i := mat.Intersection{T: 0.5, S: w.Objects[1]}

	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(i, r, &comps)
	color := rc.shadeHit(comps, 1, 1)
	assert.InEpsilon(t, 0.90498, color.Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.90498, color.Get(1), mat.Epsilon)
	assert.InEpsilon(t, 0.90498, color.Get(2), mat.Epsilon)
}

func TestColorWhenRayMiss(t *testing.T) {
	w := mat.NewDefaultWorld()
	rc := New(w)
	r := mat.NewRay(mat.NewPoint(0, 0, -5), mat.NewVector(0, 1, 0))
	color := rc.colorAt(r, 1, 5)
	assert.Equal(t, color.Elems[0], 0.0)
	assert.Equal(t, color.Elems[1], 0.0)
	assert.Equal(t, color.Elems[2], 0.0)
}

func TestColorWhenRayHits(t *testing.T) {
	w := mat.NewDefaultWorld()
	rc := New(w)
	r := mat.NewRay(mat.NewPoint(0, 0, -5), mat.NewVector(0, 0, 1))
	color := rc.colorAt(r, 1, 5)
	assert.InEpsilon(t, 0.38066, color.Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.47583, color.Get(1), mat.Epsilon)
	assert.InEpsilon(t, 0.2855, color.Get(2), mat.Epsilon)
}

// Page 97
func TestColorWhenCastWithinSphereAtInsideSphere(t *testing.T) {
	w := mat.NewDefaultWorld()
	w.Objects[0].SetMaterial(mat.NewMaterial(mat.NewColor(0.8, 1.0, 0.6), 1.0, 0.7, 0.2, 200))
	w.Objects[1].SetMaterial(mat.NewMaterial(mat.NewColor(0.8, 1.0, 0.6), 1.0, 0.7, 0.2, 200))
	rc := New(w)
	r := mat.NewRay(mat.NewPoint(0, 0, 0.75), mat.NewVector(0, 0, -1))
	color := rc.colorAt(r, 1, 5)
	assert.InEpsilon(t, w.Objects[1].GetMaterial().Color.Get(0), color.Get(0), mat.Epsilon)
	assert.InEpsilon(t, w.Objects[1].GetMaterial().Color.Get(1), color.Get(1), mat.Epsilon)
	assert.InEpsilon(t, w.Objects[1].GetMaterial().Color.Get(2), color.Get(2), mat.Epsilon)
}

func TestPointNotInShadow(t *testing.T) {
	w := mat.NewDefaultWorld()
	rc := New(w)
	p := mat.NewPoint(0, 10, 10)
	assert.False(t, rc.pointInShadow(w.Light[0], p))
}
func TestPointInShadow(t *testing.T) {
	w := mat.NewDefaultWorld()
	rc := New(w)
	p := mat.NewPoint(10, -10, 10)
	assert.True(t, rc.pointInShadow(w.Light[0], p))
}
func TestPointNotInShadowWhenBehindLight(t *testing.T) {
	w := mat.NewDefaultWorld()
	rc := New(w)
	p := mat.NewPoint(-20, 20, -20)
	assert.False(t, rc.pointInShadow(w.Light[0], p))
}
func TestPointNotInShadowWhenBehindPoint(t *testing.T) {
	w := mat.NewDefaultWorld()
	rc := New(w)
	p := mat.NewPoint(-2, 2, -2)
	assert.False(t, rc.pointInShadow(w.Light[0], p))
}

// Big one on page 114
func TestWorldWithShadowTest(t *testing.T) {
	w := mat.NewDefaultWorld()
	w.Light = []mat.Light{mat.NewLight(mat.NewPoint(0, 0, -10), mat.NewColor(1, 1, 1))}
	s := mat.NewSphere()
	w.Objects = append(w.Objects, s)
	s2 := mat.NewSphere()
	s2.Transform = mat.Multiply(s2.Transform, mat.Translate(0, 0, 10))
	w.Objects = append(w.Objects, s2)

	rc := New(w)

	r := mat.NewRay(mat.NewPoint(0, 0, 5), mat.NewVector(0, 0, 1))
	i := mat.NewIntersection(4, s2)
	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(i, r, &comps)
	color := rc.shadeHit(comps, 1, 1)
	color.Elems[3] = 1 // just a fix for me using Tuple4 to represent colors...
	assert.Equal(t, mat.NewColor(0.1, 0.1, 0.1), color)
}

func TestOpaqueRefraction(t *testing.T) {
	w := mat.NewDefaultWorld()
	rc := New(w)
	s1 := w.Objects[0]
	r := mat.NewRay(mat.NewPoint(0, 0, -5), mat.NewVector(0, 0, 1))
	xs := []mat.Intersection{
		mat.NewIntersection(4, s1), mat.NewIntersection(6, s1),
	}
	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs[0], r, &comps, xs...)
	color := rc.refractedColor(comps, 5, 5)
	assert.Equal(t, black, color)
}

func TestRefractiveColorAndMaxRecursionDepth(t *testing.T) {
	w := mat.NewDefaultWorld()
	s1 := w.Objects[0]
	material := mat.NewDefaultMaterial()
	material.Transparency = 1.0
	material.RefractiveIndex = 1.5
	s1.SetMaterial(material)
	rc := New(w)
	r := mat.NewRay(mat.NewPoint(0, 0, -5), mat.NewVector(0, 0, 1))
	xs := []mat.Intersection{
		mat.NewIntersection(4, s1), mat.NewIntersection(6, s1),
	}
	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs[0], r, &comps, xs...)
	color := rc.refractedColor(comps, 0, 0)
	assert.Equal(t, black, color)
}

func TestTotalInternalRefraction(t *testing.T) {
	w := mat.NewDefaultWorld()
	s1 := w.Objects[0]
	material := mat.NewDefaultMaterial()
	material.Transparency = 1.0
	material.RefractiveIndex = 1.5
	s1.SetMaterial(material)

	rc := New(w)

	r := mat.NewRay(mat.NewPoint(0, 0, math.Sqrt(2)/2), mat.NewVector(0, 1, 0))

	xs := []mat.Intersection{
		mat.NewIntersection(-math.Sqrt(2)/2, s1), mat.NewIntersection(math.Sqrt(2)/2, s1),
	}

	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs[1], r, &comps, xs...)
	color := rc.refractedColor(comps, 5, 5)
	assert.Equal(t, black, color)
}

func TestRefractedColorWithRefractedRay(t *testing.T) {
	w := mat.NewDefaultWorld()
	s1 := w.Objects[0]
	material := mat.NewDefaultMaterial()
	material.Ambient = 1.0
	material.Pattern = mat.NewTestPattern()
	s1.SetMaterial(material)

	s2 := w.Objects[1]
	material2 := mat.NewDefaultMaterial()
	material2.Transparency = 1.0
	material2.RefractiveIndex = 1.5
	s2.SetMaterial(material2)
	rc := New(w)
	r := mat.NewRay(mat.NewPoint(0, 0, 0.1), mat.NewVector(0, 1, 0))

	xs := []mat.Intersection{
		mat.NewIntersection(-0.9899, s1),
		mat.NewIntersection(-0.4899, s2),
		mat.NewIntersection(0.4899, s2),
		mat.NewIntersection(0.4899, s1),
	}
	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs[2], r, &comps, xs...)
	color := rc.refractedColor(comps, 5, 5)
	assert.Equal(t, 0.0, color.Get(0))
	assert.InEpsilon(t, 0.99888, color.Get(1), mat.Epsilon*2)
	assert.InEpsilon(t, 0.04725, color.Get(2), mat.Epsilon*5)
}

func TestShadeHitWithRefractedMaterial(t *testing.T) {
	w := mat.NewDefaultWorld()
	floor := mat.NewPlane()
	floor.SetTransform(mat.Translate(0, -1, 0))
	mat1 := mat.NewDefaultMaterial()
	mat1.Transparency = 0.5
	mat1.RefractiveIndex = 1.5
	floor.SetMaterial(mat1)
	w.Objects = append(w.Objects, floor)

	ball := mat.NewSphere()
	mat2 := mat.NewDefaultMaterial()
	mat2.Color = mat.NewColor(1, 0, 0)
	mat2.Ambient = 0.5
	ball.SetMaterial(mat2)
	ball.SetTransform(mat.Translate(0, -3.5, -0.5))
	w.Objects = append(w.Objects, ball)

	rc := New(w)

	ray := mat.NewRay(mat.NewPoint(0, 0, -3), mat.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	xs := []mat.Intersection{
		mat.NewIntersection(math.Sqrt(2), floor),
	}
	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs[0], ray, &comps, xs...)
	color := rc.shadeHit(comps, 5, 5)
	assert.InEpsilon(t, 0.93642, color.Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.68642, color.Get(1), mat.Epsilon)
	assert.InEpsilon(t, 0.68642, color.Get(2), mat.Epsilon)
}

func TestShadeHitWhenBothTransparentAndRefractive(t *testing.T) {
	w := mat.NewDefaultWorld()
	r := mat.NewRay(mat.NewPoint(0, 0, -3), mat.NewVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))

	floor := mat.NewPlane()
	floor.SetTransform(mat.Translate(0, -1, 0))
	mat1 := mat.NewDefaultMaterial()
	mat1.Reflectivity = 0.5
	mat1.Transparency = 0.5
	mat1.RefractiveIndex = 1.5
	floor.SetMaterial(mat1)
	w.Objects = append(w.Objects, floor)

	ball := mat.NewSphere()
	ball.SetTransform(mat.Translate(0, -3.5, -0.5))
	mat2 := mat.NewDefaultMaterial()
	mat2.Color = mat.NewColor(1, 0, 0)
	mat2.Ambient = 0.5
	ball.SetMaterial(mat2)
	w.Objects = append(w.Objects, ball)

	rc := New(w)

	xs := []mat.Intersection{
		mat.NewIntersection(math.Sqrt(2), floor),
	}

	comps := mat.NewComputation()
	mat.PrepareComputationForIntersectionPtr(xs[0], r, &comps, xs...)
	color := rc.shadeHit(comps, 5, 5)
	assert.InEpsilon(t, 0.93391, color.Get(0), mat.Epsilon)
	assert.InEpsilon(t, 0.69643, color.Get(1), mat.Epsilon)
	assert.InEpsilon(t, 0.69243, color.Get(2), mat.Epsilon)
}
