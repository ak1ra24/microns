package shell

import (
	"fmt"
	"os"

	"github.com/ak1ra24/microns/utils"
)

// func RunContainer(nodename, imagename string) string {
func RunContainer(node utils.NodeInfo) string {

	check_container := `
if docker container ls | grep %s > /dev/null 2>&1; then
	echo Container:%s is Exist
else
	%s
fi
`
	var runcmd string
	var CheckandRuncmd string

	if len(node.Sysctls) != 0 {
		runcmd = fmt.Sprintf("docker run -td --rm --net=none --privileged --name %s --hostname %s", node.Name, node.Name)
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
		CheckandRuncmd = fmt.Sprintf(check_container, node.Name, node.Name, runcmd)
	} else {
		runcmd = fmt.Sprintf("docker run -td --rm --net=none --privileged --name %s --hostname %s", node.Name, node.Name)
		if len(node.Volumes) != 0 {
			for _, volume := range node.Volumes {
				volumeconf := fmt.Sprintf(" -v %s:%s", volume.HostVolume, volume.ContainerVolume)
				runcmd += volumeconf
			}
		}
		runcmd += fmt.Sprintf(" %s", node.Image)
		CheckandRuncmd = fmt.Sprintf(check_container, node.Name, node.Name, runcmd)
	}

	return CheckandRuncmd
}

func GetContainerPid(nodename string) string {
	getpidcmd := fmt.Sprintf("PID=`docker inspect %s --format '{{.State.Pid}}'`", nodename)

	return getpidcmd
}

func SymlinkNstoContainer(nodename string) string {
	netns := fmt.Sprintf("/var/run/netns/%s", nodename)

	if _, err := os.Stat("/var/run/netns"); os.IsNotExist(err) {
		// path/to/whatever does not exist
		if err := os.MkdirAll("/var/run/netns", 0755); err != nil {
			fmt.Printf("Failed to Make Dir /var/run/netns: %v\n", err)
			os.Exit(1)
		}
	}

	var symlinkcmd string
	symlinkcmd = fmt.Sprintf("ln -s /proc/$PID/ns/net %s", netns)
	// if _, err := os.Stat(netns); os.IsNotExist(err) {
	// 	symlinkcmd = fmt.Sprintf("ln -s /proc/$PID/ns/net %s", netns)
	// } else {
	// 	symlinkcmd = fmt.Sprintf("echo %s is Exist", netns)
	// }

	return symlinkcmd
}

func AddBr(bridge utils.Switch) string {
	br := fmt.Sprintf("ip link add %s type bridge", bridge.Name)

	return br
}

func LinkAdd(node utils.NodeInfo, inf utils.InterFace) (string, string) {
	// node1 := inf.InfName
	// node2 := inf.PeerInf
	// vethname := node1 + "_to_" + node2
	// peername := node2 + "_to_" + node1
	// vethname := fmt.Sprintf("%s-%s", node.Name, inf.InfName)
	// peername := fmt.Sprintf("%s-%s", inf.PeerNode, inf.PeerInf)
	vethname := fmt.Sprintf("%s", inf.InfName)
	peername := fmt.Sprintf("%s", inf.PeerInf)

	// 	check_link := `
	// if ip link show %s > /dev/null 2>&1; then
	// 	echo %s is Exist
	// elif ip netns exec %s ip link show %s > /dev/null 2>&1; then
	// 	echo netns:%s link:%s is Exist
	// else
	// 	ip link add %s netns %s type veth peer name %s netns %s
	// fi
	// `
	// CheckandAddLinkcmd := fmt.Sprintf(check_link, vethname, vethname, node.Name, vethname, node.Name, vethname, vethname, node.Name, peername, inf.PeerNode)
	check_link := "ip link add %s netns %s type veth peer name %s netns %s"
	CheckandAddLinkcmd := fmt.Sprintf(check_link, vethname, node.Name, peername, inf.PeerNode)

	// addlinkcmd := fmt.Sprintf("ip link add %s type veth peer name %s", vethname, peername)

	return vethname, CheckandAddLinkcmd
}

func LinkAddBr(bridges []utils.Switch, node utils.NodeInfo, inf utils.InterFace) ([]string, string, string, string) {
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

func LinkUpBridge(brname string) string {
	check_link_up_bridge := fmt.Sprintf("ip link set %s up", brname)

	return check_link_up_bridge
}

func LinkUpBrLink(brlinkname string) string {
	check_link_up_brlink := fmt.Sprintf("ip link set %s up", brlinkname)

	return check_link_up_brlink
}

func LinkSetBridge(brlinkname, bridgename string) string {
	check_set_br := fmt.Sprintf("ip link set dev %s master %s", brlinkname, bridgename)

	return check_set_br
}

func LinkSetNs(vethname, nodename string) string {
	// 	check_set_link := `
	// if ip netns exec %s ip link show %s > /dev/null 2>&1; then
	// 	echo Already Set netns:%s link:%s is Exist
	// else
	// 	ip link set %s netns %s up
	// fi
	// `
	// setLinkNscmd := fmt.Sprintf(check_set_link, nodename, vethname, nodename, vethname, vethname, nodename)
	check_set_link := "ip netns exec %s ip link set %s up"
	// setLinkNscmd := fmt.Sprintf("ip link set %s netns %s up", vethname, nodename)
	setLinkNscmd := fmt.Sprintf(check_set_link, nodename, vethname)

	return setLinkNscmd
}

func AddrAddv4(nodename, vethname string, inf utils.InterFace) string {
	// 	check_addr := `
	// if ip netns exec %s ip addr show dev %s > /dev/null 2>&1; then
	// 	echo Already Add Address
	// else
	// 	ip netns exec %s ip addr add %s dev %s
	// fi
	// `

	addAddrcmd := fmt.Sprintf("ip netns exec %s ip addr add %s dev %s", nodename, inf.Ipv4, vethname)
	// addAddrcmd := fmt.Sprintf(check_addr, nodename, vethname, nodename, inf.Ipv4, vethname)

	return addAddrcmd
}

func AddrAddv6(nodename, vethname string, inf utils.InterFace) string {
	// 	check_addr := `
	// if ip netns exec %s ip addr show dev %s > /dev/null 2>&1; then
	// 	echo Already Add Address
	// else
	// 	ip netns exec %s ip -6 addr add %s dev %s
	// fi
	// `

	addAddrcmd := fmt.Sprintf("ip netns exec %s ip -6 addr add %s dev %s", nodename, inf.Ipv6, vethname)
	// addAddrcmd := fmt.Sprintf(check_addr, nodename, vethname, nodename, inf.Ipv6, vethname)

	return addAddrcmd
}

func RunCmd(config utils.Nodeconfig) []string {

	var runcmds []string
	for _, cmd := range config.Cmds {
		runcmd := fmt.Sprintf("docker exec %s %s", config.Name, cmd.Cmd)
		runcmds = append(runcmds, runcmd)
	}

	return runcmds
}

func NsDel(nodename string) string {
	delNscmd := fmt.Sprintf("ip netns delete %s", nodename)

	return delNscmd
}

func DockerDel(nodename string) string {
	delDockercmd := fmt.Sprintf("docker rm -f %s", nodename)

	return delDockercmd
}

func BridgeDel(brname string) string {
	delBrCmd := fmt.Sprintf("ip link delete %s", brname)

	return delBrCmd
}

func RunTestCmd(testcmds utils.TestCmd) []string {

	var runtestcmds []string
	echoname := fmt.Sprintf("echo '%s'", testcmds.Name)
	runtestcmds = append(runtestcmds, echoname)
	for _, testcmd := range testcmds.Cmds {
		runtestcmds = append(runtestcmds, testcmd.Cmd)
	}

	return runtestcmds
}
