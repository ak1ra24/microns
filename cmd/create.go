/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cmd is microns command line tools
package cmd

import (
	"fmt"
	"os"

	"github.com/ak1ra24/microns/api"
	"github.com/ak1ra24/microns/shell"
	"github.com/ak1ra24/microns/utils"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// var configFile string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create docker container and ns topology",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(cfgFile) == 0 {
			fmt.Println("Must Set CONFIG YAML")
			os.Exit(1)
		}
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return err
		}

		nodes := utils.ParseNodes(cfgFile)
		configs := utils.ParseConfig(cfgFile)
		switches := utils.ParseSwitch(cfgFile)

		if apion {
			c := api.NewContainer(ctx, cli)

			c.Pull(nodes)

			for _, node := range nodes {
				c.Dockertonetns(node.Name)
			}
			for _, node := range nodes {
				fmt.Println(node.Interface)
				for _, inf := range node.Interface {
					if inf.Type == "direct" {
						api.SetLink(node, inf)
					} else if inf.Type == "bridge" {
						api.SetBridge(node, inf)
					}
				}
			}

			for _, config := range configs {
				for _, cmd := range config.Cmds {
					c.SetConf(config.Name, cmd.Cmd)
				}
			}

			for _, s := range switches {
				api.LinkUp(s.Name)
				for _, inf := range s.Interfaces {
					switchNode := fmt.Sprintf("%s-%s", inf.PeerNode, s.Name)
					api.LinkUp(switchNode)
				}
			}

			fmt.Println("Success create microns!")
			return nil
		} else if shellon {
			var addAddrv4cmd string
			var addAddrv6cmd string
			for _, node := range nodes {
				runcmd := shell.RunContainer(node)
				fmt.Println(runcmd)
				getpidcmd := shell.GetContainerPid(node.Name)
				fmt.Println(getpidcmd)
				symlinkDo := shell.SymlinkNstoContainer(node.Name)
				fmt.Println(symlinkDo)
			}

			for _, s := range switches {
				addbr := shell.AddBr(s)
				fmt.Println(addbr)
			}

			for _, node := range nodes {
				for _, link := range node.Interface {
					if link.Type == "direct" {
						vethname, addlinkcmd := shell.LinkAdd(node, link)
						fmt.Println(addlinkcmd)
						setLinkNscmd := shell.LinkSetNs(vethname, node.Name)
						fmt.Println(setLinkNscmd)
						if link.Ipv4 != "" {
							addAddrv4cmd = shell.AddrAddv4(node.Name, vethname, link)
							fmt.Println(addAddrv4cmd)
						}
						if link.Ipv6 != "" {
							addAddrv6cmd = shell.AddrAddv6(node.Name, vethname, link)
							fmt.Println(addAddrv6cmd)
						}
					} else if link.Type == "bridge" {
						linkbrs, brlinkname, brname, vethname := shell.LinkAddBr(switches, node, link)
						for _, linkbr := range linkbrs {
							fmt.Println(linkbr)
						}
						setbrlink := shell.LinkSetBridge(brlinkname, brname)
						fmt.Println(setbrlink)
						setLinkNscmd := shell.LinkSetNs(vethname, node.Name)
						fmt.Println(setLinkNscmd)
						linkupbridge := shell.LinkUpBridge(brname)
						fmt.Println(linkupbridge)
						linkupbrlink := shell.LinkUpBrLink(brlinkname)
						fmt.Println(linkupbrlink)
					}
				}
			}

			for _, config := range configs {
				runcmds := shell.RunCmd(config)
				for _, runcmd := range runcmds {
					fmt.Println(runcmd)
				}
			}

			fmt.Println("echo 'Success create microns!'")
			return nil
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	// rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
