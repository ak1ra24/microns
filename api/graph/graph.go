package graph

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ak1ra24/microns/api/utils"
	"github.com/awalterschulze/gographviz"
)

func Graph(nodes []utils.Node, filename string) {
	g := gographviz.NewGraph()
	if err := g.SetName("G"); err != nil {
		panic(err)
	}

	if err := g.SetDir(true); err != nil {
		panic(err)
	}

	if err := g.AddAttr("G", "bgcolor", "\"#343434\""); err != nil {
		panic(err)
	}

	if err := g.AddAttr("G", "layout", "neato"); err != nil {
		panic(err)
	}

	// node setting
	nodeAttrs := make(map[string]string)
	nodeAttrs["colorscheme"] = "rdylgn11"
	nodeAttrs["style"] = "\"solid,filled\""
	nodeAttrs["fontcolor"] = "6"
	nodeAttrs["fontname"] = "\"Migu 1M\""
	nodeAttrs["color"] = "7"
	nodeAttrs["fillcolor"] = "11"
	nodeAttrs["shape"] = "doublecircle"

	var Nodes []string
	var Links []string
	for _, node := range nodes {
		Nodes = append(Nodes, node.Name)
		for _, inf := range node.Interface {
			link := node.Name + "-" + inf.Name
			fmt.Println(link)
			Links = append(Links, link)
		}
	}
	fmt.Println(Nodes)
	fmt.Println(Links)

	for _, Node := range Nodes {
		if err := g.AddNode("G", Node, nodeAttrs); err != nil {
			panic(err)
		}
	}

	edgeAttrs := make(map[string]string)
	edgeAttrs["color"] = "white"

	// Edge(Node関係の関係)追加
	// src(第1引数)からdst(第2引数)へ方向を付ける
	for _, Link := range Links {
		edge := strings.Split(Link, "-")
		edge1 := edge[0]
		edge2 := edge[1]
		if err := g.AddEdge(edge1, edge2, true, edgeAttrs); err != nil {
			panic(err)
		}
	}

	// dotファイル出力
	s := g.String()
	dotfile := filename + ".dot"
	file, err := os.Create(dotfile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(s))
}

func DottoPng(filename string) {
	dotfile := filename + ".dot"
	pngfile := filename + ".png"
	err := exec.Command("dot", "-T", "png", dotfile, "-o", pngfile).Run()
	if err != nil {
		panic(err)
	}
	if err := os.Remove(dotfile); err != nil {
		panic(err)
	}
}

func main() {
	filename := "hai"
	nodes := utils.ParseYaml("./config.yaml")
	Graph(nodes, filename)
	DottoPng(filename)
}
