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
	"time"

	"github.com/ak1ra24/microns/api"
	"github.com/ak1ra24/microns/shell"
	"github.com/ak1ra24/microns/utils"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// recreateCmd represents the recreate command
var recreateCmd = &cobra.Command{
	Use:   "recreate",
	Short: "reconfigure router",
	Long: `if you change conf, execute this command.
	this command is that delete and create.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(cfgFile) == 0 {
			fmt.Println("Must Set CONFIG YAML")
			os.Exit(1)
		}
		// fmt.Println(cfgFile)

		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			panic(err)
		}

		nodes := utils.ParseNodes(cfgFile)
		configs := utils.ParseConfig(cfgFile)
		// remove container and netns
		if apion {
			c := api.NewContainer(ctx, cli)
			for _, node := range nodes {
				c.RemoveNs(node.Name)
			}

			fmt.Println("Success Delete microns!")
			fmt.Println("Waiting for 10 Seconds")
			time.Sleep(10 * time.Second)

			c.Pull(nodes)

			for _, node := range nodes {
				c.Dockertonetns(node.Name)
			}
			for _, node := range nodes {
				fmt.Println(node.Interface)
				for _, inf := range node.Interface {
					api.SetLink(node, inf)
				}
			}
			fmt.Println("Success create microns!")
		} else if shellon {
			// delete ns and container
			for _, node := range nodes {
				delNscmd := shell.NsDel(node.Name)
				fmt.Println(delNscmd)
				delDockercmd := shell.DockerDel(node.Name)
				fmt.Println(delDockercmd)
			}
			fmt.Println("echo 'Success Delete microns!'")
			fmt.Println("echo 'Waiting for 10 Seconds'")
			time.Sleep(10 * time.Second)
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
			for _, node := range nodes {
				for _, link := range node.Interface {
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
				}
			}

			for _, config := range configs {
				runcmds := shell.RunCmd(config)
				for _, runcmd := range runcmds {
					fmt.Println(runcmd)
				}
			}

			fmt.Println("echo 'Success create microns!'")
		}
	},
}

func init() {
	rootCmd.AddCommand(recreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recreateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recreateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
