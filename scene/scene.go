package scene

import "github.com/eriklupander/rt/internal/pkg/mat"

type Scene struct {
	Camera     mat.Camera
	Lights     []mat.Light
	AreaLights []mat.AreaLight
	Objects    []mat.Shape
}
