package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/ak1ra24/microns/utils"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/containernetworking/plugins/pkg/utils/sysctl"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/context"
)

// Status struct is Container and Netns status
type Status struct {
	Name   string    `json:"name"`
	Status Component `json:"status"`
}

// Component for microns status
type Component struct {
	Ns        string `json:"ns"`
	Container string `json:"container"`
}

// Container Operation struct
type Container struct {
	Ctx context.Context
	Cli *client.Client
}

// NewContainer func is Create Container Client
func NewContainer(ctx context.Context, cli *client.Client) *Container {
	container := &Container{
		Ctx: ctx,
		Cli: cli,
	}

	return container
}

// CreateContainerPort func is container create port binding
func (c *Container) CreateContainerPort(imageName string, hostName string, cport string, hport string, cfgfile string) error {
	imageName = "docker.io/" + imageName

	out, err := c.Cli.ImagePull(c.Ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, out)

	var port nat.Port
	var portbindings []nat.PortBinding
	var portMap nat.PortMap
	portbindings = append(portbindings, nat.PortBinding{HostPort: hport})
	fmt.Println("[Port Binding] Container Port: ", cport, "Host Port: ", hport)
	port, err = nat.NewPort("tcp", cport)
	if err != nil {
		return err
	}
	// portMap := nat.PortMap{port: portbindings}
	portMap[port] = portbindings

	containerBindFile := "/temp"

	p, _ := os.Getwd()
	if strings.HasPrefix(cfgfile, "./") {
		cfgfile = strings.Replace(cfgfile, "./", "", -1)
	}
	splitdir := strings.Split(cfgfile, "/")

	cfgfile = p + "/" + splitdir[0]

	fmt.Println("[Mount Volume] HostPath: ", cfgfile, " : ", "ContainerPath: ", containerBindFile)

	cmd := containerBindFile
	for _, cfgdir := range splitdir[1:] {
		cmd += "/" + cfgdir
	}
	fmt.Println("Exec Cmd: ", cmd)

	resp, err := c.Cli.ContainerCreate(c.Ctx,
		&container.Config{
			Image:    imageName,
			Hostname: hostName,
			Tty:      true,
			Cmd:      []string{cmd},
		},
		&container.HostConfig{
			PortBindings: portMap,
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: cfgfile,
					Target: containerBindFile,
				},
			},
		}, nil, hostName)
	if err != nil {
		return err
	}
	if err := c.Cli.ContainerStart(c.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

// Pull and Create Docker Container from config and Set Option for sysctl and volume
func (c *Container) Pull(nodes []utils.Node) error {

	for _, node := range nodes {

		imageName := node.Image
		imageName = "docker.io/" + imageName
		fmt.Println(imageName)

		containers, err := c.Cli.ContainerList(c.Ctx, types.ContainerListOptions{})
		if err != nil {
			return err
		}

		var containerNames []string
		sysctlconfs := make(map[string]string)

		for _, sysctl := range node.Sysctls {
			sysctlconf := strings.Split(sysctl.Sysctl, "=")
			sysctlconfs[sysctlconf[0]] = sysctlconf[1]
		}

		var volumes []string
		if len(node.Volumes) != 0 {
			for _, vol := range node.Volumes {
				volume := fmt.Sprintf("%s:%s", vol.HostVolume, vol.ContainerVolume)
				volumes = append(volumes, volume)
			}
		}

		if len(containers) != 0 {
			for _, conainerr := range containers {
				containerName := strings.Replace(conainerr.Names[0], "/", "", 1)
				containerNames = append(containerNames, containerName)
			}
			exist := utils.Contains(containerNames, node.Name)
			if exist {
				fmt.Printf("%s already created!\n", node.Name)
			} else {
				out, err := c.Cli.ImagePull(c.Ctx, imageName, types.ImagePullOptions{})
				if err != nil {
					return err
				}
				io.Copy(os.Stdout, out)

				resp, err := c.Cli.ContainerCreate(c.Ctx,
					&container.Config{
						Image:    imageName,
						Hostname: node.Name,
						Tty:      true,
					},
					&container.HostConfig{
						Privileged:  true,
						NetworkMode: "none",
						Sysctls:     sysctlconfs,
						Binds:       volumes,
					}, nil, node.Name)
				if err != nil {
					return err
				}

				if err := c.Cli.ContainerStart(c.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
					return err
				}
			}
		} else {
			out, err := c.Cli.ImagePull(c.Ctx, imageName, types.ImagePullOptions{})
			if err != nil {
				return err
			}
			io.Copy(os.Stdout, out)

			resp, err := c.Cli.ContainerCreate(c.Ctx,
				&container.Config{
					Image:    imageName,
					Tty:      true,
					Hostname: node.Name,
				},
				&container.HostConfig{
					Privileged:  true,
					NetworkMode: "none",
					Sysctls:     sysctlconfs,
					Binds:       volumes,
				}, nil, node.Name)
			if err != nil {
				return err
			}

			if err := c.Cli.ContainerStart(c.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				return err
			}
		}
	}
	return nil
}

// Dockertonetns is Mount Docker network namespace tp network namespace
func (c *Container) Dockertonetns(nodename string) error {

	containers, err := c.Cli.ContainerList(c.Ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}

	var pid string
	var containerName string

	for _, container := range containers {
		containerName = strings.Replace(container.Names[0], "/", "", 1)
		if nodename == containerName {
			json, err := c.Cli.ContainerInspect(c.Ctx, container.ID)
			if err != nil {
				return err
			}
			pid = strconv.Itoa(json.State.Pid)

			fmt.Printf("Image: %s, ID: %s, Name: %s, PID: %s\n", container.Image, container.ID, containerName, pid)
			if _, err := os.Stat("/proc/" + pid); err != nil {
				return err
			}
			fmt.Printf("/proc/" + pid + "is Exist\n")
			dockerns := fmt.Sprintf("/proc/%s/ns/net", pid)

			if _, err := os.Stat("/var/run/netns"); os.IsNotExist(err) {
				// path/to/whatever does not exist
				if err := os.MkdirAll("/var/run/netns", 0755); err != nil {
					return err
				}

			}
			netns := fmt.Sprintf("/var/run/netns/%s", containerName)

			if _, err := os.Stat(netns); os.IsNotExist(err) {
				if err := os.Symlink(dockerns, netns); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// SetBridge is Create Linux Bridge, and Link Bridge and Veth, and Set ip address to interface
func SetBridge(node utils.Node, inf utils.Interface) error {
	bridge := &netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: inf.PeerNode, Flags: net.FlagUp, MTU: 1500}}

	netlink.LinkAdd(bridge)

	node1 := node.Name
	name := fmt.Sprintf("%s-%s", node.Name, inf.InfName)
	peername := fmt.Sprintf("%s-%s", node.Name, inf.PeerNode)

	// get path
	path := "/var/run/netns/"
	node1path := path + node1

	pid1, err := utils.ParsePid(node1path)
	if err != nil {
		return err
	}
	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name:  name,
			Flags: net.FlagUp,
			MTU:   1500,
		},
		PeerName: peername,
	}

	if err := netlink.LinkAdd(veth); err != nil {
		fmt.Printf("ip link already added %s -- %s\n", veth.LinkAttrs.Name, veth.PeerName)
	}
	fmt.Println("vethname: ", veth.LinkAttrs.Name)
	link1, err := netlink.LinkByName(veth.LinkAttrs.Name)
	if err != nil {
		return err
	}
	peerlink1, err := netlink.LinkByName(veth.PeerName)
	if err != nil {
		return err
	}

	// set
	if err := netlink.LinkSetNsPid(link1, pid1); err != nil {
		fmt.Printf("ip link %s already exist\n", link1.Attrs().Name)
	}

	if err := netlink.LinkSetMaster(peerlink1, bridge); err != nil {
		return err
	}

	vethNS1, err := ns.GetNS(node1path)
	if err != nil {
		return err
	}
	defer vethNS1.Close()

	err = vethNS1.Do(func(_ ns.NetNS) error {
		linkns, err := netlink.LinkByName(link1.Attrs().Name)
		if err != nil {
			return err
		}

		if err = netlink.LinkSetName(link1, inf.InfName); err != nil {
			return err
		}

		if err = netlink.LinkSetUp(linkns); err != nil {
			return err
		}

		if inf.Ipv4 != "" {
			infVal := strings.Split(inf.Ipv4, "/")
			ipv4addr := infVal[0]
			maskStr := infVal[1]
			mask, _ := strconv.Atoi(maskStr)
			ip := net.IPNet{IP: net.ParseIP(ipv4addr), Mask: net.CIDRMask(mask, 32)}

			addr := &netlink.Addr{IPNet: &ip, Label: ""}
			if err = netlink.AddrAdd(linkns, addr); err != nil {
				return err
			}
		}

		if inf.Ipv6 != "" {
			ipv6SysctlAll := fmt.Sprint("net.ipv6.conf.all.disable_ipv6")
			ipv6SysctlDefault := fmt.Sprint("net.ipv6.conf.default.disable_ipv6")
			if _, err := sysctl.Sysctl(ipv6SysctlAll, "0"); err != nil {
				return fmt.Errorf("failed to set ipv6.disable to 0 : %v", err)
			}
			if _, err := sysctl.Sysctl(ipv6SysctlDefault, "0"); err != nil {
				return fmt.Errorf("failed to set ipv6.disable to 0 : %v", err)
			}
			infVal := strings.Split(inf.Ipv6, "/")
			ipv6addr := infVal[0]
			maskStr := infVal[1]
			mask, _ := strconv.Atoi(maskStr)
			ip := net.IPNet{IP: net.ParseIP(ipv6addr), Mask: net.CIDRMask(mask, 128)}

			addr := &netlink.Addr{IPNet: &ip, Label: ""}
			if err = netlink.AddrAdd(linkns, addr); err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

// SetLink is Link Node interface and Other Node interface, and Set ip address to interface
func SetLink(node utils.Node, inf utils.Interface) error {
	node1 := node.Name
	name := fmt.Sprintf("%s-%s", node.Name, inf.InfName)
	peername := fmt.Sprintf("%s-%s", inf.PeerNode, inf.PeerInf)

	// get path
	path := "/var/run/netns/"
	node1path := path + node1

	pid1, err := utils.ParsePid(node1path)
	if err != nil {
		return err
	}

	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name:  name,
			Flags: net.FlagUp,
			MTU:   1500,
		},
		PeerName: peername,
	}

	if err := netlink.LinkAdd(veth); err != nil {
		fmt.Printf("ip link already added %s -- %s\n", veth.LinkAttrs.Name, veth.PeerName)
	}

	link1, err := netlink.LinkByName(veth.LinkAttrs.Name)
	if err != nil {
		return err
	}

	// set link
	if err := netlink.LinkSetNsPid(link1, pid1); err != nil {
		fmt.Printf("ip link %s already exist\n", link1.Attrs().Name)
	}

	vethNS1, err := ns.GetNS(node1path)
	if err != nil {
		return err
	}
	defer vethNS1.Close()

	err = vethNS1.Do(func(_ ns.NetNS) error {
		linkns, err := netlink.LinkByName(link1.Attrs().Name)
		if err != nil {
			return err
		}

		if err = netlink.LinkSetName(link1, inf.InfName); err != nil {
			return err
		}

		if err = netlink.LinkSetUp(linkns); err != nil {
			return err
		}

		if inf.Ipv4 != "" {
			infVal := strings.Split(inf.Ipv4, "/")
			ipv4addr := infVal[0]
			maskStr := infVal[1]
			mask, _ := strconv.Atoi(maskStr)
			ip := net.IPNet{IP: net.ParseIP(ipv4addr), Mask: net.CIDRMask(mask, 32)}

			addr := &netlink.Addr{IPNet: &ip, Label: ""}
			if err = netlink.AddrAdd(linkns, addr); err != nil {
				return err
			}
		}

		if inf.Ipv6 != "" {
			ipv6SysctlAll := fmt.Sprint("net.ipv6.conf.all.disable_ipv6")
			ipv6SysctlDefault := fmt.Sprint("net.ipv6.conf.default.disable_ipv6")
			if _, err := sysctl.Sysctl(ipv6SysctlAll, "0"); err != nil {
				return fmt.Errorf("failed to set ipv6.disable to 0 : %v", err)
			}
			if _, err := sysctl.Sysctl(ipv6SysctlDefault, "0"); err != nil {
				return fmt.Errorf("failed to set ipv6.disable to 0 : %v", err)
			}
			infVal := strings.Split(inf.Ipv6, "/")
			ipv6addr := infVal[0]
			maskStr := infVal[1]
			mask, _ := strconv.Atoi(maskStr)
			ip := net.IPNet{IP: net.ParseIP(ipv6addr), Mask: net.CIDRMask(mask, 128)}

			addr := &netlink.Addr{IPNet: &ip, Label: ""}
			if err = netlink.AddrAdd(linkns, addr); err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		return err
	}

	return nil
}

// LinkUp is Link up to interface
func LinkUp(linkname string) error {
	link, err := netlink.LinkByName(linkname)
	if err != nil {
		return err
	}
	if err := netlink.LinkSetUp(link); err != nil {
		return err
	}

	return nil
}

// SetConf is Set config from config
func (c *Container) SetConf(containerName string, cmd string) error {

	// convert command for docekr exec
	splitCmds := strings.Split(cmd, " ")
	var runcmd string
	var runcmds []string
	for _, splitCmd := range splitCmds {
		splitCmd = strings.TrimSpace(splitCmd)
		if strings.HasPrefix(splitCmd, "\"") && strings.HasSuffix(splitCmd, "\"") {
			runcmd = splitCmd
			runcmds = append(runcmds, runcmd)
			runcmd = ""
		} else if strings.HasPrefix(splitCmd, "\"") {
			runcmd = strings.TrimLeft(splitCmd, "\"")
		} else if strings.HasSuffix(splitCmd, "\"") {
			runcmd += " " + strings.TrimRight(splitCmd, "\"")
			runcmds = append(runcmds, runcmd)
			runcmd = ""
		} else if len(runcmd) > 0 {
			runcmd += " " + splitCmd
		} else {
			runcmd = splitCmd
			runcmds = append(runcmds, runcmd)
			runcmd = ""
		}
	}

	fmt.Println(runcmds)

	idreq, err := c.Cli.ContainerExecCreate(c.Ctx, containerName, types.ExecConfig{
		User:         "root",
		Privileged:   true,
		Tty:          true,
		Detach:       false,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          runcmds,
	})

	if err != nil {
		return err
	}

	if err := c.Cli.ContainerExecStart(c.Ctx, idreq.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    true,
	}); err != nil {
		return err
	}

	return nil
}

// RemoveNs is Remove Docker Container, and Remove Network Namespace
func (c *Container) RemoveNs(nodename string) error {
	path := "/var/run/netns/"
	nodepath := path + nodename
	if err := os.Remove(nodepath); err != nil {
		fmt.Println("Not Exist: ", nodename)
	}

	return nil
}

// RemoveContainer is Remove Docker Container
func (c *Container) RemoveContainer(nodename string) error {

	containers, err := c.Cli.ContainerList(c.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}

	for _, container := range containers {
		containerName := strings.Replace(container.Names[0], "/", "", 1)
		if nodename == containerName {
			fmt.Println("Remove ContainerName: ", containerName)
			if err := c.Cli.ContainerRemove(c.Ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				fmt.Println("Not Exist Container Name: ", nodename)
			}
		}
	}

	return nil

}

// RemoveBr is Remove Linux Bridge
func RemoveBr(bridgeName string) error {
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		fmt.Println("Not Exist Bridge: ", bridgeName)
		return nil
	}
	if err := netlink.LinkDel(br); err != nil {
		return err
	}

	return nil
}

// StatusNs is Status of containers and ns made with microns
func (c *Container) StatusNs(nodename string) (string, error) {
	path := "/var/run/netns/"
	nodepath := path + nodename
	var status Status
	if _, err := os.Stat(nodepath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		// fmt.Printf("NS:\t\t%s\t\tNot Found\n", nodename)
		status.Name = nodename
		status.Status.Ns = "Not Found"
	} else {
		// fmt.Printf("NS:\t\t%s\t\tFound\n", nodename)
		status.Name = nodename
		status.Status.Ns = "Found"
	}
	containers, err := c.Cli.ContainerList(c.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return "", err
	}
	if len(containers) == 0 {
		// fmt.Printf("Container:\t%s\t\tNot Found\n", nodename)
		status.Status.Container = "Not Found"
	} else {
		for _, container := range containers {
			containerName := strings.Replace(container.Names[0], "/", "", 1)
			if nodename == containerName {
				// fmt.Printf("Container:\t%s\t\t%s\n", containerName, container.State)
				status.Status.Container = container.State
			}
		}
	}

	jsonbyte, _ := json.Marshal(status)
	jsonstring := string(jsonbyte)

	return jsonstring, nil
}

// StartContainer is start docker container
func (c *Container) StartContainer(nodename string) error {

	containers, err := c.Cli.ContainerList(c.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	for _, container := range containers {
		containerName := strings.Replace(container.Names[0], "/", "", 1)
		if nodename == containerName {
			fmt.Println("start containerName: ", containerName)
			if err := c.Cli.ContainerStart(c.Ctx, container.ID, types.ContainerStartOptions{}); err != nil {
				return err
			}
		}
	}

	return nil

}

// StopContainer is stop docker container
func (c *Container) StopContainer(nodename string) error {

	containers, err := c.Cli.ContainerList(c.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return err
	}
	for _, container := range containers {
		containerName := strings.Replace(container.Names[0], "/", "", 1)
		if nodename == containerName {
			fmt.Println("stop containerName: ", containerName)
			if err := c.Cli.ContainerStop(c.Ctx, container.ID, nil); err != nil {
				return err
			}
		}
	}

	return nil

}
