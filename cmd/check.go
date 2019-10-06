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

	"github.com/ak1ra24/microns/utils"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check node and interface for config",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(cfgFile) == 0 {
			fmt.Println("Must Set CONFIG YAML")
			os.Exit(1)
		}

		nodes := utils.ParseNodes(cfgFile)
		confmap := map[string]string{}

		for _, node := range nodes {
			for _, inf := range node.Interface {
				// fmt.Println(node.Name, " : ", inf.InfName, "->", inf.PeerNode, " : ", inf.PeerInf)
				host := node.Name + ":" + inf.InfName
				target := inf.PeerNode + ":" + inf.PeerInf
				confmap[host] = target
			}
		}

		var matchNum int
		falseConfigMap := map[string]string{}

		for key, value := range confmap {
			if confmap[key] == value && confmap[value] == key {
				matchNum += 1
			} else {
				falseConfigMap[key] = value
			}
		}

		if len(confmap) == matchNum {
			fmt.Println("Success Check!")
		} else {
			return fmt.Errorf("Failed Check: %s\n", falseConfigMap)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
