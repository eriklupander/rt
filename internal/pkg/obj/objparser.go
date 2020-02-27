package obj

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"strconv"
	"strings"
)

func ParseObj(data string) *Obj {
	out := &Obj{
		Verticies: make([]mat.Tuple4, 0),
		Groups:    make(map[string]*mat.Group),
	}
	// fill index 0 with placeholder
	out.Verticies = append(out.Verticies, mat.NewPoint(0, 0, 0))
	out.Normals = append(out.Normals, mat.NewVector(0, 0, 0))
	rows := strings.Split(data, "\n")
	var currentGroup = "DefaultGroup"
	out.Groups[currentGroup] = mat.NewGroup()
	for _, row := range rows {
		if strings.TrimSpace(row) != "" {
			parts := strings.Fields(strings.TrimSpace(row))
			//fmt.Printf("%v\n", row)
			switch parts[0] {
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
						out.Groups[currentGroup].AddChild(tri)
					}
				}
			case "g":
				currentGroup = strings.Fields(strings.TrimSpace(row))[1]
				if _, exists := out.Groups[currentGroup]; !exists {
					out.Groups[currentGroup] = mat.NewGroup()
				}
			default:
				out.IgnoredLines++
			}
		} else {
			out.IgnoredLines++
		}
	}
	tris := 0
	for _, v := range out.Groups {
		tris += len(v.Children)
	}
	fmt.Println("Loaded object:")
	fmt.Printf("Groups:    %d\n", len(out.Groups))
	fmt.Printf("Triangles: %d\n", tris)
	fmt.Printf("Verticies: %d\n", len(out.Verticies))
	fmt.Printf("Normals:   %d\n", len(out.Normals))
	return out
}

type Obj struct {
	Verticies    []mat.Tuple4
	Normals      []mat.Tuple4
	Groups       map[string]*mat.Group
	IgnoredLines int
}

func (o *Obj) ToGroup() *mat.Group {
	g := mat.NewGroup()
	for _, v := range o.Groups {
		g.AddChild(v)
	}
	return g
}

func (o *Obj) DefaultGroup() *mat.Group {
	return o.Groups["DefaultGroup"]
}
