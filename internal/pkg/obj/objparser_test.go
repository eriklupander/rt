package obj

import (
	"github.com/eriklupander/rt/internal/pkg/mat"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseGibberish(t *testing.T) {
	gibberish := `There was a young lady named Bright
who traveled much faster than light.
She set out one day
in a relative way,
and came back the previous night.`
	result := ParseObj(gibberish)
	assert.Equal(t, 5, result.IgnoredLines)
}

func TestParseVerticies(t *testing.T) {

	data := `
v -1 1 0
v -1.0000 0.5000 0.0000
v 1 0 0
v 1 1 0
`
	res := ParseObj(data)
	assert.Equal(t, mat.NewPoint(-1, 1, 0), res.Verticies[1])
	assert.Equal(t, mat.NewPoint(-1, 0.5, 0), res.Verticies[2])
	assert.Equal(t, mat.NewPoint(1, 0, 0), res.Verticies[3])
	assert.Equal(t, mat.NewPoint(1, 1, 0), res.Verticies[4])
}

func TestParseTriangleFaces(t *testing.T) {
	data := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
f 1 2 3
f 1 3 4
`
	parser := ParseObj(data)
	gr := parser.DefaultGroup()
	t1 := gr.Children[0].(*mat.Triangle)
	t2 := gr.Children[1].(*mat.Triangle)
	assert.Equal(t, t1.P1, parser.Verticies[1])
	assert.Equal(t, t1.P2, parser.Verticies[2])
	assert.Equal(t, t1.P3, parser.Verticies[3])
	assert.Equal(t, t2.P1, parser.Verticies[1])
	assert.Equal(t, t2.P2, parser.Verticies[3])
	assert.Equal(t, t2.P3, parser.Verticies[4])
}

func TestTriangulatePolygon(t *testing.T) {
	data := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
v 0 2 0
f 1 2 3 4 5`
	parser := ParseObj(data)
	gr := parser.DefaultGroup()
	t1 := gr.Children[0].(*mat.Triangle)
	t2 := gr.Children[1].(*mat.Triangle)
	t3 := gr.Children[2].(*mat.Triangle)

	assert.Equal(t, t1.P1, parser.Verticies[1])
	assert.Equal(t, t1.P2, parser.Verticies[2])
	assert.Equal(t, t1.P3, parser.Verticies[3])
	assert.Equal(t, t2.P1, parser.Verticies[1])
	assert.Equal(t, t2.P2, parser.Verticies[3])
	assert.Equal(t, t2.P3, parser.Verticies[4])
	assert.Equal(t, t3.P1, parser.Verticies[1])
	assert.Equal(t, t3.P2, parser.Verticies[4])
	assert.Equal(t, t3.P3, parser.Verticies[5])
}

func TestTrianglesInGroups(t *testing.T) {
	data := `
v -1 1 0
v -1 0 0
v 1 0 0
v 1 1 0
g FirstGroup
f 1 2 3
g SecondGroup
f 1 3 4`

	parser := ParseObj(data)
	gr1 := parser.Groups["FirstGroup"]
	gr2 := parser.Groups["SecondGroup"]
	t1 := gr1.Children[0].(*mat.Triangle)
	t2 := gr2.Children[0].(*mat.Triangle)

	assert.Equal(t, t1.P1, parser.Verticies[1])
	assert.Equal(t, t1.P2, parser.Verticies[2])
	assert.Equal(t, t1.P3, parser.Verticies[3])
	assert.Equal(t, t2.P1, parser.Verticies[1])
	assert.Equal(t, t2.P2, parser.Verticies[3])
	assert.Equal(t, t2.P3, parser.Verticies[4])

}

func TestNormalData(t *testing.T) {
	data := `
vn 0 0 1
vn 0.707 0 -0.707
vn 1 2 3`

	parser := ParseObj(data)
	assert.Equal(t, parser.Normals[1], mat.NewVector(0, 0, 1))
	assert.Equal(t, parser.Normals[2], mat.NewVector(0.707, 0, -0.707))
	assert.Equal(t, parser.Normals[3], mat.NewVector(1, 2, 3))

}

func TestFacesWithNormals(t *testing.T) {
	data := `
v 0 1 0
v -1 0 0
v 1 0 0
vn -1 0 0
vn 1 0 0
vn 0 1 0
f 1//3 2//1 3//2
f 1/0/3 2/102/1 3/14/2`
	parser := ParseObj(data)

	g := parser.DefaultGroup()
	t1 := g.Children[0].(*mat.SmoothTriangle)
	t2 := g.Children[1].(*mat.SmoothTriangle)
	assert.Equal(t, t1.P1, parser.Verticies[1])
	assert.Equal(t, t1.P2, parser.Verticies[2])
	assert.Equal(t, t1.P3, parser.Verticies[3])
	assert.Equal(t, t1.N1, parser.Normals[3])
	assert.Equal(t, t1.N2, parser.Normals[1])
	assert.Equal(t, t1.N3, parser.Normals[2])
	dr1 := *t1
	dr2 := *t2
	assert.True(t, reflect.DeepEqual(dr1, dr2))
}
