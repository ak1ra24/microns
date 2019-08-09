package utils

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Tn struct {
	Nodes []struct {
		Name       string `yaml:"name"`
		Image      string `yaml:"image"`
		Interfaces []struct {
			Name string `yaml:"name"`
			Type string `yaml:"type"`
			Args string `yaml:"args"`
		} `yaml:"interfaces"`
	} `yaml:"nodes"`
	Switches []struct {
		Name       string        `yaml:"name"`
		Interfaces []TnInterface `yaml:"interfaces"`
	} `yaml:"switches"`
	NodeConfigs []struct {
		Name string `yaml:"name"`
		Cmds []struct {
			Cmd string `yaml:"cmd"`
		} `yaml:"cmds"`
	} `yaml:"node_configs"`
	Test []struct {
		Name string `yaml:"name"`
		Cmds []struct {
			Cmd string `yaml:"cmd"`
		} `yaml:"cmds"`
	} `yaml:"test"`
}

type TnInterface struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Args string `yaml:"args"`
}

type Microns struct {
	Nodes       []Node       `yaml:"nodes"`
	Switches    []Switch     `yaml:"switches"`
	NodeConfigs []Nodeconfig `yaml:"node_config"`
	Test        []TestCmd    `yaml:"test"`
}

type Nodes struct {
	Nodes []Node `yaml:"nodes"`
}

type Node struct {
	Name      string      `yaml:"name"`
	Image     string      `yaml:"image"`
	Interface []Interface `yaml:"interfaces"`
	Volumes   []Volume    `yaml:"volumes"`
	Sysctls   []Sysctl    `yaml:"sysctls"`
}

type Config struct {
	Config []Nodeconfig `yaml:"node_config"`
}

type Nodeconfig struct {
	Name string `yaml:"name"`
	Cmds []Cmd  `yaml:"cmds"`
}

type Test struct {
	Testcmds []TestCmd `yaml:"test"`
}

type TestCmd struct {
	Name string `yaml:"name"`
	Cmds []Cmd  `yaml:"cmds"`
}

type Interface struct {
	InfName  string `yaml:"inf"`
	Type     string `yaml:"type"`
	Ipv4     string `yaml:"ipv4"`
	Ipv6     string `yaml:"ipv6"`
	PeerNode string `yaml:"peernode"`
	PeerInf  string `yaml:"peerinf"`
}

type Sysctl struct {
	Sysctl string `yaml:"sysctl"`
}

type Cmd struct {
	Cmd string `yaml:"cmd"`
}

type Volume struct {
	HostVolume      string `yaml:"hostvolume"`
	ContainerVolume string `yaml:"containervolume"`
}

type Switches struct {
	Switches []Switch `yaml:"switches"`
}

type Switch struct {
	Name       string      `yaml:"name"`
	Interfaces []Interface `yaml:"interfaces"`
}

func ParseNodes(filepath string) []Node {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var nodes Nodes

	err = yaml.Unmarshal(buf, &nodes)
	if err != nil {
		panic(err)
	}

	return nodes.Nodes
}

func ParseSwitch(filepath string) []Switch {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var switches Switches

	err = yaml.Unmarshal(buf, &switches)
	if err != nil {
		panic(err)
	}

	return switches.Switches
}

func ParseConfig(filepath string) []Nodeconfig {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var nodeconfigs Config

	err = yaml.Unmarshal(buf, &nodeconfigs)
	if err != nil {
		panic(err)
	}

	return nodeconfigs.Config
}

func ParseTest(filepath string) []TestCmd {
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var tests Test
	err = yaml.Unmarshal(buf, &tests)
	if err != nil {
		panic(err)
	}

	return tests.Testcmds
}

func CreateCfgFile(filename string) string {
	conf := `
nodes:
  - name: 
    image: 
    interfaces:
        - inf: 
          type: 
          ipv4: 
          ipv6:
          peernode: 
          peerinf: 
    volumes:
        - hostvolume: 
          containervolume: 
    sysctls:
        - sysctl:
switches:
	- name: 
	interfaces:
		- name: 
		  type: 
		  peernode: 

node_config:
  - name: 
    cmds:
        - cmd: 
test:
  - cmds:
        - cmd: 
`
	fp, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	fp.WriteString(conf)

	return conf
}
