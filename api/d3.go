package api

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ak1ra24/microns/utils"
)

var (
	cfgFile string
)

// D3json d3node and d3link
type D3json struct {
	D3Nodes []D3Node `json:"nodes"`
	D3Links []D3Link `json:"links"`
}

// D3Node struct is d3 node
type D3Node struct {
	Id    int    `json:"id"`
	Label string `json:"label"`
	Types string `json:"types"`
}

// D3Link struct is d3 node link
type D3Link struct {
	Source int `json:"source"`
	Target int `json:"target"`
}

// HtmlHandler func is template html view
func HtmlHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./websrc/html/index.html.tpl"))

	filepath := cfgFile
	nodes := utils.ParseNodes(filepath)
	bridges := utils.ParseSwitch(filepath)

	var d3nodes []D3Node
	var d3links []D3Link

	for nodenum, node := range nodes {
		d3node := D3Node{Id: nodenum, Label: node.Name, Types: "router"}
		d3nodes = append(d3nodes, d3node)
	}

	nodesLength := len(nodes) - 1

	for _, bridge := range bridges {
		d3node := D3Node{Id: nodesLength + 1, Label: bridge.Name, Types: "bridge"}
		d3nodes = append(d3nodes, d3node)
	}

	for nodenum, node := range nodes {
		for _, inf := range node.Interface {
			for _, d3node := range d3nodes {
				if inf.PeerNode == d3node.Label {
					d3link := D3Link{Source: nodenum, Target: d3node.Id}
					d3links = append(d3links, d3link)
				}
			}
		}
	}

	d3jsonData := D3json{D3Nodes: d3nodes, D3Links: d3links}

	if err := t.Execute(w, d3jsonData); err != nil {
		log.Fatal(err)
	}

}

// JsonResHandler func is d3 force layout json response adn create frontend server
func JsonResHandler(filepath string) {
	cfgFile = filepath
	fmt.Println("config File: ", cfgFile)
	ok := Confirm("\x1b[31mConfirm use port number 8000 for webapp\x1b[0m")
	if ok {

		ipv4s, err := utils.GetIP()
		if err != nil {
			log.Fatal(err)
		}
        fmt.Printf("Please Access http://%s:8000/", ipv4s[0])

		http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("websrc/images/"))))
		http.HandleFunc("/", HtmlHandler)

		http.ListenAndServe(":8000", nil)
	} else {
		fmt.Println("Finish webview command")
	}
}

// Confirm func is Wait until yes or no is entered
func Confirm(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for i := 0; i < 3; i++ {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" || response == "Yes" {
			return true
		} else if response == "n" || response == "no" || response == "No" {
			return false
		}
	}

	return false
}

// CreateFrontend func is Image: akiranet24/microns-frontend docker container start
// func CreateFrontend() error {
//
// 	ctx := context.Background()
// 	cli, err := client.NewEnvClient()
// 	if err != nil {
// 		return err
// 	}
//
// 	c := NewContainer(ctx, cli)
//
// 	if err := c.CreateContainerPort("akiranet24/microns-frontend", "microns-frontend", "80", "8080"); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
