package canvas

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewCanvas(t *testing.T) {
	canvas := NewCanvas(10, 20)
	assert.Equal(t, 10, canvas.W)
	assert.Equal(t, 20, canvas.H)
	for _, px := range canvas.Pixels {
		assert.True(t, px.Get(0) == 0.0)
		assert.True(t, px.Get(1) == 0.0)
		assert.True(t, px.Get(2) == 0.0)
	}
}

func TestCanvas_WritePixel(t *testing.T) {
	canvas := NewCanvas(10, 20)
	canvas.WritePixel(2, 3, mat.NewColor(1, 0, 0))
	px := canvas.ColorAt(2, 3)
	assert.True(t, px.Get(0) == 1.0)
	assert.True(t, px.Get(1) == 0.0)
	assert.True(t, px.Get(2) == 0.0)
}

func TestCanvas_ToPPMHeader(t *testing.T) {
	canvas := NewCanvas(5, 3)
	ppmStr := canvas.ToPPM()
	assert.True(t, strings.HasPrefix(ppmStr, `P3
5 3
255`))
}

func TestCanvas_ToPPM(t *testing.T) {
	canvas := NewCanvas(5, 3)
	c1 := mat.NewColor(1.5, 0, 0)
	c2 := mat.NewColor(0, 0.5, 0)
	c3 := mat.NewColor(-0.5, 0, 1)
	canvas.WritePixel(0, 0, c1)
	canvas.WritePixel(2, 1, c2)
	canvas.WritePixel(4, 2, c3)

	ppmData := canvas.ToPPM()

	assert.Equal(t, `P3
5 3
255
255 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 128 0 0 0 0 0 0 0
0 0 0 0 0 0 0 0 0 0 0 0 0 0 255`, ppmData)
}

func TestCanvas_ToPPM70(t *testing.T) {
	expected := `P3
10 2
255
255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204
153 255 204 153 255 204 153 255 204 153 255 204 153
255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204
153 255 204 153 255 204 153 255 204 153 255 204 153`

	c1 := mat.NewColor(1, 0.8, 0.6)
	canvas := NewCanvas(10, 2)
	for i:=0; i < 10*2; i++ {
		canvas.WritePixelToIndex(i, c1)
	}

	assert.Equal(t, expected, canvas.ToPPM())
}