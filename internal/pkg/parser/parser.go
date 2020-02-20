package parser

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"sort"
)

// ParseYAML is a WiP parser for loading scenes given a custom YAML format. The format is similar to, but not compatible
// with the yaml files the book sometimes uses.
func ParseYAML(file string) *Scene {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err.Error())
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &m)
	if err != nil {
		panic(err.Error())
	}

	scene := NewScene()

	sc := m["scene"].([]interface{})
	for _, v := range sc {
		switch val := v.(type) {
		case map[string]interface{}:
			//fmt.Printf("%+v\n", val)
			for k1, v := range val {
				switch k1 {
				case "camera":
					//fmt.Println("Add camera")
					AddCamera(scene, v.(map[string]interface{}))
				case "light":
					AddLight(scene, v.(map[string]interface{}))
				case "material":
					AddMaterial(scene, v.(map[string]interface{}))
				case "transform":
					AddTransform(scene, v.(map[string]interface{}))

				// primitives
				case "plane":
					AddPlane(scene, v.(map[string]interface{}))
				case "sphere":
					AddSphere(scene, v.(map[string]interface{}))
				}
			}
		case []interface{}:
			fmt.Printf("Wooot ? %v\n", val)
		}
	}

	return scene
}

func AddSphere(scene *Scene, v map[string]interface{}) {
	s := mat.NewSphere()
	processShape(s, scene, v)
	scene.World.Objects = append(scene.World.Objects, s)
}

func AddPlane(scene *Scene, v map[string]interface{}) {
	p := mat.NewPlane()
	processShape(p, scene, v)
	scene.World.Objects = append(scene.World.Objects, p)
}

func processShape(p mat.Shape, scene *Scene, v map[string]interface{}) {
	for name, val := range v {
		switch name {
		case "material":
			matM := val.(map[string]interface{})
			color := parseVector(matM["color"].([]interface{}))

			ambient := matM["ambient"].(float64)
			diffuse := matM["diffuse"].(float64)
			specular := matM["specular"].(float64)

			m := mat.NewMaterial(color, ambient, diffuse, specular, 200)
			if shininess, ok := matM["shininess"]; ok {
				m.Reflectivity = float64(shininess.(int))
			}
			if reflective, ok := matM["reflective"]; ok {
				m.Reflectivity = reflective.(float64)
			}
			if transparency, ok := matM["transparency"]; ok {
				m.Transparency = transparency.(float64)
			}
			if refractiveIndex, ok := matM["refractive-index"]; ok {
				m.RefractiveIndex = refractiveIndex.(float64)
			}

			p.SetMaterial(m)
		case "transform":
			switch transf := val.(type) {
			case string:
				tf := scene.Transforms[transf]
				for _, v := range tf {
					p.SetTransform(v)
				}
			case []interface{}:

				for _, tV := range transf {
					switch tInner := tV.(type) {
					case []interface{}:
						firstElem := tInner[0].(string)
						switch firstElem {
						case "rotate-x":
							x := mat.RotateX(tInner[1].(float64))
							p.SetTransform(x)
						case "translate":
							translate := mat.Translate(tInner[1].(float64), tInner[2].(float64), tInner[3].(float64))
							p.SetTransform(translate)
						}
					}

				}
			}
		}
		fmt.Printf("%v %v\n", name, val)
	}
}

func AddTransform(scene *Scene, v map[string]interface{}) {
	for name, val := range v {
		transforms := make([]DefinedTransform, 0)
		trf := val.(map[string]interface{})
		for transformType, transformVal := range trf {
			switch transformType {
			case "scale":
				p := parsePoint(transformVal.([]interface{}))
				transforms = append(transforms, DefinedTransform{name: "scale", transform: mat.Scale(p.Get(0), p.Get(1), p.Get(2))})
			case "translate":
				p := parsePoint(transformVal.([]interface{}))
				transforms = append(transforms, DefinedTransform{name: "translate", transform: mat.Translate(p.Get(0), p.Get(1), p.Get(2))})
			}
		}

		// sort transforms in translate, scale order so we can apply them in correct order.
		sort.Slice(transforms, func(i, j int) bool {
			return transforms[i].name > transforms[j].name
		})

		out := make([]mat.Mat4x4, len(transforms))
		for i := 0; i < len(out); i++ {
			out[i] = transforms[i].transform
		}

		scene.Transforms[name] = out
	}
}

type DefinedTransform struct {
	name      string
	transform mat.Mat4x4
}

func AddMaterial(scene *Scene, v map[string]interface{}) {
	for name, val := range v {
		mtl := val.(map[string]interface{})
		clr := parseVector(mtl["color"].([]interface{}))
		m := mat.Material{
			Color:           clr,
			Ambient:         mtl["ambient"].(float64),
			Diffuse:         mtl["diffuse"].(float64),
			Specular:        mtl["specular"].(float64),
			Shininess:       220,
			Reflectivity:    mtl["reflective"].(float64),
			Transparency:    0.0,
			RefractiveIndex: 1.0,
		}
		scene.Materials[name] = m
	}
}

func AddCamera(scene *Scene, v map[string]interface{}) {
	var width, height int
	var fov float64
	var from, to, up mat.Tuple4
	for k, v := range v {
		switch k {
		case "width":
			width = v.(int)
		case "height":
			height = v.(int)
		case "field-of-view":
			fov = v.(float64)
		case "from":
			from = parseVector(v.([]interface{}))
		case "to":
			to = parseVector(v.([]interface{}))
		case "up":
			up = parseVector(v.([]interface{}))
		}
	}
	cam := mat.NewCamera(width, height, fov)
	cam.Transform = mat.ViewTransform(from, to, up)
	scene.Camera = &cam
}

func AddLight(scene *Scene, v map[string]interface{}) {
	pos := parsePoint(v["at"].([]interface{}))
	intensity := parseVector(v["intensity"].([]interface{}))
	scene.Lights = append(scene.Lights, mat.NewLight(pos, intensity))
}

func parseVector(v []interface{}) mat.Tuple4 {
	vec := mat.NewVector(0, 0, 0)
	for i := range v {
		switch val := v[i].(type) {
		case float64:
			vec.Elems[i] = val
		case int:
			vec.Elems[i] = float64(val)
		}
	}
	return vec
}

func parsePoint(v []interface{}) mat.Tuple4 {
	p := mat.NewPoint(0, 0, 0)
	for i := range v {
		switch val := v[i].(type) {
		case float64:
			p.Elems[i] = val
		case int:
			p.Elems[i] = float64(val)
		}
	}
	return p
}
