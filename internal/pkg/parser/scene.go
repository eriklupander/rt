package parser

import "github.com/eriklupander/rt/internal/pkg/mat"

type Scene struct {
	Lights     []mat.Light
	World      *mat.World
	Camera     *mat.Camera
	Materials  map[string]mat.Material
	Transforms map[string][]mat.Mat4x4
}

func NewScene() *Scene {
	w := mat.NewWorld()
	return &Scene{
		World:      &w,
		Lights:     make([]mat.Light, 0),
		Materials:  make(map[string]mat.Material),
		Transforms: make(map[string][]mat.Mat4x4),
	}
}
