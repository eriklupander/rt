package render

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/inhies/go-bytesize"
	"runtime"
	"sync"
	"time"
)

type Context struct {
	world  mat.World
	camera mat.Camera
	canvas *mat.Canvas
	jobs   chan *job
	wg     *sync.WaitGroup
}

func NewContext(world mat.World, camera mat.Camera, canvas *mat.Canvas, jobs chan *job, wg *sync.WaitGroup) Context {
	return Context{
		world: world,
		camera: camera,
		canvas: canvas,
		jobs: jobs,
		wg: wg}
}

var lock sync.Mutex

func Threaded(c []mat.Camera, w []mat.World) *mat.Canvas {
	st := time.Now()
	canvas := mat.NewCanvas(c[0].Width, c[0].Height)
	jobs := make(chan *job)

	wg := sync.WaitGroup{}
	wg.Add(canvas.W * canvas.H)

	// allocate GOMAXPROCS render Contexts
	var GOMAXPROCS = 1
	renderContexts := make([]Context, GOMAXPROCS)
	for i := 0; i < GOMAXPROCS; i++ {

		// DO WE NEED TO DEEP-COPY the world and camera?
		renderContexts[i] = NewContext(w[i], c[i], canvas, jobs, &wg)
	}

	// start workers
	for i := 0; i < GOMAXPROCS; i++ {
		go renderContexts[i].workerFuncPerPixel()
	}

	// start passing work to the workers, one pixel at a time
	for row := 0; row < c[0].Height; row++ {
		for col := 0; col < c[0].Width; col++ {
			jobs <- &job{row: row, col: col}
		}
		fmt.Printf("%d/%d\n", row, c[0].Height)
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
		ray := mat.RayForPixel(rc.camera, job.col, job.row)
		color := mat.ColorAt(rc.world, ray, 5, 5)
		lock.Lock()
		rc.canvas.WritePixelMutex(job.col, job.row, color)
		lock.Unlock()
		rc.wg.Done()
	}
}

type job struct {
	row int
	col int
}
