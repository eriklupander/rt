package persist

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
)

func WriteWorld(w mat.World) ([]byte, error) {
	buf := new(bytes.Buffer)
	gob.Register(w)
	gob.Register(mat.NewLight(mat.NewPoint(0, 0, 0), mat.NewVector(0, 0, 0)))
	gob.Register(mat.NewCylinder())
	gob.Register(mat.NewGroup())
	gob.Register(mat.NewSphere())

	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(w)
	if err != nil {
		return nil, err
	}
	out := buf.Bytes()
	fmt.Printf("successfully serialized world to .gob, size: %v\n", len(out))
	return out, nil
}
