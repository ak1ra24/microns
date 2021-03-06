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

	"github.com/ak1ra24/microns/graph"
	"github.com/ak1ra24/microns/utils"
	"github.com/spf13/cobra"
)

var imgFile string

// imageCmd represents the image command
var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "create network topology image file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(cfgFile) == 0 {
			fmt.Println("Must Set CONFIG YAML")
			os.Exit(1)
		}
		nodes := utils.ParseNodes(cfgFile)
		bridges := utils.ParseSwitch(cfgFile)
		graph.Graph(nodes, bridges, imgFile)
		graph.DottoPng(imgFile)
		fmt.Println("Success create image file")
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageCmd.PersistentFlags().String("foo", "", "A help for foo")
	rootCmd.PersistentFlags().StringVarP(&imgFile, "output", "o", "topo", "topology image (filename only)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
