package helper

import (
	"fmt"
	"github.com/eriklupander/rt/internal/pkg/mat"
	"io/ioutil"
	"os"
	"strings"
)

func DumpGraph(g *mat.Group) {
	buf := strings.Builder{}
	buf.WriteString(`digraph G {`)
	dump(g, &buf)
	buf.WriteString(`}`)

	ioutil.WriteFile("nodegraph.viz", []byte(buf.String()), os.FileMode(0644))
}

func dump(g *mat.Group, buf *strings.Builder) {
	cnt := 0
	for i := range g.Children {
		ch := g.Children[i]
		switch val := ch.(type) {
		case *mat.Group:
			buf.WriteString(fmt.Sprintf(`"Grp-%x" -> "Grp-%x"`, g.ID(), val.ID()))
			dump(val, buf)
		default:
			cnt++
		}
	}
	buf.WriteString(fmt.Sprintf(`"Grp-%x" -> "Tri-%x-%v"`, g.ID(), g.ID(), cnt))
}
