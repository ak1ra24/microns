package graph

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ak1ra24/microns/utils"
	"github.com/awalterschulze/gographviz"
)

func Graph(nodes []utils.Node, bridges []utils.Switch, filename string) {
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
	var nodes_bridges []string

	Link := make(map[string]string)
	Addrv4 := make(map[string]string)
	Addrv6 := make(map[string]string)
	for _, node := range nodes {
		Nodes = append(Nodes, node.Name)
		for _, inf := range node.Interface {
			link := node.Name + "-" + inf.InfName
			peer := inf.PeerNode + "-" + inf.PeerInf
			if len(inf.Ipv4) != 0 && len(inf.Ipv6) != 0 {
				Addrv4[link] = fmt.Sprintf("%s", inf.Ipv4)
				Addrv6[link] = fmt.Sprintf("%s", inf.Ipv6)
			} else if len(inf.Ipv4) != 0 && len(inf.Ipv6) == 0 {
				Addrv4[link] = fmt.Sprintf("%s", inf.Ipv4)
			} else if len(inf.Ipv4) == 0 && len(inf.Ipv6) != 0 {
				Addrv6[link] = fmt.Sprintf("%s", inf.Ipv6)
			}
			Link[link] = peer
			Links = append(Links, link)
			nodes_bridges = append(nodes_bridges, node.Name)
		}
	}

	for _, bridge := range bridges {
		nodes_bridges = append(nodes_bridges, bridge.Name)
	}

	for _, Node := range nodes_bridges {
		if err := g.AddNode("G", Node, nodeAttrs); err != nil {
			panic(err)
		}
	}

	edgeAttrs := make(map[string]string)
	// edgeAttrs["color"] = "white"

	// Edge(Node関係の関係)追加
	// src(第1引数)からdst(第2引数)へ方向を付ける
	for _, linkinfo := range Links {
		edge11 := strings.Split(linkinfo, "-")
		edge22 := strings.Split(Link[linkinfo], "-")
		edge1 := edge11[0]
		edge2 := edge22[0]

		var addrv4info, addrv6info string
		var addrv4, addrv6 []string
		if len(Addrv4[linkinfo]) != 0 && len(Addrv6[linkinfo]) != 0 {
			addrv4 = strings.Split(Addrv4[linkinfo], "-")
			addrv4info = addrv4[0]
			addrv6 = strings.Split(Addrv6[linkinfo], "-")
			addrv6info = addrv6[0]
			edgeAttrs["label"] = fmt.Sprintf("\"%s, IPv4: %s, IPv6: %s\"", edge11[1], addrv4info, addrv6info)
		} else if len(Addrv4[linkinfo]) != 0 && len(Addrv6[linkinfo]) == 0 {
			addrv4 = strings.Split(Addrv4[linkinfo], "-")
			addrv4info = addrv4[0]
			edgeAttrs["label"] = fmt.Sprintf("\"%s, IPv4: %s\"", edge11[1], addrv4info)
		} else if len(Addrv4[linkinfo]) == 0 && len(Addrv6[linkinfo]) != 0 {
			addrv6 = strings.Split(Addrv6[linkinfo], "-")
			addrv6info = addrv6[0]
			edgeAttrs["label"] = fmt.Sprintf("\"%s, IPv6: %s\"", edge11[1], addrv6info)
		}

		if err := g.AddEdge(edge1, edge2, true, edgeAttrs); err != nil {
			panic(err)
		}
		delete(edgeAttrs, "label")
	}

	// dotファイル出力
	s := g.String()
	fmt.Println(s)
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
	// if err := os.Remove(dotfile); err != nil {
	// 	panic(err)
	// }
}
