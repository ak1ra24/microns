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
	"github.com/docker/docker/client"
	"github.com/vishvananda/netlink"
	"golang.org/x/net/context"
)

type Status struct {
	Name   string    `json:"name"`
	Status Component `json:"status"`
}

type Component struct {
	Ns        string `json:"ns"`
	Container string `json:"container"`
}

func Pull(ctx context.Context, cli *client.Client, nodes []utils.Node) {

	for _, node := range nodes {

		imageName := node.Image
		imageName = "docker.io/" + imageName
		fmt.Println(imageName)

		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
		if err != nil {
			fmt.Printf("Failed to Get Container List: %v\n", err)
			os.Exit(1)
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
				out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
				if err != nil {
					fmt.Printf("Failed to Pull Image: %v\n", err)
					os.Exit(1)
				}
				io.Copy(os.Stdout, out)

				resp, err := cli.ContainerCreate(ctx,
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
					fmt.Printf("Failed to Create Container: %v\n", err)
					os.Exit(1)
				}

				if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
					fmt.Printf("Failed to Start Container: %v\n", err)
					os.Exit(1)
				}
			}
		} else {
			out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
			if err != nil {
				fmt.Printf("Failed to Pull Image: %v\n", err)
				os.Exit(1)
			}
			io.Copy(os.Stdout, out)

			resp, err := cli.ContainerCreate(ctx,
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
				fmt.Printf("Failed to Create Container: %v\n", err)
				os.Exit(1)
			}

			if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				fmt.Printf("Failed to Start Container: %v\n", err)
				os.Exit(1)
			}
		}

	}
}

func Dockertonetns(ctx context.Context, cli *client.Client, nodename string) {

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		fmt.Printf("Failed to Get Container List: %v\n", err)
		os.Exit(1)
	}

	var pid string

	var containerName string

	for _, container := range containers {
		containerName = strings.Replace(container.Names[0], "/", "", 1)
		if nodename == containerName {
			json, err := cli.ContainerInspect(ctx, container.ID)
			if err != nil {
				fmt.Printf("Failed to Inspect Container: %v\n", err)
				os.Exit(1)
			}
			pid = strconv.Itoa(json.State.Pid)

			fmt.Printf("Image: %s, ID: %s, Name: %s, PID: %s\n", container.Image, container.ID, containerName, pid)
			if _, err := os.Stat("/proc/" + pid); err != nil {
				fmt.Printf("Not Found /proc/ %s: %v\n", pid, err)
				os.Exit(1)
			}
			fmt.Printf("/proc/" + pid + "is Exist\n")
			dockerns := fmt.Sprintf("/proc/%s/ns/net", pid)

			if _, err := os.Stat("/var/run/netns"); os.IsNotExist(err) {
				// path/to/whatever does not exist
				if err := os.MkdirAll("/var/run/netns", 0755); err != nil {
					fmt.Printf("Failed to Make Dir /var/run/netns: %v\n", err)
					os.Exit(1)
				}

			}
			netns := fmt.Sprintf("/var/run/netns/%s", containerName)

			if _, err := os.Stat(netns); os.IsNotExist(err) {
				if err := os.Symlink(dockerns, netns); err != nil {
					fmt.Printf("Failed to symlink %s -> %s: %v", dockerns, netns, err)
					os.Exit(1)
				}
			}
		}
	}
}

func SetBridge(node utils.Node, inf utils.Interface) {
	bridge := &netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: inf.PeerNode, Flags: net.FlagUp, MTU: 1500}}

	netlink.LinkAdd(bridge)

	node1 := node.Name
	name := fmt.Sprintf("%s", inf.InfName)
	peername := fmt.Sprintf("%s-%s", node.Name, inf.PeerNode)

	// get path
	path := "/var/run/netns/"
	node1path := path + node1

	pid1, err := utils.ParsePid(node1path)
	if err != nil {
		fmt.Printf("Failed to utils.ParsePid: %v\n", err)
		os.Exit(1)
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
		panic(err)
	}
	peerlink1, err := netlink.LinkByName(veth.PeerName)
	if err != nil {
		panic(err)
	}

	// set
	if err := netlink.LinkSetNsPid(link1, pid1); err != nil {
		fmt.Printf("ip link %s already exist\n", link1.Attrs().Name)
	}

	if err := netlink.LinkSetMaster(peerlink1, bridge); err != nil {
		panic(err)
	}

	vethNS1, err := ns.GetNS(node1path)
	if err != nil {
		fmt.Printf("Failed to get NS %s: %v\n", node1path, err)
		os.Exit(1)
	}
	defer vethNS1.Close()

	err = vethNS1.Do(func(_ ns.NetNS) error {
		linkns, err := netlink.LinkByName(link1.Attrs().Name)
		if err != nil {
			return err
		}

		if err = netlink.LinkSetUp(linkns); err != nil {
			return err
		}

		if inf.Ipv4 != "" {
			inf_val := strings.Split(inf.Ipv4, "/")
			ipv4addr := inf_val[0]
			mask_str := inf_val[1]
			mask, _ := strconv.Atoi(mask_str)
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
			inf_val := strings.Split(inf.Ipv6, "/")
			ipv6addr := inf_val[0]
			mask_str := inf_val[1]
			mask, _ := strconv.Atoi(mask_str)
			ip := net.IPNet{IP: net.ParseIP(ipv6addr), Mask: net.CIDRMask(mask, 128)}

			addr := &netlink.Addr{IPNet: &ip, Label: ""}
			if err = netlink.AddrAdd(linkns, addr); err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		fmt.Printf("Failed to Configure NS: %v\n", err)
	}
}

func SetLink(node utils.Node, inf utils.Interface) {
	node1 := node.Name
	name := fmt.Sprintf("%s", inf.InfName)
	peername := fmt.Sprintf("%s", inf.PeerInf)

	// get path
	path := "/var/run/netns/"
	node1path := path + node1

	pid1, err := utils.ParsePid(node1path)
	if err != nil {
		fmt.Printf("Failed to utils.ParsePid: %v\n", err)
		os.Exit(1)
	}

	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name:      name,
			Flags:     net.FlagUp,
			MTU:       1500,
			Namespace: pid1,
		},
		PeerName: peername,
	}

	if err := netlink.LinkAdd(veth); err != nil {
		fmt.Printf("ip link already added %s -- %s\n", veth.LinkAttrs.Name, veth.PeerName)
	}

	link1, err := netlink.LinkByName(veth.LinkAttrs.Name)
	if err != nil {
		fmt.Printf("Not Found Link: %v\n", err)
		os.Exit(1)
	}

	// set link
	if err := netlink.LinkSetNsPid(link1, pid1); err != nil {
		fmt.Printf("ip link %s already exist\n", link1.Attrs().Name)
	}

	vethNS1, err := ns.GetNS(node1path)
	if err != nil {
		fmt.Printf("Failed to get NS %s: %v\n", node1path, err)
		os.Exit(1)
	}
	defer vethNS1.Close()

	err = vethNS1.Do(func(_ ns.NetNS) error {
		linkns, err := netlink.LinkByName(link1.Attrs().Name)
		if err != nil {
			return err
		}

		if err = netlink.LinkSetUp(linkns); err != nil {
			return err
		}

		if inf.Ipv4 != "" {
			inf_val := strings.Split(inf.Ipv4, "/")
			ipv4addr := inf_val[0]
			mask_str := inf_val[1]
			mask, _ := strconv.Atoi(mask_str)
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
			inf_val := strings.Split(inf.Ipv6, "/")
			ipv6addr := inf_val[0]
			mask_str := inf_val[1]
			mask, _ := strconv.Atoi(mask_str)
			ip := net.IPNet{IP: net.ParseIP(ipv6addr), Mask: net.CIDRMask(mask, 128)}

			addr := &netlink.Addr{IPNet: &ip, Label: ""}
			if err = netlink.AddrAdd(linkns, addr); err != nil {
				return err
			}
		}

		return nil

	})

	if err != nil {
		fmt.Printf("Failed to Configure NS: %v\n", err)
	}
}

func LinkUp(linkname string) {
	link, err := netlink.LinkByName(linkname)
	if err != nil {
		panic(err)
	}
	if err := netlink.LinkSetUp(link); err != nil {
		panic(err)
	}
}

func SetConf(ctx context.Context, cli *client.Client, container_name string, cmd string) {

	var runcmd []string

	if strings.Contains(cmd, "vtysh") {
		for i, str := range strings.Split(cmd, " -c ") {
			if i != 0 {
				runcmd = append(runcmd, "-c")
			}
			str = strings.Replace(str, "\"", "", -1)
			runcmd = append(runcmd, str)
		}
	} else {
		runcmd = strings.Split(cmd, " ")
	}

	idreq, err := cli.ContainerExecCreate(ctx, container_name, types.ExecConfig{
		User:         "root",
		Privileged:   true,
		Tty:          true,
		Detach:       false,
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          runcmd,
	})
	//
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerExecStart(ctx, idreq.ID, types.ExecStartCheck{
		Detach: false,
		Tty:    true,
	}); err != nil {
		panic(err)
	}

}

func RemoveNs(ctx context.Context, cli *client.Client, nodename string) {
	path := "/var/run/netns/"
	nodepath := path + nodename
	if err := os.Remove(nodepath); err != nil {
		fmt.Printf("Failed to Remove %s: %v\n", nodepath, err)
		os.Exit(1)
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		fmt.Printf("Failed to Get ContainerList: %v\n", err)
		os.Exit(1)
	}

	for _, container := range containers {
		fmt.Println(container.ID)
		containerName := strings.Replace(container.Names[0], "/", "", 1)
		if nodename == containerName {
			if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
				fmt.Printf("Failed to Remove Container %s: %v\n", containerName, err)
				os.Exit(1)
			}
		}
	}
}

func RemoveBr(bridgeName string) {
	br, err := netlink.LinkByName(bridgeName)
	if err != nil {
		fmt.Printf("Not Found Link: %v\n", err)
		os.Exit(1)
	}
	if err := netlink.LinkDel(br); err != nil {
		panic(err)
	}
}

func StatusNs(ctx context.Context, cli *client.Client, nodename string) string {
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
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
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

	return jsonstring
}

// func main() {
// 	ctx := context.Background()
// 	cli, err := client.NewEnvClient()
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	nodes := utils.ParseYaml("./config.yaml.bk2")
//
// 	Pull(ctx, cli, nodes)
//
// 	for _, node := range nodes {
// 		Dockertonetns(ctx, cli, node.Name)
// 	}
// 	for _, node := range nodes {
// 		fmt.Println(node.Interface)
// 		for _, inf := range node.Interface {
// 			SetLink(node, inf)
// 		}
// 	}
//
// 	// remove container and netns
// 	// for _, node := range nodes {
// 	// 	RemoveNs(ctx, cli, node.Name)
// 	// }
// }
