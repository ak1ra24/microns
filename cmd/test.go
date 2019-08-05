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

	"github.com/ak1ra24/microns/shell"
	"github.com/ak1ra24/microns/utils"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Execute test from config",
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(cfgFile) == 0 {
			fmt.Println("Must Set CONFIG YAML")
			os.Exit(1)
		}
		// fmt.Println(cfgFile)
		tests := utils.ParseTest(cfgFile)

		fmt.Println("echo '----------------------------------------------'")
		fmt.Println("echo '                   test                       '")
		fmt.Println("echo '----------------------------------------------'")

		var runtestcmds []string
		for _, test := range tests {
			if len(args) == 0 {
				runtestcmds = shell.RunTestCmd(test)
				for _, runtestcmd := range runtestcmds {
					fmt.Println(runtestcmd)
				}
			} else if test.Name == args[0] {
				runtestcmds = shell.RunTestCmd(test)
				for _, runtestcmd := range runtestcmds {
					fmt.Println(runtestcmd)
				}
			}
		}
		fmt.Println("echo 'Success test microns!'")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
