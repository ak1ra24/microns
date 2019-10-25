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

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "up docker container and ns topology",
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

		switches := utils.ParseSwitch(cfgFile)

		if apion {
			c := api.NewContainer(ctx, cli)

			for _, node := range nodes {
				c.StartContainer(node.Name)
			}

			fmt.Println("Success up microns!")

			return nil

		} else if shellon {
			var addAddrv4cmd string
			var addAddrv6cmd string

			for _, node := range nodes {
				startDockercmd := shell.DockerStart(node.Name)
				fmt.Println(startDockercmd)
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

			fmt.Println("echo 'Success up microns!'")
			return nil
		} else {
			fmt.Println("Plese set --api(-a) or --shell(-s) flag")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upCmd.PersistentFlags().String("foo", "", "A help for foo")
	// rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "config file name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
