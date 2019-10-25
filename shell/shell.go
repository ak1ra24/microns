package shell

import (
	"fmt"
	"os"

	"github.com/ak1ra24/microns/utils"
)

// RunContainer func is Output docker run command
func RunContainer(node utils.Node) string {

	checkContainer := `
if docker container ls | grep %s > /dev/null 2>&1; then
	echo Container:%s is Exist
else
	%s
fi
`
	var runcmd string
	var CheckandRuncmd string

	if len(node.Sysctls) != 0 {
		runcmd = fmt.Sprintf("docker run -td --net=none --privileged --name %s --hostname %s", node.Name, node.Name)
		for _, sysctl := range node.Sysctls {
			sysctlconf := fmt.Sprintf(" --sysctl %s", sysctl.Sysctl)
			runcmd += sysctlconf
		}
		if len(node.Volumes) != 0 {
			for _, volume := range node.Volumes {
				volumeconf := fmt.Sprintf(" -v %s:%s", volume.HostVolume, volume.ContainerVolume)
				runcmd += volumeconf
			}
		}
		runcmd += fmt.Sprintf(" %s", node.Image)
		CheckandRuncmd = fmt.Sprintf(checkContainer, node.Name, node.Name, runcmd)
	} else {
		runcmd = fmt.Sprintf("docker run -td --net=none --privileged --name %s --hostname %s", node.Name, node.Name)
		if len(node.Volumes) != 0 {
			for _, volume := range node.Volumes {
				volumeconf := fmt.Sprintf(" -v %s:%s", volume.HostVolume, volume.ContainerVolume)
				runcmd += volumeconf
			}
		}
		runcmd += fmt.Sprintf(" %s", node.Image)
		CheckandRuncmd = fmt.Sprintf(checkContainer, node.Name, node.Name, runcmd)
	}

	return CheckandRuncmd
}

// GetContainerPid func is Output get Docker PID Command
func GetContainerPid(nodename string) string {
	getpidcmd := fmt.Sprintf("PID=`docker inspect %s --format '{{.State.Pid}}'`", nodename)

	return getpidcmd
}

// SymlinkNstoContainer func is Output mount Docker network namespace to network namespace command
func SymlinkNstoContainer(nodename string) string {
	netns := fmt.Sprintf("/var/run/netns/%s", nodename)

	if _, err := os.Stat("/var/run/netns"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		if err := os.MkdirAll("/var/run/netns", 0755); err != nil {
			fmt.Printf("Failed to Make Dir /var/run/netns: %v\n", err)
			os.Exit(1)
		}
	}

	symlinkcmd := fmt.Sprintf("ln -s /proc/$PID/ns/net %s", netns)

	return symlinkcmd
}

// AddBr func is Output add linux bridge command
func AddBr(bridge utils.Switch) string {
	br := fmt.Sprintf("ip link add %s type bridge", bridge.Name)

	return br
}

// LinkAdd func is Output connect node interface and other node interface
func LinkAdd(node utils.Node, inf utils.Interface) (string, string) {
	vethname := fmt.Sprintf("%s", inf.InfName)
	peername := fmt.Sprintf("%s", inf.PeerInf)

	checkLink := "ip link add %s netns %s type veth peer name %s netns %s"
	CheckandAddLinkcmd := fmt.Sprintf(checkLink, vethname, node.Name, peername, inf.PeerNode)

	return vethname, CheckandAddLinkcmd
}

// LinkAddBr func id Output connect ns and linux bridge command
func LinkAddBr(bridges []utils.Switch, node utils.Node, inf utils.Interface) ([]string, string, string, string) {
	var checklinks []string
	var brlinkname string
	var brname string

	vethname := fmt.Sprintf("%s", inf.InfName)
	for _, br := range bridges {
		for _, intface := range br.Interfaces {
			if node.Name == intface.PeerNode {
				checklink := fmt.Sprintf("ip link add name %s netns %s type veth peer name %s-%s", vethname, node.Name, br.Name, intface.PeerNode)
				brlinkname = fmt.Sprintf("%s-%s", br.Name, intface.PeerNode)
				brname = fmt.Sprintf("%s", br.Name)
				checklinks = append(checklinks, checklink)
			}
		}
	}
	return checklinks, brlinkname, brname, vethname
}

// LinkUpBridge func is Output link up linux bridge command
func LinkUpBridge(brname string) string {

	checkLinkUpBridge := fmt.Sprintf("ip link set %s up", brname)

	return checkLinkUpBridge
}

// LinkUpBrLink func is Output link up linux bridge link command
func LinkUpBrLink(brlinkname string) string {

	checkLinkUpBrlink := fmt.Sprintf("ip link set %s up", brlinkname)

	return checkLinkUpBrlink
}

// LinkSetBridge func is Output link set linux bridge link to linux bridge command
func LinkSetBridge(brlinkname, bridgename string) string {

	checkSetBr := fmt.Sprintf("ip link set dev %s master %s", brlinkname, bridgename)

	return checkSetBr
}

// LinkSetNs func is Output link set link to ns command
func LinkSetNs(vethname, nodename string) string {

	setLinkNscmd := fmt.Sprintf("ip netns exec %s ip link set %s up", nodename, vethname)

	return setLinkNscmd
}

// AddrAddv4 func is Output set ip address v4 to ns
func AddrAddv4(nodename, vethname string, inf utils.Interface) string {

	addAddrcmd := fmt.Sprintf("ip netns exec %s ip addr add %s dev %s", nodename, inf.Ipv4, vethname)

	return addAddrcmd
}

// AddrAddv6 func is Output set ip address v6 to ns
func AddrAddv6(nodename, vethname string, inf utils.Interface) string {

	addAddrcmd := fmt.Sprintf("ip netns exec %s ip -6 addr add %s dev %s", nodename, inf.Ipv6, vethname)

	return addAddrcmd
}

// RunCmd func is Output configure cmds to docker container
func RunCmd(config utils.Nodeconfig) []string {

	var runcmds []string
	for _, cmd := range config.Cmds {
		runcmd := fmt.Sprintf("docker exec %s %s", config.Name, cmd.Cmd)
		runcmds = append(runcmds, runcmd)
	}

	return runcmds
}

// NsDel func is Output delete to ns command
func NsDel(nodename string) string {
	delNscmd := fmt.Sprintf("ip netns delete %s", nodename)

	return delNscmd
}

// DockerDel func is Output delete to docker container command
func DockerDel(nodename string) string {
	delDockercmd := fmt.Sprintf("docker rm -f %s", nodename)

	return delDockercmd
}

// BridgeDel func is Output delete to linux bridge command
func BridgeDel(brname string) string {
	delBrCmd := fmt.Sprintf("ip link delete %s", brname)

	return delBrCmd
}

// DockerStart func is Output start docker
func DockerStart(nodename string) string {
	startDockercmd := fmt.Sprintf("docker start %s", nodename)

	return startDockercmd
}

// DockerStop func is Output stop docker
func DockerStop(nodename string) string {
	stopDockercmd := fmt.Sprintf("docker stop %s", nodename)

	return stopDockercmd
}

// RunTestCmd func is Output testcmd
func RunTestCmd(testcmds utils.TestCmd) []string {

	var runtestcmds []string
	echoname := fmt.Sprintf("echo '%s'", testcmds.Name)
	runtestcmds = append(runtestcmds, echoname)
	for _, testcmd := range testcmds.Cmds {
		runtestcmds = append(runtestcmds, testcmd.Cmd)
	}

	return runtestcmds
}
