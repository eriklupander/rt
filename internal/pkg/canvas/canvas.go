package canvas

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"math"
	"strconv"
	"strings"
)

type Canvas struct {
	W        int
	H        int
	MaxIndex int
	Pixels   []*mat.Tuple4
}

func NewCanvas(w int, h int) *Canvas {
	pixels := make([]*mat.Tuple4, w*h)
	for i, _ := range pixels {
		pixels[i] = mat.NewColor(0, 0, 0)
	}
	return &Canvas{W: w, H: h, Pixels: pixels, MaxIndex: w * h}
}

func (c *Canvas) WritePixel(col, row int, color *mat.Tuple4) {
	if row < 0 || col < 0 || row >= c.H || col > c.W {
		fmt.Println("pixel was out of bounds")
		return
	}
	if row*col > c.MaxIndex {
		fmt.Println("pixel was out of bounds total bounds")
		return
	}
	c.Pixels[c.toIdx(col, row)] = color
}

func (c *Canvas) WritePixelToIndex(idx int, color *mat.Tuple4) {
	c.Pixels[idx] = color
}

func (c *Canvas) ColorAt(col, row int) *mat.Tuple4 {
	return c.Pixels[c.toIdx(col, row)]
}

func (c *Canvas) ToPPM() string {

	final := ""

	for row := 0; row < c.H; row++ {
		buf := strings.Builder{}
		written := 0
		for col := 0; col < c.W; col++ {
			clamp(c.Pixels[c.toIdx(col, row)], &buf, &written)
		}
		out := buf.String()
		out = strings.TrimSuffix(out, " ")
		final += out + "\n"
	}
	final = strings.TrimSuffix(final, "\n")
	return fmt.Sprintf("P3\n%d %d\n255\n", c.W, c.H) + final
}

func clamp(color *mat.Tuple4, buf *strings.Builder, written *int) {
	for i := 0; i < 3; i++ {
		c := color.Get(i) * 255.0
		rounded := math.Round(c)
		if rounded > 255.0 {
			rounded = 255.0
		} else if rounded < 0.0 {
			rounded = 0.0
		}
		if *written+3 > 69 {
			buf.WriteString("\n")
			buf.WriteString(strconv.Itoa(int(rounded)))
			buf.WriteString(" ")
			*written = 3
		} else {
			buf.WriteString(strconv.Itoa(int(rounded)))
			if *written+6 < 69 {
				buf.WriteString(" ")
			}
			*written = *written + 4
		}

	}
}

func (c *Canvas) toIdx(x, y int) int {
	return y*c.W + x
}
