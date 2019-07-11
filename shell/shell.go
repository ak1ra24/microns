package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/ak1ra24/microns/api/utils"
)

// func RunContainer(nodename, imagename string) string {
func RunContainer(node utils.Node) string {

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
		runcmd += fmt.Sprintf(" %s", node.Image)
		CheckandRuncmd = fmt.Sprintf(check_container, node.Name, node.Name, runcmd)
	} else {
		runcmd = fmt.Sprintf("docker run -td --rm --net=none --privileged --name %s --hostname %s %s", node.Name, node.Name, node.Image)
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
	if _, err := os.Stat(netns); os.IsNotExist(err) {
		symlinkcmd = fmt.Sprintf("ln -s /proc/$PID/ns/net %s", netns)
	} else {
		symlinkcmd = fmt.Sprintf("echo %s is Exist", netns)
	}

	return symlinkcmd
}

func LinkAdd(node utils.Node, inf utils.InterFace) (string, string) {
	nodelink := fmt.Sprintf("%s:%s", node.Name, inf.Name)
	node1 := strings.Split(nodelink, ":")[0]
	node2 := strings.Split(inf.Args, ":")[0]
	vethname := node1 + "_to_" + node2
	peername := node2 + "_to_" + node1

	check_link := `
if ip link show %s > /dev/null 2>&1; then
	echo %s is Exist 
elif ip netns exec %s ip link show %s > /dev/null 2>&1; then
	echo netns:%s link:%s is Exist 
else
	ip link add %s type veth peer name %s
fi
`

	CheckandAddLinkcmd := fmt.Sprintf(check_link, vethname, vethname, node1, vethname, node1, vethname, vethname, peername)

	// addlinkcmd := fmt.Sprintf("ip link add %s type veth peer name %s", vethname, peername)

	return vethname, CheckandAddLinkcmd
}

func LinkSetNs(vethname, nodename string) string {
	check_set_link := `
if ip netns exec %s ip link show %s > /dev/null 2>&1; then
	echo Already Set netns:%s link:%s is Exist 
else
	ip link set %s netns %s up
fi
`
	// setLinkNscmd := fmt.Sprintf("ip link set %s netns %s up", vethname, nodename)
	setLinkNscmd := fmt.Sprintf(check_set_link, nodename, vethname, nodename, vethname, vethname, nodename)

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

func NsDel(nodename string) string {
	delNscmd := fmt.Sprintf("ip netns delete %s", nodename)

	return delNscmd
}

func DockerDel(nodename string) string {
	delDockercmd := fmt.Sprintf("docker rm -f %s", nodename)

	return delDockercmd
}
