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

// var configFile2 string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete docker container and ns topology",
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
		switches := utils.ParseSwitch(cfgFile)
		// remove container and netns
		if apion {
			c := api.NewContainer(ctx, cli)
			for _, node := range nodes {
				c.RemoveNs(node.Name)
			}
			for _, s := range switches {
				api.RemoveBr(s.Name)
			}
			fmt.Println("Success Delete microns!")
		} else if shellon {
			// delete ns and container
			for _, node := range nodes {
				delNscmd := shell.NsDel(node.Name)
				fmt.Println(delNscmd)
				delDockercmd := shell.DockerDel(node.Name)
				fmt.Println(delDockercmd)
			}
			for _, s := range switches {
				delBridge := shell.BridgeDel(s.Name)
				fmt.Println(delBridge)
				result_Bridge := fmt.Sprintf("Delete Bridge %s", s.Name)
				fmt.Println(fmt.Sprintf("echo %s", result_Bridge))
			}
			fmt.Println("echo 'Success Delete microns!'")
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")
	// rootCmd.PersistentFlags().StringVarP(&configFile2, "config2", "d", "", "config file name")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
