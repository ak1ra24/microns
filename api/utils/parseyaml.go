package utils

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Data struct {
	Nodes []Node `yaml:"nodes"`
}

type Node struct {
	Name      string      `yaml:"name"`
	Image     string      `yaml:"image"`
	Interface []InterFace `yaml:"interfaces"`
	Volume    string      `yaml:"volume"`
}

type InterFace struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Ipv4 string `yaml:"ipv4"`
	Ipv6 string `yaml:"ipv6"`
	Args string `yaml:"args"`
}

func ParseYaml(filepath string) []Node {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	fmt.Println("----------------------------------------------")
	fmt.Println("                   CONFIG                     ")
	fmt.Println("----------------------------------------------")
	fmt.Printf("%+v\n", string(buf))

	var nodes Data

	err = yaml.Unmarshal(buf, &nodes)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("nodes: %+v\n", nodes)

	return nodes.Nodes
}
