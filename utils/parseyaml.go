package utils

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// Tn struct is Tinet Config
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

// TnInterface struct is Tinet Interface config struct
type TnInterface struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Args string `yaml:"args"`
}

// Microns struct is Microns config struct
type Microns struct {
	Nodes       []Node       `yaml:"nodes"`
	Switches    []Switch     `yaml:"switches"`
	NodeConfigs []Nodeconfig `yaml:"node_config"`
	Test        []TestCmd    `yaml:"test"`
}

// Nodes struct is Microns Nodes
type Nodes struct {
	Nodes []Node `yaml:"nodes"`
}

// Node struct is Node config
type Node struct {
	Name      string      `yaml:"name"`
	Image     string      `yaml:"image"`
	Interface []Interface `yaml:"interfaces"`
	Volumes   []Volume    `yaml:"volumes"`
	Sysctls   []Sysctl    `yaml:"sysctls"`
}

// Config struct is many NodeConfigs
type Config struct {
	Config []Nodeconfig `yaml:"node_config"`
}

// Nodeconfig struct is Set NodeName and config Command
type Nodeconfig struct {
	Name string `yaml:"name"`
	Cmds []Cmd  `yaml:"cmds"`
}

// Test struct is many testcmds
type Test struct {
	Testcmds []TestCmd `yaml:"test"`
}

// TestCmd struct is Set NodeName and test command
type TestCmd struct {
	Name string `yaml:"name"`
	Cmds []Cmd  `yaml:"cmds"`
}

// Interface struct is Interface Config
type Interface struct {
	InfName  string `yaml:"inf"`
	Type     string `yaml:"type"`
	Ipv4     string `yaml:"ipv4"`
	Ipv6     string `yaml:"ipv6"`
	PeerNode string `yaml:"peernode"`
	PeerInf  string `yaml:"peerinf"`
}

// Sysctl struct is configure sysctl for docker container
type Sysctl struct {
	Sysctl string `yaml:"sysctl"`
}

// Cmd struct is command
type Cmd struct {
	Cmd string `yaml:"cmd"`
}

// Volume struct is Mount to Docker Settings
type Volume struct {
	HostVolume      string `yaml:"hostvolume"`
	ContainerVolume string `yaml:"containervolume"`
}

// Switches struct is Bridge Settings
type Switches struct {
	Switches []Switch `yaml:"switches"`
}

// Switch struct is BridgeName and Interface settings
type Switch struct {
	Name       string      `yaml:"name"`
	Interfaces []Interface `yaml:"interfaces"`
}

// ParseNodes func is parse nodes from microns config yaml file
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

// ParseSwitch func is parse switches from microns config yaml file
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

// ParseConfig func is parse node config from microns config yaml file
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

// ParseTest func is parse testcmds from microns config yaml file
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

// CreateCfgFile func is Create Template config yaml file
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
