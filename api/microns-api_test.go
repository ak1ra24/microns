package api

import (
	"context"
	"testing"

	"github.com/ak1ra24/microns/utils"
	"github.com/docker/docker/client"
)

func TestNewContainer(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	NewContainer(ctx, cli)
}

func TestPull(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	c := NewContainer(ctx, cli)
	inf01 := utils.Interface{InfName: "net0", Type: "direct", Ipv4: "192.168.0.1/24", PeerNode: "node02", PeerInf: "net0"}
	inf02 := utils.Interface{InfName: "net0", Type: "direct", Ipv4: "192.168.0.2/24", PeerNode: "node01", PeerInf: "net0"}
	node01 := utils.Node{Name: "node01", Image: "akiranet24/frr:1.0", Interface: []utils.Interface{inf01}}
	node02 := utils.Node{Name: "node02", Image: "akiranet24/frr:1.0", Interface: []utils.Interface{inf02}}
	nodes := []utils.Node{node01, node02}
	if err := c.Pull(nodes); err != nil {
		t.Fatalf("failed test %#v", err)
	}
}

func TestDockertonetnsAddDel(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	c := NewContainer(ctx, cli)

	inf01 := utils.Interface{InfName: "net0", Type: "direct", Ipv4: "192.168.0.1/24", PeerNode: "node02", PeerInf: "net0"}
	inf02 := utils.Interface{InfName: "net0", Type: "direct", Ipv4: "192.168.0.2/24", PeerNode: "node01", PeerInf: "net0"}
	node01 := utils.Node{Name: "node01", Image: "akiranet24/frr:1.0", Interface: []utils.Interface{inf01}}
	node02 := utils.Node{Name: "node02", Image: "akiranet24/frr:1.0", Interface: []utils.Interface{inf02}}
	nodes := []utils.Node{node01, node02}

	for _, node := range nodes {
		if err = c.Dockertonetns(node.Name); err != nil {
			t.Fatalf("failed test %#v", err)
		}
	}

	for _, node := range nodes {
		if err := c.RemoveNs(node.Name); err != nil {
			t.Fatalf("failed test %#v", err)
		}
	}
}

func TestLinkAddDel(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	c := NewContainer(ctx, cli)

	inf01 := utils.Interface{InfName: "net0", Type: "direct", Ipv4: "192.168.0.1/24", PeerNode: "node02", PeerInf: "net0"}
	inf02 := utils.Interface{InfName: "net0", Type: "direct", Ipv4: "192.168.0.2/24", PeerNode: "node01", PeerInf: "net0"}
	node01 := utils.Node{Name: "node01", Image: "akiranet24/frr:1.0", Interface: []utils.Interface{inf01}}
	node02 := utils.Node{Name: "node02", Image: "akiranet24/frr:1.0", Interface: []utils.Interface{inf02}}
	nodes := []utils.Node{node01, node02}

	if err := c.Pull(nodes); err != nil {
		t.Fatalf("failed test %#v", err)
	}

	for _, node := range nodes {
		if err = c.Dockertonetns(node.Name); err != nil {
			t.Fatalf("failed test %#v", err)
		}
	}

	for _, node := range nodes {
		for _, inf := range node.Interface {
			if inf.Type == "direct" {
				if err := SetLink(node, inf); err != nil {
					t.Fatalf("failed test %#v", err)
				}
			}
		}
	}

	for _, node := range nodes {
		if err := c.RemoveNs(node.Name); err != nil {
			t.Fatalf("failed test %#v", err)
		}
	}
}

// TestSetConfAddDel
func TestSetConfSuccess(t *testing.T) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	c := NewContainer(ctx, cli)

	inf01 := utils.Interface{InfName: "net0", Type: "direct", Ipv4: "192.168.0.1/24", PeerNode: "node02", PeerInf: "net0"}
	inf02 := utils.Interface{InfName: "net0", Type: "direct", Ipv4: "192.168.0.2/24", PeerNode: "node01", PeerInf: "net0"}
	node01 := utils.Node{Name: "node01", Image: "akiranet24/frr:1.0", Interface: []utils.Interface{inf01}}
	node02 := utils.Node{Name: "node02", Image: "akiranet24/frr:1.0", Interface: []utils.Interface{inf02}}
	nodes := []utils.Node{node01, node02}

	if err := c.Pull(nodes); err != nil {
		t.Fatalf("failed test %#v", err)
	}

	for _, node := range nodes {
		if err = c.Dockertonetns(node.Name); err != nil {
			t.Fatalf("failed test %#v", err)
		}
	}

	for _, node := range nodes {
		for _, inf := range node.Interface {
			if inf.Type == "direct" {
				if err := SetLink(node, inf); err != nil {
					t.Fatalf("failed test %#v", err)
				}
			}
		}
	}

	node01Cmd := []utils.Cmd{
		utils.Cmd{Cmd: "/etc/init.d/frr start"},
		utils.Cmd{Cmd: "vtysh -c 'conf t' -c 'router bgp 100' -c 'bgp router-id 1.1.1.1' -c 'neighbor 192.168.0.2 remote-as 200'"},
	}
	node01Config := utils.Nodeconfig{Name: "node01", Cmds: node01Cmd}
	node02Cmd := []utils.Cmd{
		utils.Cmd{Cmd: "/etc/init.d/frr start"},
		utils.Cmd{Cmd: "vtysh -c 'conf t' -c 'router bgp 200' -c 'bgp router-id 2.2.2.2' -c 'neighbor 192.168.0.1 remote-as 1000'"},
	}
	node02Config := utils.Nodeconfig{Name: "node02", Cmds: node02Cmd}
	configs := []utils.Nodeconfig{node01Config, node02Config}

	for _, config := range configs {
		for _, cmd := range config.Cmds {
			if err := c.SetConf(config.Name, cmd.Cmd); err != nil {
				t.Fatalf("failed test %#v", err)
			}
		}
	}

	for _, node := range nodes {
		if err := c.RemoveNs(node.Name); err != nil {
			t.Fatalf("failed test %#v", err)
		}
	}
}

// TestSetBrAddDel
// TestLinkUp
