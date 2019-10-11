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

// confCmd represents the conf command
var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "set config to container",
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

		configs := utils.ParseConfig(cfgFile)

		if apion {
			c := api.NewContainer(ctx, cli)

			for _, config := range configs {
				for _, cmd := range config.Cmds {
					c.SetConf(config.Name, cmd.Cmd)
				}
			}

			fmt.Println("Success conf microns!")
			return nil
		} else if shellon {

			for _, config := range configs {
				runcmds := shell.RunCmd(config)
				for _, runcmd := range runcmds {
					fmt.Println(runcmd)
				}
			}

			fmt.Println("echo 'Success conf microns!'")
			return nil
		} else {
			fmt.Println("Plese set --api(-a) or --shell(-s) flag")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(confCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// confCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// confCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
