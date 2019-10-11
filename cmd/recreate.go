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
	Short: "reconfigure router (stop -> up -> conf)",
	Run: func(cmd *cobra.Command, args []string) {
		if len(cfgFile) == 0 {
			fmt.Println("Must Set CONFIG YAML")
			os.Exit(1)
		}

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
				c.StopContainer(node.Name)
			}

			fmt.Println("Success stop microns!")
			fmt.Println("Waiting for 10 Seconds")
			time.Sleep(10 * time.Second)

			for _, node := range nodes {
				c.StartContainer(node.Name)
			}

			for _, config := range configs {
				for _, cmd := range config.Cmds {
					c.SetConf(config.Name, cmd.Cmd)
				}
			}

			fmt.Println("Success recreate microns!")
		} else if shellon {
			// stop ns and container
			for _, node := range nodes {
				stopDockercmd := shell.DockerStop(node.Name)
				fmt.Println(stopDockercmd)
			}
			fmt.Println("echo 'Success stop microns!'")
			fmt.Println("echo 'Waiting for 10 Seconds'")
			time.Sleep(10 * time.Second)
			for _, node := range nodes {
				startDockercmd := shell.DockerStart(node.Name)
				fmt.Println(startDockercmd)
			}

			for _, config := range configs {
				runcmds := shell.RunCmd(config)
				for _, runcmd := range runcmds {
					fmt.Println(runcmd)
				}
			}

			fmt.Println("echo 'Success recreate microns!'")
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
