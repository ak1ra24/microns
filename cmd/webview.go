// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/ak1ra24/microns/api"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

// webviewCmd represents the webview command
var webviewCmd = &cobra.Command{
	Use:   "webview",
	Short: "webview",
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

		ok := api.Confirm("\x1b[31mConfirm use port number 8000 for webapp\x1b[0m")

		if ok {
			c := api.NewContainer(ctx, cli)

			if err := c.CreateContainerPort("akiranet24/microns-frontend", "microns-frontend", "8000", "8000", cfgFile); err != nil {
				return err
			}

			fmt.Println("http://<your machine address>:8000/")
		} else {
			fmt.Println("Finish webview command...")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(webviewCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webviewCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webviewCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
