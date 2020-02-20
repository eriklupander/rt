package mat

import (
	"github.com/eriklupander/rt/internal/pkg/identity"
	"math"
)

type Camera struct {
	Width      int
	Height     int
	Fov        float64
	Transform  Mat4x4
	Inverse    Mat4x4
	PixelSize  float64
	HalfWidth  float64
	HalfHeight float64
}

func NewCamera(width int, height int, fov float64) Camera {
	// Get the length of half the opposite part of the triangle
	halfView := math.Tan(fov / 2)
	aspect := float64(width) / float64(height)
	var halfWidth, halfHeight float64
	if aspect >= 1.0 {
		halfWidth = halfView
		halfHeight = halfView / aspect
	} else {
		halfWidth = halfView * aspect
		halfHeight = halfView
	}
	pixelSize := (halfWidth * 2) / float64(width)

	//transform := IdentityMatrix //make([]float64, 16)
	//copy(transform, IdentityMatrix.Elems)

	//inverse := IdentityMatrix //make([]float64, 16)
	//copy(inverse, IdentityMatrix.Elems)
	return Camera{
		Width:      width,
		Height:     height,
		Fov:        fov,
		Transform:  identity.Matrix,
		Inverse:    identity.Matrix,
		PixelSize:  pixelSize,
		HalfWidth:  halfWidth,
		HalfHeight: halfHeight,
	}
}

func ViewTransform(from, to, up Tuple4) Mat4x4 {
	// Create a new matrix from the identity matrix.
	vt := IdentityMatrix //Mat4x4{Elems: make([]float64, 16)}
	//copy(vt.Elems, IdentityMatrix.Elems)

	// Sub creates the initial vector between the eye and what we're looking at.
	forward := Normalize(Sub(to, from))

	// Normalize the up vector
	upN := Normalize(up)

	// Use the cross product to get the "third" axis (in this case, not the forward or up one)
	left := Cross(forward, upN)

	// Again, use cross product between the just computed left and forward to get the "true" up.
	trueUp := Cross(left, forward)

	// copy each axis into the matrix
	vt[0] = left.Get(0)
	vt[1] = left.Get(1)
	vt[2] = left.Get(2)

	vt[4] = trueUp.Get(0)
	vt[5] = trueUp.Get(1)
	vt[6] = trueUp.Get(2)

	vt[8] = -forward.Get(0)
	vt[9] = -forward.Get(1)
	vt[10] = -forward.Get(2)

	// finally, move the view matrix opposite the camera position to emulate that the camera has moved.
	return Multiply(vt, Translate(-from.Get(0), -from.Get(1), -from.Get(2)))
}

//
//func Render(c Camera, w World) *Canvas {
//
//	st := time.Now()
//	canvas := NewCanvas(c.Width, c.Height)
//	for row := 0; row < c.Height; row++ {
//		for col := 0; col < c.Width; col++ {
//			ray := RayForPixel(c, col, row)
//			color := ColorAt(w, ray, 5, 5)
//			canvas.WritePixel(col, row, color)
//		}
//		fmt.Printf("%d / %d\n", row+1, c.Height)
//	}
//	stats := runtime.MemStats{}
//	runtime.ReadMemStats(&stats)
//	fmt.Printf("Memory: %v ", bytesize.New(float64(stats.Alloc)).String())
//	fmt.Printf("Mallocs: %v ", stats.Mallocs)
//	fmt.Printf("Total alloc: %v\n", bytesize.New(float64(stats.TotalAlloc)).String())
//	fmt.Printf("%v", time.Now().Sub(st))
//	return canvas
//}
//
//func RenderThreaded(c Camera, w World) *Canvas {
//	st := time.Now()
//	canvas := NewCanvas(c.Width, c.Height)
//	jobs := make(chan *job)
//	wg := sync.WaitGroup{}
//	stats := runtime.MemStats{}
//	wg.Add(canvas.W * canvas.H)
//	for i := 0; i < 8; i++ {
//		go workerFuncPerPixel(canvas, c, w, jobs, &wg)
//	}
//	for row := 0; row < c.Height; row++ {
//		for col := 0; col < c.Width; col++ {
//			jobs <- &job{row: row, col: col}
//		}
//		fmt.Printf("%d/%d\n", row, c.Height)
//	}
//	wg.Wait()
//	fmt.Println("All done")
//	runtime.ReadMemStats(&stats)
//	fmt.Printf("Memory: %v ", bytesize.New(float64(stats.Alloc)).String())
//	fmt.Printf("Mallocs: %v ", stats.Mallocs)
//	fmt.Printf("Total alloc: %v\n", bytesize.New(float64(stats.TotalAlloc)).String())
//	fmt.Printf("%v", time.Now().Sub(st))
//	return canvas
//}
//
//func workerFuncPerPixel(canvas *Canvas, c Camera, w World, jobs chan *job, wg *sync.WaitGroup) {
//	for job := range jobs {
//		ray := RayForPixel(c, job.col, job.row)
//		color := ColorAt(w, ray, 5, 5)
//		canvas.WritePixelMutex(job.col, job.row, color)
//		wg.Done()
//	}
//}
//
//type job struct {
//	row int
//	col int
//}
