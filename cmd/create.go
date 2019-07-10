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
	"github.com/ak1ra24/microns/api/utils"
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
		fmt.Println("create called")

		if len(cfgFile) == 0 {
			fmt.Println("Must Set CONFIG YAML")
			os.Exit(1)
		}
		fmt.Println(cfgFile)
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			return err
		}

		nodes := utils.ParseYaml(cfgFile)

		fmt.Println("----------------------------------------------")
		fmt.Println("                   CREATE                     ")
		fmt.Println("----------------------------------------------")
		api.Pull(ctx, cli, nodes)

		for _, node := range nodes {
			api.Dockertonetns(ctx, cli, node.Name)
		}
		for _, node := range nodes {
			fmt.Println(node.Interface)
			for _, inf := range node.Interface {
				api.SetLink(node, inf)
			}
		}
		fmt.Println("Success create microns!")
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
