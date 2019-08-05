package api

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/ak1ra24/microns/utils"
	yaml "gopkg.in/yaml.v2"
)

func Convert(tncfg, mncfg string) error {
	var tn utils.Tn
	mn := new(utils.Microns)

	buf, err := ioutil.ReadFile(tncfg)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, &tn)
	if err != nil {
		return err
	}

	for _, n := range tn.Nodes {
		var intfaces []utils.InterFace
		for _, inf := range n.Interfaces {
			if strings.Contains(inf.Args, "#") {
				peer := strings.Split(inf.Args, "#")
				peernode := peer[0]
				peerinf := peer[1]
				var intface utils.InterFace = utils.InterFace{
					InfName:  inf.Name,
					Type:     inf.Type,
					PeerNode: peernode,
					PeerInf:  peerinf,
				}
				intfaces = append(intfaces, intface)
			} else {
				var intface utils.InterFace = utils.InterFace{
					InfName:  inf.Name,
					Type:     inf.Type,
					PeerNode: inf.Args,
				}
				intfaces = append(intfaces, intface)
			}
		}
		var m utils.NodeInfo = utils.NodeInfo{
			Name:      n.Name,
			Image:     n.Image,
			Interface: intfaces,
		}

		mn.Nodes = append(mn.Nodes, m)
	}

	for _, c := range tn.NodeConfigs {
		var nodecmds []utils.Cmd
		for _, cmd := range c.Cmds {
			nodecmds = append(nodecmds, cmd)
		}
		var nodeconfig utils.Nodeconfig = utils.Nodeconfig{
			Name: c.Name,
			Cmds: nodecmds,
		}
		mn.NodeConfigs = append(mn.NodeConfigs, nodeconfig)
	}

	for _, s := range tn.Switches {
		var sinfs []utils.InterFace
		for _, tninf := range s.Interfaces {
			sinf := utils.InterFace{
				InfName:  tninf.Name,
				Type:     tninf.Type,
				PeerNode: tninf.Args,
			}
			sinfs = append(sinfs, sinf)
		}
		var sw utils.Switch = utils.Switch{
			Name:       s.Name,
			Interfaces: sinfs,
		}
		mn.Switches = append(mn.Switches, sw)
	}

	for _, t := range tn.Test {
		var testcmds []utils.Cmd
		for _, testcmd := range t.Cmds {
			testcmds = append(testcmds, testcmd)
		}
		var testCmd utils.TestCmd = utils.TestCmd{
			Name: t.Name,
			Cmds: testcmds,
		}
		mn.Test = append(mn.Test, testCmd)
	}

	mncfg_output, err := yaml.Marshal(&mn)
	if err != nil {
		return err
	}

	file, err := os.Create(mncfg)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(mncfg_output)

	return nil
}

// func main() {
//
// 	flag.Parse()
// 	if flag.NArg() < 2 {
// 		fmt.Println("go run main.go <tn configfile name> <converted file name>")
// 		os.Exit(1)
// 	}
//
// 	tncfgfile := flag.Arg(0)
// 	mncfgfile := flag.Arg(1)
// 	Convert(tncfgfile, mncfgfile)
// }
