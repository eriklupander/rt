package persist

import (
	"encoding/json"
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"testing"
)

func TestWriteWorld(t *testing.T) {
	s1 := mat.NewSphere()
	s1.SetTransform(mat.Multiply(mat.Translate(-2, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))

	data, err := json.Marshal(s1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(string(data))

	gr := mat.NewGroup()
	gr.AddChild(s1)
	data, err = json.Marshal(gr)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(string(data))

	//light := mat.NewLight(mat.NewPoint(-5, 2.5, -3), mat.NewColor(1, 1, 1))
	//w := mat.NewDefaultWorld()
	//w.Light = append(w.Light, light)
	//gr := mat.NewGroup()
	//gr.AddChild(s1)
	//w.Objects = append(w.Objects, gr)
	//
	//gob.Register(s1)
	//gob.Register(light)
	//gob.Register(w)
	//gob.Register(gr)
	//buf := new(bytes.Buffer)
	//encoder := gob.NewEncoder(buf)
	//
	//err := encoder.Encode(w)
	//assert.NoError(t, err)
	//assert.True(t, len(buf.Bytes()) > 0)
}

func createWorldForTest() mat.World {
	w := mat.NewWorld()
	light := mat.NewLight(mat.NewPoint(-5, 2.5, -3), mat.NewColor(1, 1, 1))
	w.Light = append(w.Light, light)

	//  Cylinder
	cyl := mat.NewCylinderMMC(0.0, 3.0, true)
	cyl.SetTransform(mat.Translate(1, 0.0, -1)) //1, 0.25, -1
	cyl.SetTransform(mat.Scale(0.2, 0.4, 0.2))
	cyl.Material.Color = mat.NewColor(1, 0.88, 0.63)
	w.Objects = append(w.Objects, cyl)

	gr := mat.NewGroup()

	s1 := mat.NewSphere()
	s1.SetTransform(mat.Multiply(mat.Translate(-2, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
	gr.AddChild(s1)

	s2 := mat.NewSphere()
	s2.SetTransform(mat.Multiply(mat.Translate(-1, 0.25, -1), mat.Scale(0.25, 0.25, 0.25)))
	gr.AddChild(s2)

	s3 := mat.NewSphere()
	gr.AddChild(s3)

	gr.SetTransform(mat.RotateY(0.67))
	gr.Bounds() // For now, important to always call Bounds on Group once set up.
	mat.Divide(gr, 1)
	w.Objects = append(w.Objects, gr)

	return w
}
