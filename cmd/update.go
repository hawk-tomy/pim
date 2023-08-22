/*
MIT License

# Copyright (c) 2023 - present hawk-tomy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"fmt"

	"github.com/hawk-tomy/python_install_manager/lib"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update python",
	Long:  "update python",

	Args: cobra.MaximumNArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		isAllVer, err := cmd.Flags().GetBool("all")
		cobra.CheckErr(err)
		if !isAllVer {
			if len(args) != 1 {
				fmt.Println("accept only 1 argument.")
			}
		} else {
			if len(args) != 0 {
				fmt.Println("accept no argument.")
			}
		}

		if isAllVer {
			lib.UpdateAll(config)
		} else {
			version, err := lib.NewVersion(args[0])
			cobra.CheckErr(err)
			lib.UpdateLatest(config, version)
		}

	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	updateCmd.Flags().BoolP("all", "A", false, "update all updatable python")
}
