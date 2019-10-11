package api

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ak1ra24/microns/utils"
	"github.com/docker/docker/client"
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

// MicronstoD3Handler func is microns config yaml to d3 json
func MicronstoD3Handler(w http.ResponseWriter, r *http.Request) {

	// filepath := "/home/akira/go/src/github.com/ak1ra24/microns/examples/clos/config.yaml"
	// filepath := "./examples/clos/config.yaml"
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

	var d3jsonData D3json
	d3jsonData.D3Nodes = d3nodes
	d3jsonData.D3Links = d3links

	d3json, err := json.Marshal(d3jsonData)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(d3json)

}

// JsonResHandler func is d3 force layout json response adn create frontend server
func JsonResHandler(filepath string) {
	cfgFile = filepath
	fmt.Println("config File: ", cfgFile)
	ok := Confirm("\x1b[31mConfirm use port number 8080 for frontend nginx container and 8000 for backend webapp\x1b[0m")
	if ok {
		if err := CreateFrontend(); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Please Access FrontEnd: localhost:8080, BackEnd: localhost:8000/api/d3")

		http.HandleFunc("/api/d3", MicronstoD3Handler)
		http.ListenAndServe(":8000", nil)
	} else {
		fmt.Println("Finish webview command")
	}
}

// CreateFrontend func is Image: akiranet24/microns-frontend docker container start
func CreateFrontend() error {

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	c := NewContainer(ctx, cli)

	if err := c.CreateContainerPort("akiranet24/microns-frontend", "microns-frontend", "80", "8080"); err != nil {
		return err
	}

	return nil
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
