package utils

import (
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
	Cmds      []Cmd       `yaml:"cmds"`
	Volumes   []Volume    `yaml:"volumes"`
	Sysctls   []Sysctl    `yaml:"sysctls"`
}

type InterFace struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Ipv4 string `yaml:"ipv4"`
	Ipv6 string `yaml:"ipv6"`
	Args string `yaml:"args"`
}

type Sysctl struct {
	Sysctl string `yaml:"sysctl"`
}

type Cmd struct {
	Cmd string `yaml:"cmd"`
}

type Volume struct {
	Volume string `yaml:"volume"`
}

func ParseYaml(filepath string) []Node {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	// fmt.Println("----------------------------------------------")
	// fmt.Println("                   CONFIG                     ")
	// fmt.Println("----------------------------------------------")
	// fmt.Printf("%+v\n", string(buf))

	var nodes Data

	err = yaml.Unmarshal(buf, &nodes)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("nodes: %+v\n", nodes)

	return nodes.Nodes
}
