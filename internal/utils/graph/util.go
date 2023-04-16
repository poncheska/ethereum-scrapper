package graph

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"math/big"
	"os"
)

func SaveScheme(path string, root string, in, out map[string]*big.Int) error {
	g := graph.New(graph.StringHash, graph.Directed())

	_ = g.AddVertex(root, graph.VertexAttribute("colorscheme", "blues3"), graph.VertexAttribute("style", "filled"), graph.VertexAttribute("color", "2"), graph.VertexAttribute("fillcolor", "1"))

	for node, val := range out {
		_ = g.AddVertex(node, graph.VertexAttribute("colorscheme", "greens3"), graph.VertexAttribute("style", "filled"), graph.VertexAttribute("color", "2"), graph.VertexAttribute("fillcolor", "1"))
		_ = g.AddEdge(root, node, graph.EdgeAttribute("label", fmt.Sprintf("%s wei", val.String())))
	}

	for node, val := range in {
		_ = g.AddVertex(node, graph.VertexAttribute("colorscheme", "greens3"), graph.VertexAttribute("style", "filled"), graph.VertexAttribute("color", "2"), graph.VertexAttribute("fillcolor", "1"))
		_ = g.AddEdge(node, root, graph.EdgeAttribute("label", fmt.Sprintf("%s wei", val.String())))
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	err = draw.DOT(g, file)
	if err != nil {
		return err
	}

	return nil
}
