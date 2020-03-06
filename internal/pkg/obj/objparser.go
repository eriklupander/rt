package obj

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"io/ioutil"
	"strconv"
	"strings"
)

func ParseObj(data string) *Obj {
	out := &Obj{
		Verticies: make([]mat.Tuple4, 0),
		Groups:    make(map[string]*mat.Group),
	}

	var mats map[string]*mat.Mtl

	// fill index 0 with placeholder
	out.Verticies = append(out.Verticies, mat.NewPoint(0, 0, 0))
	out.Normals = append(out.Normals, mat.NewVector(0, 0, 0))
	rows := strings.Split(data, "\n")
	var currentGroup = "DefaultGroup"
	var currentMaterial = mat.NewDefaultMaterial()
	out.Groups[currentGroup] = mat.NewGroup()
	out.Groups[currentGroup].Label = currentGroup

	for _, row := range rows {
		if strings.TrimSpace(row) != "" {
			parts := strings.Fields(strings.TrimSpace(row))
			switch parts[0] {
			case "mtllib":
				fileName := parts[1]
				matData, err := ioutil.ReadFile("./assets/models/" + fileName)
				if err != nil {
					panic(err.Error())
				}
				mats = ParseMtl(string(matData))
			case "usemtl":
				currentMaterial = toMaterial(mats[parts[1]])
				fmt.Printf("Set material '%v' on object '%v'\n", mats[parts[1]].Name, currentGroup)
			case "v":

				x, _ := strconv.ParseFloat(parts[1], 64)
				y, _ := strconv.ParseFloat(parts[2], 64)
				z, _ := strconv.ParseFloat(parts[3], 64)
				out.Verticies = append(out.Verticies, mat.NewPoint(x, y, z))

			case "vn":
				x, _ := strconv.ParseFloat(parts[1], 64)
				y, _ := strconv.ParseFloat(parts[2], 64)
				z, _ := strconv.ParseFloat(parts[3], 64)
				out.Normals = append(out.Normals, mat.NewVector(x, y, z))

			case "f":

				if len(out.Normals) == 1 {
					for i := 2; i < len(parts)-1; i++ {
						idx1, _ := strconv.Atoi(parts[1])
						idx2, _ := strconv.Atoi(parts[i])
						idx3, _ := strconv.Atoi(parts[i+1])
						tri := mat.NewTriangle(
							out.Verticies[idx1],
							out.Verticies[idx2],
							out.Verticies[idx3])
						out.Groups[currentGroup].AddChild(tri)
					}
				} else {

					for i := 2; i < len(parts)-1; i++ {
						subparts1 := strings.Split(parts[1], "/")
						subparts2 := strings.Split(parts[i], "/")
						subparts3 := strings.Split(parts[i+1], "/")

						idx1, _ := strconv.Atoi(subparts1[0])
						idx2, _ := strconv.Atoi(subparts2[0])
						idx3, _ := strconv.Atoi(subparts3[0])

						normIdx1, _ := strconv.Atoi(subparts1[2])
						normIdx2, _ := strconv.Atoi(subparts2[2])
						normIdx3, _ := strconv.Atoi(subparts3[2])

						tri := mat.NewSmoothTriangle(
							out.Verticies[idx1],
							out.Verticies[idx2],
							out.Verticies[idx3],
							out.Normals[normIdx1],
							out.Normals[normIdx2],
							out.Normals[normIdx3])
						tri.Material = currentMaterial
						out.Groups[currentGroup].AddChild(tri)
					}
				}
			case "g":
				fallthrough
			case "o":
				currentGroup = strings.Fields(strings.TrimSpace(row))[1]
				if _, exists := out.Groups[currentGroup]; !exists {
					out.Groups[currentGroup] = mat.NewGroup()
					if len(parts) > 1 {
						out.Groups[currentGroup].Label = parts[1]
					}
				}
			default:
				out.IgnoredLines++
			}
		} else {
			out.IgnoredLines++
		}
	}
	tris := 0
	for i := range out.Groups {
		tris += len(out.Groups[i].Children)
	}
	fmt.Println("Loaded object:")
	fmt.Printf("Groups:    %d\n", len(out.Groups))
	fmt.Printf("Triangles: %d\n", tris)
	fmt.Printf("Verticies: %d\n", len(out.Verticies))
	fmt.Printf("Normals:   %d\n", len(out.Normals))
	return out
}

// toMaterial is a temp fix to convert our legacy materials as MTL materials
func toMaterial(mtl *mat.Mtl) mat.Material {
	m := mat.Material{}
	m.Name = mtl.Name

	r := (mtl.Ambient[0] + mtl.Diffuse[0] + mtl.Specular[0])
	g := (mtl.Ambient[1] + mtl.Diffuse[1] + mtl.Specular[1])
	b := (mtl.Ambient[2] + mtl.Diffuse[2] + mtl.Specular[2])
	m.Color = mat.NewColor(r, g, b)
	m.Ambient = avg(mtl.Ambient)
	m.Diffuse = avg(mtl.Diffuse)
	m.Specular = avg(mtl.Specular)
	m.Transparency = mtl.Transparency
	m.RefractiveIndex = mtl.RefractiveIndex
	m.Shininess = mtl.Shininess
	return m
}
func avg(t mat.Tuple4) float64 {
	return (t[1] + t[2]) / 2
}

type Obj struct {
	Verticies    []mat.Tuple4
	Normals      []mat.Tuple4
	Groups       map[string]*mat.Group
	IgnoredLines int
}

func (o *Obj) ToGroup() *mat.Group {
	g := mat.NewGroup()
	g.Label = "ROOT"
	for _, v := range o.Groups {
		g.AddChild(v)
	}
	return g
}

func (o *Obj) DefaultGroup() *mat.Group {
	return o.Groups["DefaultGroup"]
}

/*
Ns 96.078431
Ka 0.000000 0.000000 0.000000
Kd 0.000000 0.429367 0.640000
Ks 0.500000 0.500000 0.500000
Ni 1.000000
d 1.000000
illum 2
*/
func ParseMtl(data string) map[string]*mat.Mtl {
	rows := strings.Split(data, "\n")
	out := make(map[string]*mat.Mtl)

	var current string
	for _, row := range rows {
		if strings.TrimSpace(row) != "" {
			parts := strings.Fields(strings.TrimSpace(row))
			switch parts[0] {
			case "newmtl":
				name := parts[1]
				m := &mat.Mtl{}
				m.Name = name
				out[name] = m
				current = name
			case "Ns":
				out[current].Shininess, _ = strconv.ParseFloat(parts[1], 64)
			case "Ka":
				r, _ := strconv.ParseFloat(parts[1], 64)
				g, _ := strconv.ParseFloat(parts[2], 64)
				b, _ := strconv.ParseFloat(parts[3], 64)
				out[current].Ambient = mat.NewColor(r, g, b)
			case "Kd":
				r, _ := strconv.ParseFloat(parts[1], 64)
				g, _ := strconv.ParseFloat(parts[2], 64)
				b, _ := strconv.ParseFloat(parts[3], 64)
				out[current].Diffuse = mat.NewColor(r, g, b)
			case "Ks":
				r, _ := strconv.ParseFloat(parts[1], 64)
				g, _ := strconv.ParseFloat(parts[2], 64)
				b, _ := strconv.ParseFloat(parts[3], 64)
				out[current].Specular = mat.NewColor(r, g, b)
			case "Ni":
				out[current].RefractiveIndex, _ = strconv.ParseFloat(parts[1], 64)
			case "d":
				n, _ := strconv.ParseFloat(parts[1], 64)
				out[current].Transparency = 1 - n
			default:
				// ignore..
			}
		}
	}
	return out
}
