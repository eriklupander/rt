package gui

type Pixel struct {
	X int
	Y int
	R int
	G int
	B int
}

var PixelChan = make(chan Pixel, 1000000)